package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (c *Client) ListSFMembers(id string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, id, "members"}, "/")
	newHeader := http.Header{}
	if len(etags) > 0 {
		newHeader = http.Header{"If-None-Match": etags}
	}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetSFMember(id, email string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, id, "members", email}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	if len(etags) > 0 {
		newHeader = http.Header{"If-None-Match": etags}
	}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) AddSFMember(id, email string, permissions []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, id, "members"}, "/")
	link := c.getURL(route, "")

	newMember := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal new ShareFolder member")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) SetSFMemberPermissions(id, email string, permissions, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	newPerms := map[string]interface{}{
		"permissions": permissions,
	}
	data, err := json.Marshal(newPerms)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) RemoveSFMember(id, email string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{SF_ROUTE, id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DEL", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}
