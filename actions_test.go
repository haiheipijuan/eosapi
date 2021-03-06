package eosapi

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/eosioca/eosapi/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionNewAccount(t *testing.T) {
	pubKey, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	require.NoError(t, err)
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{AccountName("eosio"), PermissionName("active")},
		},
		Data: NewAccount{
			Creator: AccountName("eosio"),
			Name:    AccountName("abourget"),
			Owner: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: pubKey,
						Weight:    1,
					},
				},
			},
			Active: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: pubKey,
						Weight:    1,
					},
				},
			},
			Recovery: Authority{
				Threshold: 1,
				Accounts: []PermissionLevelWeight{
					PermissionLevelWeight{
						Permission: PermissionLevel{AccountName("eosio"), PermissionName("active")},
						Weight:     1,
					},
				},
			},
		},
	}
	tx := &Transaction{
		Actions: []*Action{a},
	}

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)

	assert.Equal(t, `00096e8800000000000000000000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}`, string(buf))

}

func TestMarshalTransactionAndSigned(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{AccountName("eosio"), PermissionName("active")},
		},
		Data: NewAccount{
			Creator: AccountName("eosio"),
			Name:    AccountName("abourget"),
		},
	}
	tx := &SignedTransaction{Transaction: &Transaction{
		Actions: []*Action{a},
	}}

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)
	// 00096e88 0000 0000 00000000 0000 0000 00
	// actions: 01
	// 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
	// len: 22
	// 0000000000ea3055 00000059b1abe931 000000000000000000000000000000000000

	assert.Equal(t, `00096e8800000000000000000000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232220000000000ea305500000059b1abe9310000000000000000000000000000000000000000`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000059b1abe931000000000000000000000000000000000000","name":"newaccount"}`, string(buf))
}

// FETCHED FROM A SIMILAR TRANSACTION VIA `eosioc`, includes the Transaction headers though:
// This was BEFORE the `keys` and `accounts` were swapped on `Authority`.
// transaction header:
//    expiration epoch: 1e76ac5a
//    region: 0000
//    blocknum: 62cf
//    blockprefix: 50090bd8
//    packedbandwidthwords: 0000
//    contexfreecpubandwidth: 0000
//    []ContextFreeActions: 00
//    []Actions: 01
// Action idx 0:
//  account: 0000000000ea3055 (eosio)
//  name: 00409e9a2264b89a (newaccount)
//  []authorizations: 01
//   - actor: 0000000000ea3055 (eosio)
//     permission: 00000000a8ed3232 (active)
//  data len: 7c (124, indeed the length of the following...
//  creator: 0000000000ea3055
//  name: 0000001e4d75af46
//  owner authority:
//   threshold: 01000000
//   []accounts: 00
//   []keys: 01
//     - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf // fixed width.
//       weight: 0100
//  active authority:
//   threshold: 01000000
//   []accounts: 00
//   []keys: 01
//     - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//       weight: 0100
//  recovery authority:
//    threshold: 01000000
//    []accounts: 01
//    - name: 0000000000ea3055 (eosio)
//    []keys: 00
// now the `newaccount` struct is done.. what,s that ?
// a list of a new object: 01
// an account name:
// a permission name: 00000000a8ed3232 (active)
// some list with one thing: 01
//   - an empty list: 00
//   - another empty list: 00

// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32326a0000000000ea305500000059b1abe9310100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000000000
// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100
// Generated by the first run:
// account: 0000000000ea3055 (eosio)
// name: 00409e9a2264b89a (newaccount)
// []authorizations: 01
//  - actor: 0000000000ea3055 (eosio)
//    permission: 00000000a8ed3232 (active)
// data length: 6a (106) which MATCHES the lengths to follow.
// NewAccount:
//  creator: 0000000000ea3055 (eosio)
//  name: 00000059b1abe931 (abourget)
// owner-authority:
//  threshold: 01000000
//  []keys: 01
//  - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//    weight: 0100
//  []accounts: 00
// active-authority:
//  threshold: 01000000
//  []keys: 01
//  - pubkey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//    weight: 0100
//  []accounts: 00
// recovery-authority:  // the last bit is the Recovery authority.. it works :)
//  threshold: 00000000
//  []keys: 00
//  []accounts: 00

// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe9310100000001   BINARY SERIALIZER FAILED: 35454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100
