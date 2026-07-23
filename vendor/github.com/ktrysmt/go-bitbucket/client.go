package bitbucket

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/bitbucket"
	"golang.org/x/oauth2/clientcredentials"
)

const DEFAULT_PAGE_LENGTH = 10
const DEFAULT_LIMIT_PAGES = 0
const DEFAULT_MAX_DEPTH = 1
const DEFAULT_BITBUCKET_API_BASE_URL = "https://api.bitbucket.org/2.0"

func apiBaseUrlEnv() (*url.URL, error) {
	ev := os.Getenv("BITBUCKET_API_BASE_URL")
	if ev == "" {
		ev = DEFAULT_BITBUCKET_API_BASE_URL
	}

	return url.Parse(ev)
}

func appendCaCerts(caCerts []byte) (*http.Client, error) {
	// 1. If the system standard cert pool exists, create a copy that can be modified.
	caCertPool, err := x509.SystemCertPool()
	// The system cert pool does not exist, so we are going to create a new one.
	if err != nil {
		// The system standard cert pool does not exist so create a new empty one.
		caCertPool = x509.NewCertPool()
	}
	// 2. Append the custom CA certs to the pool.
	if success := caCertPool.AppendCertsFromPEM(caCerts); !success {
		return nil, fmt.Errorf("unable to append CA Certs to cert pool: %w", err)
	}
	// 3. Create a new http.Transport copying http.DefaultTransport
	newTransport := http.DefaultTransport.(*http.Transport).Clone()
	// 4. Append the custom CA certs to the new transport.
	newTransport.TLSClientConfig = &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
	}
	// 5. Create a new http client
	return &http.Client{Transport: newTransport}, nil
}

type Client struct {
	Auth         *auth
	Users        *Users
	User         user
	Teams        teams
	Repositories *Repositories
	Workspaces   *Workspace
	Pagelen      int
	MaxDepth     int
	// LimitPages limits the number of pages for a request
	//	default value as 0 -- disable limits
	LimitPages int
	// DisableAutoPaging allows you to disable the default behavior of automatically requesting
	// all the pages for a paginated response.
	DisableAutoPaging bool
	apiBaseURL        *url.URL

	HttpClient *http.Client
}

type auth struct {
	appID, secret  string
	user, password string
	token          oauth2.Token
	bearerToken    string
	caCerts        []byte
	apiBaseUrl     *url.URL
}

type Response struct {
	*http.Response `json:"-"`
	Size           int           `json:"size"`
	Page           int           `json:"page"`
	Pagelen        int           `json:"pagelen"`
	Next           string        `json:"next"`
	Previous       string        `json:"previous"`
	Values         []interface{} `json:"values"`
}

// Uses the Client Credentials Grant oauth2 flow to authenticate to Bitbucket
func NewOAuthClientCredentials(i, s string) (*Client, error) {
	return NewOAuthClientCredentialsWithEndpoint(i, s, bitbucket.Endpoint.TokenURL)
}

// NewOAuthClientCredentialsWithEndpoint is like NewOAuthClientCredentials but
// targets a custom OAuth token endpoint (e.g. an Isolated Cloud Instance with
// a customer-specific hostname).
func NewOAuthClientCredentialsWithEndpoint(i, s, tokenURL string) (*Client, error) {
	a := &auth{appID: i, secret: s}
	ctx := context.Background()
	conf := &clientcredentials.Config{
		ClientID:     i,
		ClientSecret: s,
		TokenURL:     tokenURL,
	}

	tok, err := conf.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}
	a.token = *tok
	return injectClient(a)

}

// NewOAuth performs an interactive OAuth flow using stdin/stdout.
//
// Deprecated: This function uses stdin/stdout directly, making it unsuitable for
// non-interactive environments (e.g., web servers, background jobs). Instead, use
// NewOAuthWithCode after obtaining the authorization code through your own UI/CLI.
// You can generate the authorization URL using oauth2.Config.AuthCodeURL() directly.
func NewOAuth(i, s string) (*Client, error) {
	return NewOAuthWithEndpoint(i, s, bitbucket.Endpoint)
}

