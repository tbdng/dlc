// Code generated by mockery v1.0.0. DO NOT EDIT.

package rpcmock

import btcjson "github.com/btcsuite/btcd/btcjson"
import btcutil "github.com/btcsuite/btcutil"
import chainhash "github.com/btcsuite/btcd/chaincfg/chainhash"
import json "encoding/json"
import mock "github.com/stretchr/testify/mock"

import wire "github.com/btcsuite/btcd/wire"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Generate provides a mock function with given fields: numBlocks
func (_m *Client) Generate(numBlocks uint32) ([]*chainhash.Hash, error) {
	ret := _m.Called(numBlocks)

	var r0 []*chainhash.Hash
	if rf, ok := ret.Get(0).(func(uint32) []*chainhash.Hash); ok {
		r0 = rf(numBlocks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*chainhash.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint32) error); ok {
		r1 = rf(numBlocks)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockCount provides a mock function with given fields:
func (_m *Client) GetBlockCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImportAddressRescan provides a mock function with given fields: address, account, rescan
func (_m *Client) ImportAddressRescan(address string, account string, rescan bool) error {
	ret := _m.Called(address, account, rescan)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool) error); ok {
		r0 = rf(address, account, rescan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListUnspentMinMaxAddresses provides a mock function with given fields: minConf, maxConf, addrs
func (_m *Client) ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	ret := _m.Called(minConf, maxConf, addrs)

	var r0 []btcjson.ListUnspentResult
	if rf, ok := ret.Get(0).(func(int, int, []btcutil.Address) []btcjson.ListUnspentResult); ok {
		r0 = rf(minConf, maxConf, addrs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, []btcutil.Address) error); ok {
		r1 = rf(minConf, maxConf, addrs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RawRequest provides a mock function with given fields: method, params
func (_m *Client) RawRequest(method string, params []json.RawMessage) (json.RawMessage, error) {
	ret := _m.Called(method, params)

	var r0 json.RawMessage
	if rf, ok := ret.Get(0).(func(string, []json.RawMessage) json.RawMessage); ok {
		r0 = rf(method, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(json.RawMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []json.RawMessage) error); ok {
		r1 = rf(method, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendRawTransaction provides a mock function with given fields: tx, allowHighFees
func (_m *Client) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	ret := _m.Called(tx, allowHighFees)

	var r0 *chainhash.Hash
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, bool) *chainhash.Hash); ok {
		r0 = rf(tx, allowHighFees)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, bool) error); ok {
		r1 = rf(tx, allowHighFees)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendToAddress provides a mock function with given fields: address, amount
func (_m *Client) SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error) {
	ret := _m.Called(address, amount)

	var r0 *chainhash.Hash
	if rf, ok := ret.Get(0).(func(btcutil.Address, btcutil.Amount) *chainhash.Hash); ok {
		r0 = rf(address, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(btcutil.Address, btcutil.Amount) error); ok {
		r1 = rf(address, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
