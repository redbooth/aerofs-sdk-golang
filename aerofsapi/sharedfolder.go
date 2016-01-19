package aerofsapi

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

const (
	SF_ROUTE = "shares"
)

func (c *Client) ListSharedFolders(email string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "shares"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) ListSharedFolderMetadata(sid string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, sid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}
	return unpackageResponse(res)
}

func (c *Client) CreateSharedFolder(name string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE}, "/")
	link := c.getURL(route, "")
	data := []byte(fmt.Sprintf(`{"name" : %s"}`, name))

	res, err := c.post(link, bytes.NewBuffer(data))

	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}
