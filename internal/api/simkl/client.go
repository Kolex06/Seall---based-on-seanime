package simkl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

const (
	DefaultAPIBaseURL  = "https://api.simkl.com"
	DefaultAuthBaseURL = "https://simkl.com"
)

var (
	ErrNotAuthenticated = errors.New("simkl: not authenticated")
	ErrMissingClientID  = errors.New("simkl: client_id is required")
	ErrNotFound         = errors.New("simkl: media not found")
)

type Client struct {
	token        string
	clientID     string
	clientSecret string
	redirectURI  string
	cacheDir     string
	apiBaseURL   string
	authBaseURL  string
	httpClient   *http.Client
	logger       *zerolog.Logger
}

type NewClientOptions struct {
	Token        string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	CacheDir     string
	APIBaseURL   string
	AuthBaseURL  string
	HTTPClient   *http.Client
	Logger       *zerolog.Logger
}

func NewClient(opts NewClientOptions) *Client {
	clientID := firstNonEmpty(opts.ClientID, os.Getenv("SIMKL_CLIENT_ID"))
	clientSecret := firstNonEmpty(opts.ClientSecret, os.Getenv("SIMKL_CLIENT_SECRET"))
	redirectURI := firstNonEmpty(opts.RedirectURI, os.Getenv("SIMKL_REDIRECT_URI"))
	apiBaseURL := firstNonEmpty(opts.APIBaseURL, DefaultAPIBaseURL)
	authBaseURL := firstNonEmpty(opts.AuthBaseURL, DefaultAuthBaseURL)
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		token:        opts.Token,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		cacheDir:     opts.CacheDir,
		apiBaseURL:   strings.TrimRight(apiBaseURL, "/"),
		authBaseURL:  strings.TrimRight(authBaseURL, "/"),
		httpClient:   httpClient,
		logger:       opts.Logger,
	}
}

func (c *Client) IsAuthenticated() bool {
	return c != nil && c.token != ""
}

func (c *Client) GetCacheDir() string {
	if c == nil {
		return ""
	}
	return c.cacheDir
}

func (c *Client) ClientID() string {
	if c == nil {
		return ""
	}
	return c.clientID
}

func (c *Client) Token() string {
	if c == nil {
		return ""
	}
	return c.token
}

func (c *Client) SetToken(token string) {
	if c == nil {
		return
	}
	c.token = token
}

func (c *Client) AuthorizeURL(state string) (string, error) {
	if c == nil || c.clientID == "" {
		return "", ErrMissingClientID
	}
	if c.redirectURI == "" {
		return "", errors.New("simkl: redirect_uri is required for OAuth authorization")
	}
	u, err := url.Parse(c.authBaseURL + "/oauth/authorize")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.clientID)
	q.Set("redirect_uri", c.redirectURI)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	if c == nil || c.clientID == "" {
		return nil, ErrMissingClientID
	}
	if c.clientSecret == "" {
		return nil, errors.New("simkl: client_secret is required to exchange an OAuth code")
	}
	if c.redirectURI == "" {
		return nil, errors.New("simkl: redirect_uri is required to exchange an OAuth code")
	}

	payload := map[string]string{
		"code":          code,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  c.redirectURI,
		"grant_type":    "authorization_code",
	}

	var ret TokenResponse
	if err := c.doJSON(ctx, http.MethodPost, "/oauth/token", nil, payload, &ret, false, false); err != nil {
		return nil, err
	}
	if ret.AccessToken != "" {
		c.token = ret.AccessToken
	}
	return &ret, nil
}

