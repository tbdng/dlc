// Package wallet project wallet.go
package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
)

// Wallet is hierarchical deterministic wallet
type Wallet struct {
	params chaincfg.Params
	// rpc    *rpc.BtcRPC

	db               walletdb.DB
	Manager          *waddrmgr.Manager
	publicPassphrase []byte
}

// CreateWallet returns a new Wallet, also creates db where wallet resides
// TODO: separate db creation and Manager creation, creature loader script for
// wallet init
func CreateWallet(params chaincfg.Params, seed, pubPass, privPass []byte, dbFilePath, walletName string) (*Wallet, error) {
	wallet := &Wallet{}
	wallet.params = params
	// wallet.rpc = rpc
	wallet.publicPassphrase = pubPass

	// TODO: add prompts for dbDirPath, walletDBname
	dbDirPath := filepath.Join(dbFilePath, params.Name)
	walletDBname := walletName + ".db"
	dbPath := filepath.Join(dbDirPath, walletDBname)
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("something already exists on this filepath")
	}
	err = os.MkdirAll(dbDirPath, 0700)
	if err != nil {
		return nil, err
	}

	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		return nil, err
	}
	wallet.db = db

	var mgr *waddrmgr.Manager
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
		if e != nil {
			return e
		}

		birthday := time.Now()
		e = waddrmgr.Create(
			addrmgrNs, seed, pubPass, privPass, &params, nil,
			birthday,
		)
		if e != nil {
			// TODO: figure out how to gracefully close db
			//   possibly defer db.Close() ?
			db.Close()
			return e
		}
		mgr, e = waddrmgr.Open(addrmgrNs, pubPass, &params)
		wallet.Manager = mgr

		return e
	})
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// TODO: add Open wallet function
// TODO: add Close wallet function that will gracefully close db, Manager...

// CreateAccount creates a new account in ScopedKeyManagar of scope
func (w *Wallet) CreateAccount(scope waddrmgr.KeyScope, name string, privPass []byte) (uint32, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.Manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return 0, err
	}

	scopedMgr, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		return 0, err
	}

	var account uint32
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		account, e = scopedMgr.NewAccount(ns, name)
		return e
	})
	if err != nil {
		return 0, err
	}

	return account, nil
}

// NewAddress returns a new ManagedAddress for a given scope and account number.
// NOTE: this function callsNextExternalAddresses to generate a ManagadAdddress.
func (w *Wallet) NewAddress(scope waddrmgr.KeyScope, privPass []byte,
	account uint32, numAddresses uint32) ([]waddrmgr.ManagedAddress, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.Manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return nil, err
	}

	// get ScopedKeyManager
	scopedMgr, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		return nil, err
	}

	var addrs []waddrmgr.ManagedAddress
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		addrs, e = scopedMgr.NextExternalAddresses(ns, account, numAddresses)
		return e
	})
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// Helper function, TODO: move somewhere else?
func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}