package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/p2pderivatives/dlc/pkg/script"
)

// ContractExecutionTx constructs a contract execution tx (CET) using pubkeys and given condition.
// Both parties have different transactions signed by the other side.
//
// txins:
//   [0]:fund transaction output[0]
// txouts:
//   [0]:settlement script
//   [1]:p2wpkh (option)
func (d *DLC) ContractExecutionTx(
	party Contractor, deal *Deal, dID int) (*wire.MsgTx, error) {
	cparty := counterparty(party)

	// out values
	damt1 := deal.Amts[party]
	damt2 := deal.Amts[cparty]

	if damt1 == 0 {
		return d.ContractAbandonmentTx(party)
	}

	// txout1: contract execution script
	pub1 := d.Pubs[party]
	if pub1 == nil {
		return nil, errors.New("missing pubkey")
	}
	pub2 := d.Pubs[cparty]
	if pub2 == nil {
		return nil, errors.New("missing pubkey")
	}

	C := d.Oracle.Commitments[dID]
	if C == nil {
		return nil, errors.New("missing oracle's commitment")
	}

	sc, err := script.ContractExecutionScript(pub1, pub2, C)
	if err != nil {
		return nil, err
	}
	pkScript, err := script.P2WSHpkScript(sc)
	if err != nil {
		return nil, err
	}

	tx, err := d.newRedeemTx()
	if err != nil {
		return nil, err
	}

	outAmt1 := damt1 + d.closignTxFee()
	txout1 := wire.NewTxOut(int64(outAmt1), pkScript)
	tx.AddTxOut(txout1)

	// txout2: counterparty's p2wpkh
	if damt2 > 0 {
		txout2, err := d.distTxOut(cparty, damt2)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txout2)
	}
	return tx, nil
}

// ContractAbandonmentTx creates tx that sends all fund to the counterparty
// Note: This transaction isn't useful in the realworld, but is necessary for PoC
func (d *DLC) ContractAbandonmentTx(p Contractor) (*wire.MsgTx, error) {
	famt, err := d.fundAmount()
	if err != nil {
		return nil, err
	}

	tx, err := d.newRedeemTx()
	if err != nil {
		return nil, err
	}

	// return closing tx fee to the counterparty
	cp := counterparty(p)
	exterFee := d.closignTxFee()
	txout, err := d.distTxOut(cp, famt+exterFee)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txout)

	return tx, nil
}

// SignContractExecutionTxs signs contract execution txs for all deals
func (b *Builder) SignContractExecutionTxs() ([][]byte, error) {
	var sigs [][]byte
	for idx, deal := range b.Contract.Conds.Deals {
		sign, err := b.SignContractExecutionTx(deal, idx)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, sign)
	}
	return sigs, nil
}

// SignContractExecutionTx signs a contract execution tx for a given party
func (b *Builder) SignContractExecutionTx(deal *Deal, idx int) ([]byte, error) {
	cparty := counterparty(b.party)

	tx, err := b.Contract.ContractExecutionTx(cparty, deal, idx)
	if err != nil {
		return nil, err
	}

	return b.witsigForFundScript(tx)
}

// AcceptCETxSignatures accepts CETx signatures received from the counterparty
func (b *Builder) AcceptCETxSignatures(sigs [][]byte) error {
	for idx, sig := range sigs {
		err := b.Contract.AcceptCETxSignature(b.party, idx, sig)
		if err != nil {
			return err
		}
	}
	return nil
}

// AcceptCETxSignature sets a signature if it's valid for an identified CETx
func (d *DLC) AcceptCETxSignature(party Contractor, idx int, sig []byte) error {
	deal, err := d.Deal(idx)
	if err != nil {
		return err
	}

	tx, err := d.ContractExecutionTx(party, deal, idx)
	if err != nil {
		return err
	}

	err = d.verifyCETxSignature(party, tx, sig)
	if err != nil {
		return err
	}

	d.ExecSigs[idx] = sig
	return nil
}

func (d *DLC) verifyCETxSignature(
	p Contractor, tx *wire.MsgTx, sig []byte) error {

	cparty := counterparty(p)

	fsc, err := d.fundScript()
	if err != nil {
		return err
	}

	sighashes := txscript.NewTxSigHashes(tx)

	ftx, err := d.FundTx()
	if err != nil {
		return err
	}

	fout := ftx.TxOut[fundTxOutAt]

	hash, err := txscript.CalcWitnessSigHash(
		fsc, sighashes, txscript.SigHashAll, tx, fundTxInAt, fout.Value)
	if err != nil {
		return err
	}

	s, err := btcec.ParseDERSignature(sig, btcec.S256())
	if err != nil {
		return err
	}

	if !s.Verify(hash, d.Pubs[cparty]) {
		return errors.New("failed to verify")
	}

	return nil
}

// FixedContractExecutionTx returns fixed CETx
func (d *DLC) FixedContractExecutionTx(p Contractor) (*wire.MsgTx, error) {
	if !d.HasDealFixed() {
		return nil, newNoFixedDealError()
	}

	dID, deal, err := d.FixedDeal()
	if err != nil {
		return nil, err
	}

	return d.ContractExecutionTx(p, deal, dID)
}

// SignedContractExecutionTx returns a contract execution tx signed by both parties
// TODO: separate SignedContractExecutionTx and SignContractExecutionTx
func (b *Builder) SignedContractExecutionTx() (*wire.MsgTx, error) {
	tx, err := b.Contract.FixedContractExecutionTx(b.party)
	if err != nil {
		return nil, err
	}

	sig, err := b.witsigForFundScript(tx)
	if err != nil {
		return nil, err
	}

	dID, _, err := b.Contract.FixedDeal()
	if err != nil {
		return nil, err
	}

	cpSig := b.Contract.ExecSigs[dID]

	var sig1, sig2 []byte
	switch b.party {
	case FirstParty:
		sig1, sig2 = sig, cpSig
	case SecondParty:
		sig1, sig2 = cpSig, sig
	}

	wit, err := b.Contract.witnessForFundScript(sig1, sig2)
	if err != nil {
		return nil, err
	}

	tx.TxIn[fundTxInAt].Witness = wit

	return tx, nil
}