func (c *Client) RequestPin(ctx context.Context, redirect string) (*PinCode, error) {
	if c == nil || c.clientID == "" {
		return nil, ErrMissingClientID
	}
	q := url.Values{}
	q.Set("client_id", c.clientID)
	if redirect != "" {
		q.Set("redirect", redirect)
	}

	var ret PinCode
	if err := c.doJSON(ctx, http.MethodGet, "/oauth/pin", q, nil, &ret, false, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) CheckPin(ctx context.Context, userCode string) (*PinStatus, error) {
	if c == nil || c.clientID == "" {
		return nil, ErrMissingClientID
	}
	q := url.Values{}
	q.Set("client_id", c.clientID)

	var ret PinStatus
	if err := c.doJSON(ctx, http.MethodGet, "/oauth/pin/"+url.PathEscape(userCode), q, nil, &ret, false, false); err != nil {
		return nil, err
	}
	if ret.AccessToken != "" {
		c.token = ret.AccessToken
	}
	return &ret, nil
}

func (c *Client) Settings(ctx context.Context) (*UserSettings, error) {
	var ret UserSettings
	if err := c.doJSON(ctx, http.MethodPost, "/users/settings", nil, nil, &ret, true, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) Activities(ctx context.Context) (*Activities, error) {
	var ret Activities
	if err := c.doJSON(ctx, http.MethodPost, "/sync/activities", nil, nil, &ret, true, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) AllItems(ctx context.Context, mediaType MediaType, status WatchStatus, query url.Values) (*AllItems, error) {
	if query == nil {
		query = url.Values{}
	}
	requestPath := "/sync/all-items/"
	if mediaType != MediaTypeAll {
		requestPath = path.Join(requestPath, string(mediaType)) + "/"
	}
	if status != "" {
		requestPath = path.Join(requestPath, string(status)) + "/"
	}

	var ret AllItems
	if err := c.doJSON(ctx, http.MethodGet, requestPath, query, nil, &ret, true, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) MediaDetails(ctx context.Context, mediaType MediaType, id string, extended string) (*StandardMedia, error) {
	if mediaType == MediaTypeAll {
		return nil, errors.New("simkl: media type is required")
	}
	if id == "" {
		return nil, errors.New("simkl: media id is required")
	}

	q := url.Values{}
	if extended != "" {
		q.Set("extended", extended)
	}

	var ret StandardMedia
	requestPath := path.Join("/", mediaDetailsPath(mediaType), id)
	if err := c.doJSON(ctx, http.MethodGet, requestPath, q, nil, &ret, false, true); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) MediaEpisodes(ctx context.Context, mediaType MediaType, id string, extended string) ([]Episode, error) {
	if mediaType != MediaTypeAnime && mediaType != MediaTypeShows {
		return nil, errors.New("simkl: episode lists are only available for anime and tv shows")
	}
	if id == "" {
		return nil, errors.New("simkl: media id is required")
	}

	q := url.Values{}
	if extended != "" {
		q.Set("extended", extended)
	}

	var ret []Episode
	requestPath := path.Join("/", mediaEpisodesPath(mediaType), id)
	if err := c.doJSON(ctx, http.MethodGet, requestPath, q, nil, &ret, false, true); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) AddItems(ctx context.Context, payload AddItemsRequest) (*AddItemsResponse, error) {
	var ret AddItemsResponse
	if err := c.doJSON(ctx, http.MethodPost, "/sync/history", nil, payload, &ret, true, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) AddToList(ctx context.Context, payload AddItemsRequest) (*AddItemsResponse, error) {
	var ret AddItemsResponse
	if err := c.doJSON(ctx, http.MethodPost, "/sync/add-to-list", nil, payload, &ret, true, false); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) RemoveItems(ctx context.Context, payload AddItemsRequest) error {
	return c.doJSON(ctx, http.MethodPost, "/sync/history/remove", nil, payload, nil, true, false)
}

func (c *Client) doJSON(ctx context.Context, method string, requestPath string, query url.Values, body interface{}, out interface{}, authRequired bool, clientIDInQuery bool) error {
	if c == nil {
		return errors.New("simkl: nil client")
	}
	if c.clientID == "" {
		return ErrMissingClientID
	}
	if authRequired && c.token == "" {
		return ErrNotAuthenticated
	}

	fullURL, err := c.buildURL(requestPath, query, clientIDInQuery)
	if err != nil {
		return err
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("simkl-api-key", c.clientID)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if c.logger != nil {
		c.logger.Debug().
			Str("method", method).
			Str("path", requestPath).
			Int("status", resp.StatusCode).
			Str("duration", time.Since(start).Truncate(time.Millisecond).String()).
			Msg("simkl: request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("simkl: %s %s failed with %d: %s", method, requestPath, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if out == nil || len(respBody) == 0 || string(respBody) == "null" {
		return nil
	}

	bodyText := strings.TrimSpace(string(respBody))
	if bodyText == "[]" && !jsonTargetIsSlice(out) {
		return fmt.Errorf("%w: %s", ErrNotFound, requestPath)
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("simkl: decode %s: %w", requestPath, err)
	}
	return nil
}

func jsonTargetIsSlice(out interface{}) bool {
	if out == nil {
		return false
	}
	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return false
	}
	elem := value.Elem()
	return elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array
}

func (c *Client) buildURL(requestPath string, query url.Values, clientIDInQuery bool) (string, error) {
	base := c.apiBaseURL
	u, err := url.Parse(base + "/" + strings.TrimLeft(requestPath, "/"))
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, vals := range query {
		for _, v := range vals {
			q.Add(k, v)
		}
	}
	if clientIDInQuery {
		q.Set("client_id", c.clientID)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func mediaDetailsPath(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeMovies:
		return "movies"
	case MediaTypeShows:
		return "tv"
	case MediaTypeAnime:
		return "anime"
	default:
		return string(mediaType)
	}
}

func mediaEpisodesPath(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeAnime:
		return "anime/episodes"
	case MediaTypeShows:
		return "tv/episodes"
	default:
		return string(mediaType)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func ImageURL(kind ImageKind, imagePath string, size ImageSize) string {
	if imagePath == "" {
		return ""
	}
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		return imagePath
	}
	imagePath = strings.TrimPrefix(imagePath, "/")
	ext := ".webp"
	if kind == ImageKindAvatar {
		ext = ".jpg"
	}
	return fmt.Sprintf("https://wsrv.nl/?url=https://simkl.in/%s/%s%s%s", kind, imagePath, size, ext)
}
