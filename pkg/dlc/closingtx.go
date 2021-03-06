package dlc

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/pkg/script"
	"github.com/p2pderivatives/dlc/pkg/wallet"
)

// closingTxOutAt is a txout index of contract execution tx
const closingTxOutAt = 0

// ClosingTx constructs a tx that redeems a given CET
func (d *DLC) ClosingTx(
	p Contractor, cetx *wire.MsgTx) (*wire.MsgTx, error) {

	tx := wire.NewMsgTx(txVersion)

	// txin
	txid := cetx.TxHash()
	txin := wire.NewTxIn(
		wire.NewOutPoint(&txid, closingTxOutAt), nil, nil)

	tx.AddTxIn(txin)

	in := btcutil.Amount(cetx.TxOut[closingTxOutAt].Value)
	fee := d.closignTxFee()
	out := in - fee

	if out <= 0 {
		return nil, newNotEnoughFeesError(in, fee)
	}

	txout, err := d.distTxOut(p, out)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txout)

	return tx, nil
}

// SignedClosingTx constructs a closing tx with witness
func (b *Builder) SignedClosingTx(cetx *wire.MsgTx) (*wire.MsgTx, error) {
	dID, _, err := b.Contract.FixedDeal()
	if err != nil {
		return nil, err
	}
	C := b.Contract.Oracle.Commitments[dID]

	tx, err := b.Contract.ClosingTx(b.party, cetx)
	if err != nil {
		return nil, err
	}

	wit, err := b.witnessForCEScript(tx, cetx, C)
	if err != nil {
		return nil, err
	}
	tx.TxIn[0].Witness = wit

	return tx, nil
}

func (b *Builder) witnessForCEScript(
	tx, cetx *wire.MsgTx, C *btcec.PublicKey) (wire.TxWitness, error) {
	cetxout := cetx.TxOut[closingTxOutAt]
	amt := btcutil.Amount(cetxout.Value)

	cparty := counterparty(b.party)
	pub1, pub2 := b.Contract.Pubs[b.party], b.Contract.Pubs[cparty]

	sc, err := script.ContractExecutionScript(
		pub1, pub2, C)
	if err != nil {
		return nil, err
	}

	// callback function that adds message sig to private key
	osig := b.Contract.Oracle.Sig
	privkeyConverter := genAddSigToPrivkeyFunc(osig)

	sig, err := b.wallet.WitnessSignatureWithCallback(
		tx, closingTxOutAt, amt, sc, pub1, privkeyConverter)
	if err != nil {
		return nil, err
	}

	wit := script.WitnessForCEScript(sig, sc)
	return wit, nil
}

func genAddSigToPrivkeyFunc(
	sig []byte) wallet.PrivateKeyConverter {
	return func(priv *btcec.PrivateKey) (*btcec.PrivateKey, error) {
		n := new(big.Int).Add(priv.D, new(big.Int).SetBytes(sig))
		n = new(big.Int).Mod(n, btcec.S256().N)
		p, _ := btcec.PrivKeyFromBytes(btcec.S256(), n.Bytes())
		return p, nil
	}
}

// ExecuteContract sends CETx and closing tx
func (b *Builder) ExecuteContract() error {
	cetx, err := b.SignedContractExecutionTx()
	if err != nil {
		return err
	}
	cltx, err := b.SignedClosingTx(cetx)
	if err != nil {
		return err
	}

	_, err = b.wallet.SendRawTransaction(cetx)
	if err != nil {
		return err
	}

	_, err = b.wallet.SendRawTransaction(cltx)
	return err
}
