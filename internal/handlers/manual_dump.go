package handlers

import (
	"seall/internal/api/mediaapi"
	"seall/internal/library/scanner"
	"seall/internal/util/limiter"

	"github.com/labstack/echo/v4"
)

// DUMMY HANDLER

type RequestBody struct {
	Dir      string `json:"dir"`
	Username string `json:"userName"`
}

// HandleTestDump
//
//	@summary this is a dummy handler for testing purposes.
//	@route /api/v1/test-dump [POST]
func (h *Handler) HandleTestDump(c echo.Context) error {

	body := new(RequestBody)
	if err := c.Bind(body); err != nil {
		return h.RespondWithError(c, err)
	}

	if err := h.guardStrictLocalOnlyAction(c); err != nil {
		return err
	}

	localFiles, err := scanner.GetLocalFilesFromDir(body.Dir, h.App.Logger)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	completeAnimeCache := mediaapi.NewCompleteAnimeCache()

	mc, err := scanner.NewMediaFetcher(c.Request().Context(), &scanner.MediaFetcherOptions{
		Enhanced:               false,
		PlatformRef:            h.App.MediaPlatformRef,
		MetadataProviderRef:    h.App.MetadataProviderRef,
		LocalFiles:             localFiles,
		CompleteAnimeCache:     completeAnimeCache,
		Logger:                 h.App.Logger,
		MediaApiRateLimiter:    limiter.NewMediaApiLimiter(),
		DisableAnimeCollection: false,
		ScanLogger:             nil,
	})

	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, mc.AllMedia)
}