// NewOAuthWithEndpoint is like NewOAuth but targets a custom OAuth endpoint
// (e.g. an Isolated Cloud Instance with a customer-specific hostname).
//
// Deprecated: This function uses stdin/stdout directly, making it unsuitable
// for non-interactive environments. Prefer NewOAuthWithCodeWithEndpoint.
func NewOAuthWithEndpoint(i, s string, ep oauth2.Endpoint) (*Client, error) {
	a := &auth{appID: i, secret: s}
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     i,
		ClientSecret: s,
		Endpoint:     ep,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog:\n%v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	fmt.Printf("Enter the code in the return URL: ")
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("failed to read authorization code: %w", err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	a.token = *tok
	return injectClient(a)
}

// NewOAuthWithCode finishes the OAuth handshake with a given code
// and returns a *Client
func NewOAuthWithCode(i, s, c string) (*Client, string, error) {
	return NewOAuthWithCodeWithEndpoint(i, s, c, bitbucket.Endpoint)
}

// NewOAuthWithCodeWithEndpoint is like NewOAuthWithCode but targets a custom
// OAuth endpoint (e.g. an Isolated Cloud Instance with a customer-specific
// hostname).
func NewOAuthWithCodeWithEndpoint(i, s, c string, ep oauth2.Endpoint) (*Client, string, error) {
	a := &auth{appID: i, secret: s}
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     i,
		ClientSecret: s,
		Endpoint:     ep,
	}

	tok, err := conf.Exchange(ctx, c)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	a.token = *tok
	client, err := injectClient(a)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create client: %w", err)
	}
	return client, tok.AccessToken, nil
}

// NewOAuthWithRefreshToken obtains a new access token with a given refresh token
// and returns a *Client
func NewOAuthWithRefreshToken(i, s, rt string) (*Client, string, error) {
	return NewOAuthWithRefreshTokenWithEndpoint(i, s, rt, bitbucket.Endpoint)
}

// NewOAuthWithRefreshTokenWithEndpoint is like NewOAuthWithRefreshToken but
// targets a custom OAuth endpoint (e.g. an Isolated Cloud Instance with a
// customer-specific hostname).
func NewOAuthWithRefreshTokenWithEndpoint(i, s, rt string, ep oauth2.Endpoint) (*Client, string, error) {
	a := &auth{appID: i, secret: s}
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     i,
		ClientSecret: s,
		Endpoint:     ep,
	}

	tokenSource := conf.TokenSource(ctx, &oauth2.Token{
		RefreshToken: rt,
	})
	tok, err := tokenSource.Token()
	if err != nil {
		return nil, "", fmt.Errorf("failed to refresh token: %w", err)
	}
	a.token = *tok
	client, err := injectClient(a)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create client: %w", err)
	}
	return client, tok.AccessToken, nil
}

func NewOAuthbearerToken(t string) (*Client, error) {
	a := &auth{bearerToken: t}
	return injectClient(a)
}

func NewOAuthbearerTokenWithBaseUrlStr(t, u string) (*Client, error) {
	apiBaseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	a := &auth{bearerToken: t, apiBaseUrl: apiBaseURL}
	return injectClient(a)
}

func NewOAuthbearerTokenWithCaCert(t string, c []byte) (*Client, error) {
	a := &auth{bearerToken: t, caCerts: c}
	return injectClient(a)
}

func NewOAuthbearerTokenWithBaseUrlStrCaCert(t, u string, c []byte) (*Client, error) {
	apiBaseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	a := &auth{bearerToken: t, caCerts: c, apiBaseUrl: apiBaseURL}
	return injectClient(a)
}

// NewBasicAuth returns a Client authenticated via HTTP Basic auth.
//
// Atlassian has deprecated Bitbucket Cloud app passwords in favor of
// Atlassian API tokens. For new integrations, prefer NewAPITokenAuth, which
// delegates to NewBasicAuth but documents the email + API token usage.
func NewBasicAuth(u, p string) (*Client, error) {
	a := &auth{user: u, password: p}
	return injectClient(a)
}

func NewBasicAuthWithBaseUrlStr(u, p, urlStr string) (*Client, error) {
	apiBaseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	a := &auth{user: u, password: p, apiBaseUrl: apiBaseURL}
	return injectClient(a)
}

func NewBasicAuthWithCaCert(u, p string, c []byte) (*Client, error) {
	a := &auth{user: u, password: p, caCerts: c}
	return injectClient(a)
}

func NewBasicAuthWithBaseUrlStrCaCert(u, p, urlStr string, c []byte) (*Client, error) {
	apiBaseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	a := &auth{user: u, password: p, apiBaseUrl: apiBaseURL, caCerts: c}
	return injectClient(a)
}

// NewAPITokenAuth returns a Client authenticated with an Atlassian API token.
// Pass the Atlassian account email as the first argument and the API token as
// the second. Bitbucket Cloud accepts API tokens via HTTP Basic auth, so this
// is a thin alias over NewBasicAuth that makes the intent explicit.
//
// See https://support.atlassian.com/bitbucket-cloud/docs/using-api-tokens/
func NewAPITokenAuth(email, token string) (*Client, error) {
	return NewBasicAuth(email, token)
}

// NewAPITokenAuthWithBaseUrlStr is like NewAPITokenAuth but targets a custom
// API base URL (e.g. an Isolated Cloud Instance with a customer-specific
// hostname). Equivalent to BITBUCKET_API_BASE_URL but configured per client.
func NewAPITokenAuthWithBaseUrlStr(email, token, urlStr string) (*Client, error) {
	return NewBasicAuthWithBaseUrlStr(email, token, urlStr)
}

