package vrf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"

	sdkkr "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptocodec "github.com/evmos/evmos/v16/crypto/codec"
	enccodec "github.com/evmos/evmos/v16/encoding/codec"
)

var TestCodec amino.Codec

func init() {
	cdc := amino.NewLegacyAmino()
	cryptocodec.RegisterCrypto(cdc)

	interfaceRegistry := types.NewInterfaceRegistry()
	TestCodec = amino.NewProtoCodec(interfaceRegistry)
	enccodec.RegisterInterfaces(interfaceRegistry)
}

func TestKeyring(t *testing.T) {
	dir := t.TempDir()
	mockIn := strings.NewReader("")
	kr, err := sdkkr.New("evmos", sdkkr.BackendTest, dir, mockIn, TestCodec, VrfOption())
	require.NoError(t, err)

	// fail in retrieving key
	info, err := kr.Key("foo")
	require.Error(t, err)
	require.Nil(t, info)

	keyringAlgos, _ := kr.SupportedAlgorithms()
	algo, err := sdkkr.NewSigningAlgoFromString("vrf", keyringAlgos)
	require.NoError(t, err)

	mockIn.Reset("password\npassword\n")
	newRecord, err := kr.NewAccount("foo", "", "", "", algo)
	require.NoError(t, err)
	require.Equal(t, "foo", newRecord.Name)
	require.Equal(t, "local", newRecord.GetType().String())
	pubKey, err := newRecord.GetPubKey()
	require.NoError(t, err)
	require.Equal(t, string(VrfType), pubKey.Type())

	bz, err := VrfAlgo.Derive()("", "", "")
	require.NoError(t, err)
	require.NotEmpty(t, bz)
}

func TestDerivation(t *testing.T) {
	bz, err := VrfAlgo.Derive()("", "", "")
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	badBz, err := VrfAlgo.Derive()("", "", "")
	require.NoError(t, err)
	require.NotEmpty(t, badBz)

	require.NotEqual(t, bz, badBz)

	privkey := VrfAlgo.Generate()(bz)
	badPrivKey := VrfAlgo.Generate()(badBz)

	require.False(t, privkey.Equals(badPrivKey))
}
