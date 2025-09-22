package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"blogchain/x/blog/keeper"
	module "blogchain/x/blog/module"
	"blogchain/x/blog/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	// Set up the bech32 prefix
	sdk.GetConfig().SetBech32PrefixForAccount("blogchain", "blogchainpub")
	
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec("blogchain")
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	k, err := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}

// GenerateTestAddresses creates valid bech32 addresses for testing
func GenerateTestAddresses(t *testing.T, count int) []string {
	t.Helper()
	
	addresses := make([]string, count)
	for i := 0; i < count; i++ {
		// Generate a new private key
		privKey := secp256k1.GenPrivKey()
		
		// Get the public key
		pubKey := privKey.PubKey()
		
		// Derive the address
		addr := sdk.AccAddress(pubKey.Address())
		
		// Convert to bech32 string with blogchain prefix
		addresses[i] = addr.String()
	}
	
	return addresses
}

// GetTestAddress returns a single valid test address
func GetTestAddress(t *testing.T) string {
	t.Helper()
	addresses := GenerateTestAddresses(t, 1)
	return addresses[0]
}
