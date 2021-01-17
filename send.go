package gomsteams

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Known webhook URL prefixes for submitting messages to Microsoft Teams
const (
	WebhookURLOfficeComPrefix = "https://outlook.office.com"
	WebhookURLOffice365Prefix = "https://outlook.office365.com"
)

var (
	// ErrUserAccessDenied access denied error
	ErrUserAccessDenied = errors.New("you do not have access to the requested resource")
	// ErrNotFound 404 error
	ErrNotFound = errors.New("the requested resource not found")
	// ErrTooManyRequests error too many requests
	ErrTooManyRequests = errors.New("you have exceeded throttle")
)

// API - interface of MS Teams notify
type API interface {
	Send(webhookURL string, webhookMessage MessageCard) error
}

// Options - options for the API httpClient
type Options struct {
	Timeout time.Duration
	Verbose bool
}

// Client MS teams Http client
type Client struct {
	httpClient *http.Client
	options    *Options
}

// NewClient create a brand new client for MS Teams notify
func NewClient(options Options) *Client {
	if options.Timeout.String() == "" {
		options.Timeout = 30 * time.Second
	}

	teamsClient := &Client{
		httpClient: &http.Client{
			Timeout: options.Timeout,
		},
		options: &options,
	}

	return teamsClient
}

func (c *Client) newRequest(ctx context.Context, method, reqURL string, payload interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	if c.options.Verbose {
		body, _ := httputil.DumpRequest(req, true)
		log.Println(string(body))
	}

	req = req.WithContext(ctx)
	return req, nil
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to make request [%s:%s]: %v", r.Method, r.URL.String(), err)
	}

	if c.options.Verbose {
		body, _ := httputil.DumpResponse(resp, true)
		log.Println(string(body))
	}

	switch resp.StatusCode {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent:
		return resp, nil
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusUnauthorized,
		http.StatusForbidden:
		return nil, ErrUserAccessDenied
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	}

	return nil, fmt.Errorf("failed to do request, %d status code received", resp.StatusCode)
}

func (c *Client) doRequest(r *http.Request, v interface{}) error {
	resp, err := c.do(r)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println("Error: error in closing response body, ", err)
		}
	}()

	if v == nil {
		return nil
	}

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("could not parse response body: %v [%s:%s] %s", err, r.Method, r.URL.String(), buf.String())
	}

	return nil
}

// Send - will post a notification to MS Teams webhook URL
func (c *Client) Send(webhookURL string, webhookMessage MessageCard) error {
	// Validate input data
	if valid, err := IsValidInput(webhookMessage, webhookURL); !valid {
		return err
	}

	// make new request
	req, err := c.newRequest(context.Background(), http.MethodPost, webhookURL, webhookMessage)
	if err != nil {
		return fmt.Errorf("error in creating request, %s", err)
	}

	// do the request
	err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// helper --------------------------------------------------------------------------------------------------------------

// IsValidInput is a validation "wrapper" function. This function is intended
// to run current validation checks and offer easy extensibility for future
// validation requirements.
func IsValidInput(webhookMessage MessageCard, webhookURL string) (bool, error) {
	// validate url
	if valid, err := IsValidWebhookURL(webhookURL); !valid {
		return false, err
	}

	// validate message
	if valid, err := IsValidMessageCard(webhookMessage); !valid {
		return false, err
	}

	return true, nil
}

// IsValidWebhookURL performs validation checks on the webhook URL used to submit messages to Microsoft Teams.
func IsValidWebhookURL(webhookURL string) (bool, error) {
	switch {
	case strings.HasPrefix(webhookURL, WebhookURLOfficeComPrefix):
	case strings.HasPrefix(webhookURL, WebhookURLOffice365Prefix):
	default:
		u, err := url.Parse(webhookURL)
		if err != nil {
			return false, fmt.Errorf(
				"unable to parse webhook URL %q: %v",
				webhookURL,
				err,
			)
		}
		userProvidedWebhookURLPrefix := u.Scheme + "://" + u.Host

		return false, &url.Error{Err: fmt.Errorf("webhook URL does not contain expected prefix; got %q, expected one of %q or %q",
			userProvidedWebhookURLPrefix,
			WebhookURLOfficeComPrefix,
			WebhookURLOffice365Prefix),
		}
	}

	return true, nil
}

// IsValidMessageCard performs validation/checks for known issues with
// MessageCard values.
func IsValidMessageCard(webhookMessage MessageCard) (bool, error) {
	if (webhookMessage.Text == "") && (webhookMessage.Summary == "") {
		// This scenario results in:
		// 400 Bad Request
		// Summary or Text is required.
		return false, fmt.Errorf("invalid message card: summary or text field is required")
	}

	return true, nil
}
