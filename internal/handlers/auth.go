package handlers

import (
	"errors"
	"strings"

	"seall/internal/api/simkl"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// HandleLogin
//
//	@summary logs in the user by saving the SIMKL token in the database.
//	@desc This is called when the token is obtained from SIMKL after logging in.
//	@desc It also fetches the SIMKL user data and saves it in the database.
//	@desc It creates a new handlers.Status and refreshes App modules.
//	@route /api/v1/auth/login [POST]
//	@returns handlers.Status
func (h *Handler) HandleLogin(c echo.Context) error {

	type body struct {
		Token string `json:"token"`
	}

	var b body

	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	if err := h.App.LoginToSimkl(b.Token); err != nil {
		return h.RespondWithError(c, err)
	}

	// Create a new status
	status := h.NewStatus(c)

	// Return new status
	return h.RespondWithData(c, status)

}

// HandleStartSimklPinLogin starts SIMKL PIN authentication.
//
//	@summary starts SIMKL PIN auth.
//	@route /api/v1/auth/simkl/pin [POST]
//	@returns simkl.PinCode
func (h *Handler) HandleStartSimklPinLogin(c echo.Context) error {
	pin, err := h.App.SimklClientRef.Get().RequestPin(c.Request().Context(), "")
	if err != nil {
		if errors.Is(err, simkl.ErrMissingClientID) {
			return h.RespondWithError(c, errors.New("SIMKL client ID is missing. Add one in the login window first."))
		}
		return h.RespondWithError(c, err)
	}
	return h.RespondWithData(c, pin)
}

// HandleCheckSimklPinLogin checks SIMKL PIN authentication and logs in once approved.
//
//	@summary checks SIMKL PIN auth.
//	@route /api/v1/auth/simkl/pin/check [POST]
//	@returns handlers.Status when approved, otherwise simkl.PinStatus
func (h *Handler) HandleCheckSimklPinLogin(c echo.Context) error {
	type body struct {
		UserCode string `json:"userCode"`
	}

	var b body
	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}
	if b.UserCode == "" {
		return h.RespondWithError(c, errors.New("missing SIMKL user code"))
	}

	pinStatus, err := h.App.SimklClientRef.Get().CheckPin(c.Request().Context(), b.UserCode)
	if err != nil {
		if errors.Is(err, simkl.ErrMissingClientID) {
			return h.RespondWithError(c, errors.New("SIMKL client ID is missing. Add one in the login window first."))
		}
		return h.RespondWithError(c, err)
	}

	if pinStatus.AccessToken == "" {
		return h.RespondWithData(c, pinStatus)
	}

	if err := h.App.LoginToSimkl(pinStatus.AccessToken); err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, h.NewStatus(c))
}

// HandleSaveSimklClientConfig saves the SIMKL client configuration used by PIN login.
//
//	@summary saves SIMKL client config.
//	@route /api/v1/auth/simkl/client [PATCH]
//	@returns handlers.Status
func (h *Handler) HandleSaveSimklClientConfig(c echo.Context) error {
	type body struct {
		ClientID     string  `json:"clientId"`
		ClientSecret *string `json:"clientSecret"`
		RedirectURI  *string `json:"redirectUri"`
	}

	var b body
	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	clientID := strings.TrimSpace(b.ClientID)
	if clientID == "" {
		return h.RespondWithError(c, errors.New("SIMKL client ID is required"))
	}

	h.App.Config.Simkl.ClientID = clientID
	viper.Set("simkl.clientId", clientID)

	if b.ClientSecret != nil {
		clientSecret := strings.TrimSpace(*b.ClientSecret)
		h.App.Config.Simkl.ClientSecret = clientSecret
		viper.Set("simkl.clientSecret", clientSecret)
	}

	if b.RedirectURI != nil {
		redirectURI := strings.TrimSpace(*b.RedirectURI)
		h.App.Config.Simkl.RedirectURI = redirectURI
		viper.Set("simkl.redirectUri", redirectURI)
	}

	if err := viper.WriteConfig(); err != nil {
		return h.RespondWithError(c, err)
	}

	h.App.UpdateSimklClientToken(h.App.GetUserSimklToken())

	return h.RespondWithData(c, h.NewStatus(c))
}

// HandleLogout
//
//	@summary logs out the user by removing JWT token from the database.
//	@desc It removes SIMKL token and viewer data from the database.
//	@desc It creates a new handlers.Status and refreshes App modules.
//	@route /api/v1/auth/logout [POST]
//	@returns handlers.Status
func (h *Handler) HandleLogout(c echo.Context) error {
	h.App.LogoutFromSimkl()

	status := h.NewStatus(c)
	return h.RespondWithData(c, status)
}
