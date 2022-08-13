package client

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/tim3-p/gophkeeper/internal/store"
)

const (
	defaultClientTimeout = time.Second * 1
)

// Client describes general client configuration
type Client struct {
	ServerAddr    string
	UserName      string
	UserPass      string
	CacheFile     string
	Store         *store.Store
	Timeout       time.Duration
	HTTPSInsecure bool
}

// NewClient returns new client
func NewClient(serverAddr string,
	userName string,
	userPass string,
	cacheFile string,
	httpsInsecure bool,
) *Client {
	var s *store.Store
	var err error
	if cacheFile != "" {
		s, err = store.NewStore(cacheFile)
		if err != nil {
			log.Print(err)
			return nil
		}
	}
	return &Client{
		ServerAddr:    serverAddr,
		UserName:      userName,
		UserPass:      userPass,
		CacheFile:     cacheFile,
		Timeout:       defaultClientTimeout,
		HTTPSInsecure: httpsInsecure,
		Store:         s,
	}
}

func (c *Client) prepaReq(method, path string, body []byte) (*http.Request, error) {
	b := bytes.NewReader(body)
	req, err := http.NewRequest(method, c.ServerAddr+path, b)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(c.UserName, c.UserPass)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.HTTPSInsecure,
		},
	}
	client := &http.Client{
		Timeout:   c.Timeout,
		Transport: tr,
	}
	return client
}
