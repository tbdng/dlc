package oracle

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// Oracle is a struct
type Oracle struct {
	name      string                  // display name
	nRpoints  int                     // number of committed R-points
	masterKey *hdkeychain.ExtendedKey // master HD extended key (private)
	db        *memdb                  // memory db for testing
}

// New creates a oracle
func New(name string, params *chaincfg.Params, nRpoints int) (*Oracle, error) {
	// TODO: pass master key from outside
	mKey, err := randMasterKey(name, params)
	if err != nil {
		return nil, err
	}

	oracle := &Oracle{name: name, nRpoints: nRpoints, masterKey: mKey}
	return oracle, nil
}

func isMainNet(params *chaincfg.Params) bool {
	return params.Net == chaincfg.MainNetParams.Net
}

// randMasterKey creates oracle's random master key
func randMasterKey(name string, params *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	seed := chainhash.DoubleHashB([]byte(name))
	return hdkeychain.NewMaster(seed, params)
}
