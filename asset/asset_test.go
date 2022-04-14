package asset

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

var (
	hashBytes1     = [32]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	hashBytes2     = [32]byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	pubKeyBytes, _ = hex.DecodeString(
		"a0afeb165f0ec36880b68e0baabd9ad9c62fd1a69aa998bc30e9a346202e078f",
	)
	pubKey, _   = schnorr.ParsePubKey(pubKeyBytes)
	sigBytes, _ = hex.DecodeString(
		"e907831f80848d1069a5371b402410364bdf1c5f8307b0084c55f1ce2dca821525f66a4a85ea8b71e482a74f382d2ce5ebeee8fdb2172f477df4900d310536c0",
	)
	sig, _ = schnorr.ParseSignature(sigBytes)
)

func assertAssetEqual(t *testing.T, a, b *Asset) {
	t.Helper()

	require.Equal(t, a.Version, b.Version)
	require.Equal(t, a.Genesis, b.Genesis)
	require.Equal(t, a.Type, b.Type)
	require.Equal(t, a.Amount, b.Amount)
	require.Equal(t, a.LockTime, b.LockTime)
	require.Equal(t, a.RelativeLockTime, b.RelativeLockTime)
	require.Equal(t, a.PrevWitnesses, b.PrevWitnesses)
	require.Equal(t, a.SplitCommitmentRoot, b.SplitCommitmentRoot)
	require.Equal(t, a.ScriptVersion, b.ScriptVersion)
	require.Equal(t, a.ScriptKey, b.ScriptKey)
	require.Equal(t, a.FamilyKey, b.FamilyKey)
}

// TestAssetEncoding asserts that we can properly encode and decode assets
// through their TLV serialization.
func TestAssetEncoding(t *testing.T) {
	t.Parallel()

	assertAssetEncoding := func(a *Asset) {
		t.Helper()

		assertAssetEqual(t, a, a.Copy())
		var buf bytes.Buffer
		require.NoError(t, a.Encode(&buf))
		var b Asset
		require.NoError(t, b.Decode(&buf))
		assertAssetEqual(t, a, &b)
	}

	assertAssetEncoding(&Asset{
		Version: 1,
		Genesis: Genesis{
			FirstPrevOut: wire.OutPoint{
				Hash:  hashBytes1,
				Index: 1,
			},
			Tag:         "asset",
			Metadata:    []byte{1, 2, 3},
			OutputIndex: 1,
		},
		Type:             1,
		Amount:           1,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes1,
					Index: 1,
				},
				ID:        hashBytes1,
				ScriptKey: *pubKey,
			},
			TxWitness:            wire.TxWitness{{1}, {1}},
			SplitCommitmentProof: []byte{1, 1},
		}, {
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes2,
					Index: 2,
				},
				ID:        hashBytes2,
				ScriptKey: *pubKey,
			},
			TxWitness:            wire.TxWitness{{2}, {2}},
			SplitCommitmentProof: []byte{2, 2},
		}},
		SplitCommitmentRoot: &SplitCommitmentRoot{
			Key: hashBytes1,
			Sum: 1337,
		},
		ScriptVersion: 1,
		ScriptKey:     *pubKey,
		FamilyKey:     &FamilyKey{Key: *pubKey, Sig: *sig},
	})

	assertAssetEncoding(&Asset{
		Version: 2,
		Genesis: Genesis{
			FirstPrevOut: wire.OutPoint{
				Hash:  hashBytes2,
				Index: 2,
			},
			Tag:         "asset",
			Metadata:    []byte{1, 2, 3},
			OutputIndex: 2,
		},
		Type:             2,
		Amount:           2,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID:               nil,
			TxWitness:            nil,
			SplitCommitmentProof: nil,
		}, {
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes2,
					Index: 2,
				},
				ID:        hashBytes2,
				ScriptKey: *pubKey,
			},
			TxWitness:            wire.TxWitness{{2}, {2}},
			SplitCommitmentProof: nil,
		}},
		SplitCommitmentRoot: nil,
		ScriptVersion:       2,
		ScriptKey:           *pubKey,
		FamilyKey:           nil,
	})
}
