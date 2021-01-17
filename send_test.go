package gomsteams

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	opts := Options{
		Timeout: 60 * time.Second,
		Verbose: true,
	}
	client := NewClient(opts)
	assert.IsType(t, &Client{}, client)
}

func TestTeamsClientSend(t *testing.T) {
	simpleMsgCard := NewMessageCard()
	simpleMsgCard.Text = "Hello World"
	webHookUrl := "https://outlook.office.com/webhook/xxx"
	var tests = []struct {
		reqURL    string
		reqMsg    MessageCard
		resStatus int   // httpClient response status
		resError  error // httpClient error
		error     error // method error
	}{
		// invalid webhookURL - url.Parse error
		{
			reqURL:    "http://",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resError:  nil,
			error:     &url.Error{},
		},
		// invalid webhookURL - missing prefix in webhook URL
		{
			reqURL:    "",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resError:  nil,
			error:     &url.Error{},
		},
		// invalid httpClient.Do call
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  errors.New("pling"),
			error:     errors.New(""),
		},
		// invalid httpClient.Do call
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  errors.New("pling"),
			error:     errors.New(""),
		},
		// invalid response status code
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 400,
			resError:  nil,
			error:     errors.New(""),
		},
		// invalid response status code
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 400,
			resError:  nil,
			error:     errors.New(""),
		},
		// valid
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  nil,
			error:     nil,
		},
		// valid
		{
			reqURL:    webHookUrl,
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  nil,
			error:     nil,
		},
	}
	for _, test := range tests {
		client := NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), test.reqURL)
			return &http.Response{
				StatusCode: test.resStatus,
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}, test.resError
		})

		opts := Options{
			Timeout: 60 * time.Second,
			Verbose: true,
		}
		c := &Client{httpClient: client, options: &opts}

		err := c.Send(test.reqURL, test.reqMsg)
		assert.IsType(t, test.error, err)
	}
}

// helper for testing --------------------------------------------------------------------------------------------------

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient returns *http.API with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
