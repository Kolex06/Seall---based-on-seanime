package mediaapi

import (
	"seall/internal/testutil"
	"testing"
)

func newLiveMediaApiClient(t testing.TB) MediaApiClient {
	t.Helper()

	cfg := testutil.InitTestProvider(t, testutil.MediaApi(), testutil.Live())
	if cfg.Provider.MediaApiJwt == "" {
		t.Skip("simkl live tests require simkl_jwt")
	}

	return NewMediaApiClient(cfg.Provider.MediaApiJwt, "")
}
