package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/schema"
	"github.com/pghq/go-tea"
)

const (
	// Version of the client
	Version = "v2"

	// defaultBaseURL for sending DeepL Pro API requests
	defaultBaseURL = "https://api-free.deepl.com"

	// userAgent for DeepL Pro API requests
	userAgent = "go-ink/v0"
)

// Client for the DeepL Pro API
type Client struct {
	http          *http.Client
	baseURL       *url.URL
	encoder       *schema.Encoder
	authorization string
}

// NewClient creates a new client instance for the DeepL Pro API
func NewClient(authKey string, opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		http:          http.DefaultClient,
		encoder:       schema.NewEncoder(),
		baseURL:       baseURL,
		authorization: fmt.Sprintf("DeepL-Auth-Key %s", authKey),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// do a http request to DeepL API
func (c Client) do(req *http.Request, v interface{}) error {
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx := req.Context(); ctx != nil && ctx.Err() != nil {
			return tea.Stack(ctx.Err())
		}

		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var res struct {
			Message string `json:"message"`
		}
		res.Message = errorMessage(resp.StatusCode)
		_ = json.NewDecoder(resp.Body).Decode(&res)
		return tea.AsErrTransfer(resp.StatusCode, tea.Err(res.Message))
	}

	defer resp.Body.Close()
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
			return tea.Stack(err)
		}
	}

	return nil
}

// newRequest creates a http request
func (c Client) newRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(Version, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", c.authorization)

	return req, nil
}

// ClientOption for custom client configuration
type ClientOption func(c *Client)

// BaseURL sets a custom base URL for the client
func BaseURL(o string) ClientOption {
	return func(c *Client) {
		c.baseURL, _ = url.Parse(o)
	}
}

// HttpClient sets a custom http client for the client
func HttpClient(o *http.Client) ClientOption {
	return func(c *Client) {
		c.http = o
	}
}

// errorMessage obtains an error message from the status code
func errorMessage(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Bad request. Please check error message and your parameters."
	case http.StatusForbidden:
		return "Authorization failed. Please supply a valid auth_key parameter."
	case http.StatusNotFound:
		return "The requested resource could not be found."
	case http.StatusRequestEntityTooLarge:
		return "The request size exceeds the limit."
	case http.StatusTooManyRequests:
		return "Too many requests. Please wait and resend your request."
	case 456:
		return "Quota exceeded. The character limit has been reached."
	case http.StatusServiceUnavailable:
		return "Resource currently unavailable. Try again later."
	case 529:
		return "Too many requests. Please wait and resend your request."
	default:
		return http.StatusText(statusCode)
	}
}