// NewAPITokenAuthWithCaCert is like NewAPITokenAuth but trusts the supplied
// PEM-encoded CA certificates in addition to the system roots.
func NewAPITokenAuthWithCaCert(email, token string, caCerts []byte) (*Client, error) {
	return NewBasicAuthWithCaCert(email, token, caCerts)
}

// NewAPITokenAuthWithBaseUrlStrCaCert combines a custom API base URL with
// extra trusted CA certificates. Suited for Isolated Cloud Instances behind
// an internal certificate authority.
func NewAPITokenAuthWithBaseUrlStrCaCert(email, token, urlStr string, caCerts []byte) (*Client, error) {
	return NewBasicAuthWithBaseUrlStrCaCert(email, token, urlStr, caCerts)
}

func injectClient(a *auth) (*Client, error) {
	c := &Client{Auth: a, Pagelen: DEFAULT_PAGE_LENGTH, MaxDepth: DEFAULT_MAX_DEPTH,
		LimitPages: DEFAULT_LIMIT_PAGES}
	if a.apiBaseUrl != nil {
		c.apiBaseURL = a.apiBaseUrl
	} else {
		bitbucketUrl, err := apiBaseUrlEnv()
		if err != nil {
			return nil, fmt.Errorf("unable to parse Bitbucket base Url from environment: %w", err)
		}
		c.apiBaseURL = bitbucketUrl
	}
	if a.caCerts != nil {
		httpClient, err := appendCaCerts(a.caCerts)
		if err != nil {
			return nil, fmt.Errorf("unable to create http client with passed in CA certificates: %w", err)
		}
		c.HttpClient = httpClient
	} else {
		c.HttpClient = new(http.Client)
	}
	c.Repositories = &Repositories{
		c:                  c,
		PullRequests:       &PullRequests{c: c},
		Pipelines:          &Pipelines{c: c},
		Repository:         &Repository{c: c},
		Issues:             &Issues{c: c},
		Commits:            &Commits{c: c},
		Diff:               &Diff{c: c},
		BranchRestrictions: &BranchRestrictions{c: c},
		Webhooks:           &Webhooks{c: c},
		Downloads:          &Downloads{c: c},
		DeployKeys:         &DeployKeys{c: c},
	}
	c.Users = &Users{
		c:       c,
		SSHKeys: &SSHKeys{c: c},
	}
	c.User = &User{c: c}
	c.Teams = &Teams{c: c}
	c.Workspaces = &Workspace{c: c, Repositories: c.Repositories, Permissions: &Permission{c: c}}
	return c, nil
}

func (c *Client) GetOAuthToken() oauth2.Token {
	return c.Auth.token
}

func (c *Client) GetApiBaseURL() string {
	return fmt.Sprintf("%s%s", c.GetApiHostnameURL(), c.apiBaseURL.Path)
}

func (c *Client) GetApiHostnameURL() string {
	return fmt.Sprintf("%s://%s", c.apiBaseURL.Scheme, c.apiBaseURL.Host)
}

func (c *Client) SetApiBaseURL(urlStr url.URL) {
	c.apiBaseURL = &urlStr
}

func (c *Client) executeRaw(method string, urlStr string, text string) (io.ReadCloser, error) {
	body := strings.NewReader(text)

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	c.authenticateRequest(req)
	return c.doRawRequest(req, false)
}

func (c *Client) execute(method string, urlStr string, text string) (interface{}, error) {
	return c.executeWithContext(method, urlStr, text, context.Background())
}

