package aerofsapi

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// Default size used when uploading a file
	CHUNKSIZE = 1000000

	API = "api/v1.3"
)

// A Client is used to communicate with an AeroFS Appliance
type Client struct {
	// The hostname/IP of the AeroFS Appliance
	// Used when constructing the default API Prefix for all subsequent API calls
	// Ie. share.syncfs.com
	Host string

	// The OAuth token
	Token string

	// Default header containing Token, Content-type and Endpoint-Consistency
	// For conditional file, and folder requests, the header is populated
	// with an ETag
	Header http.Header

	// Stored http-connection to prevent multile TLS, TCP handshakes
	hClient http.Client
}

// API-Client Constructor
func NewClient(token, host string) (*Client, error) {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("Endpoint-Consistency", "strict")

	c := Client{Host: host,
		Header: header,
		Token:  token}

	return &c, nil
}

// Extract the response body and header
func unpackageResponse(res *http.Response) ([]byte, *http.Header, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, errors.New("Unable to read body of HTTP response")
	}
	header := res.Header

	// For each API call, unpackage the HTTP response and return an error if a non
	// 2XX status code is retrieved
	if res.StatusCode >= 300 {
		err := errors.New(res.Status)
		return body, &header, err
	}
	return body, &header, nil
}

// Construct a URL given a route and query parameters
func (c *Client) getURL(route, query string) string {
	link := url.URL{Scheme: "https",
		Path: strings.Join([]string{API, route}, "/"),
		Host: c.Host,
	}

	if query != "" {
		link.RawQuery = query
	}

	return link.String()
}

// Resets the token for a given client
// Allows the third-party developer to construct 1 SDK-Client used to retrieve
// the values for multiple users
func (c *Client) SetToken(token string) {
	c.Header.Set("Authorization", "Bearer "+token)
}

//
// Wrappers for basic HTTP functions
//

// HTTP-GET
func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP GET Request")
	}

	request.Header = c.Header
	return c.hClient.Do(request)
}

// HTTP-POST
func (c *Client) post(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP POST request")
	}

	request.Header = c.Header
	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	return c.hClient.Do(request)
}

// HTTP-PUT
func (c *Client) put(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP PUT request")
	}

	request.Header = c.Header

	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	return c.hClient.Do(request)
}

// HTTP-DELETE
func (c *Client) del(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP DELETE Request")
	}

	request.Header = c.Header
	return c.hClient.Do(request)
}

// Generic Handler for HTTP request
// Allows the passing of additional HTTP request header K/V pairs
func (c *Client) request(req, url string, options *http.Header, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(req, url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP " + req + " Request")
	}

	// If header map passed in , add additional KV pairs
	request.Header = c.Header
	if options != nil && len(*options) > 0 {
		for k, v := range *options {
			for _, el := range v {
				request.Header.Add(k, el)
			}
		}
	}

	// If we are not sending data, delete the default content-Type
	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	// TODO : Add extra field to signal serializing
	// Note : Determine if this has actual effect
	contentType := options.Get("Content-Type")
	if options.Get("Content-Type") != "" {
		request.Header.Set("Content-Type", contentType)
	}

	return c.hClient.Do(request)
}

// Unmarshalls data from an HTTP Response into a given entity
func GetEntity(res *http.Response, entity interface{}) error {
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &entity)
}