func (c *Client) executeWithContext(method string, urlStr string, text string, ctx context.Context) (interface{}, error) {
	body := strings.NewReader(text)
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if ctx != nil {
		req.WithContext(ctx)
	}
	c.authenticateRequest(req)
	result, err := c.doRequest(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) executePaginated(method string, urlStr string, text string, page *int) (interface{}, error) {
	if c.Pagelen != DEFAULT_PAGE_LENGTH {
		urlObj, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		q := urlObj.Query()
		q.Set("pagelen", strconv.Itoa(c.Pagelen))
		urlObj.RawQuery = q.Encode()
		urlStr = urlObj.String()
	}

	body := strings.NewReader(text)
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	c.authenticateRequest(req)
	result, err := c.doPaginatedRequest(req, page, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// executeFormURLEncoded posts an application/x-www-form-urlencoded body and
// discards the response. It's intended for endpoints (notably the source
// commit endpoint) that accept either multipart or URL-encoded payloads.
func (c *Client) executeFormURLEncoded(method, urlStr string, form url.Values, ctx context.Context) error {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if ctx != nil {
		req.WithContext(ctx)
	}
	c.authenticateRequest(req)
	_, err = c.doRequest(req, true)
	return err
}

func (c *Client) executeFileUpload(method string, urlStr string, files []File, filesToDelete []string, params map[string]string, ctx context.Context) (interface{}, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	for _, file := range files {
		fileReader, err := os.Open(file.Path)
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()

		if fw, err = w.CreateFormFile(file.Name, file.Name); err != nil {
			return nil, err
		}

		if _, err = io.Copy(fw, fileReader); err != nil {
			return nil, err
		}
	}

	for key, value := range params {
		err := w.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}

	for _, filename := range filesToDelete {
		err := w.WriteField("files", filename)
		if err != nil {
			return nil, err
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest(method, urlStr, &b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	if ctx != nil {
		req.WithContext(ctx)
	}
	c.authenticateRequest(req)
	return c.doRequest(req, true)

}

func (c *Client) authenticateRequest(req *http.Request) {
	if c.Auth.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.Auth.bearerToken)
	}

	if c.Auth.user != "" && c.Auth.password != "" {
		req.SetBasicAuth(c.Auth.user, c.Auth.password)
	} else if c.Auth.token.Valid() {
		c.Auth.token.SetAuthHeader(req)
	}
}

func (c *Client) doRequest(req *http.Request, emptyResponse bool) (interface{}, error) {
	resBody, err := c.doRawRequest(req, emptyResponse)
	if err != nil {
		return nil, err
	}
	if emptyResponse || resBody == nil {
		return nil, nil
	}

	defer resBody.Close()

	responseBytes, err := ioutil.ReadAll(resBody)
	if err != nil {
		return resBody, err
	}

	var result interface{}
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return responseBytes, err
	}
	return result, nil
}

func (c *Client) doPaginatedRequest(req *http.Request, page *int, emptyResponse bool) (interface{}, error) {
	disableAutoPaging := c.DisableAutoPaging
	curPage := 1
	if page != nil {
		disableAutoPaging = true
		curPage = *page
		q := req.URL.Query()
		q.Set("page", strconv.Itoa(curPage))
		req.URL.RawQuery = q.Encode()
	}
	// q.Encode() does not encode "~".
	req.URL.RawQuery = strings.ReplaceAll(req.URL.RawQuery, "~", "%7E")

	resBody, err := c.doRawRequest(req, emptyResponse)
	if err != nil {
		return nil, err
	}
	if emptyResponse || resBody == nil {
		return nil, nil
	}

	defer resBody.Close()

	responseBytes, err := ioutil.ReadAll(resBody)
	if err != nil {
		return resBody, err
	}

	responsePaginated := &Response{}
	err = json.Unmarshal(responseBytes, responsePaginated)
	if err == nil && len(responsePaginated.Values) > 0 {
		values := responsePaginated.Values
		for {
			if disableAutoPaging || responsePaginated.Next == "" ||
				(curPage >= c.LimitPages && c.LimitPages != 0) {
				break
			}
			curPage++
			newReq, err := http.NewRequest(req.Method, responsePaginated.Next, nil)
			if err != nil {
				return resBody, err
			}
			c.authenticateRequest(newReq)
			resp, err := c.doRawRequest(newReq, false)
			if err != nil {
				return resBody, err
			}

			responsePaginated = &Response{}
			json.NewDecoder(resp).Decode(responsePaginated)
			values = append(values, responsePaginated.Values...)
		}
		responsePaginated.Values = values
		responseBytes, err = json.Marshal(responsePaginated)
		if err != nil {
			return resBody, err
		}
	}

	var result interface{}
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return resBody, err
	}
	return result, nil
}

func (c *Client) doRawRequest(req *http.Request, emptyResponse bool) (io.ReadCloser, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if unexpectedHttpStatusCode(resp.StatusCode) {
		defer resp.Body.Close()

		out := &UnexpectedResponseStatusError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			out.Body = []byte(fmt.Sprintf("could not read the response body: %v", err))
		} else {
			out.Body = body
		}

		return nil, out
	}

	if emptyResponse || resp.StatusCode == http.StatusNoContent {
		resp.Body.Close()
		return nil, nil
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	return resp.Body, nil
}

func unexpectedHttpStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusAccepted:
		return false
	default:
		return true
	}
}

func (c *Client) requestUrl(template string, args ...interface{}) string {

	if len(args) == 1 && args[0] == "" {
		return c.GetApiBaseURL() + template
	}
	return c.GetApiBaseURL() + fmt.Sprintf(template, args...)
}

func (c *Client) addMaxDepthParam(params *url.Values, customMaxDepth *int) {
	maxDepth := c.MaxDepth
	if customMaxDepth != nil && *customMaxDepth > 0 {
		maxDepth = *customMaxDepth
	}

	if maxDepth != DEFAULT_MAX_DEPTH {
		params.Set("max_depth", strconv.Itoa(maxDepth))
	}
}
