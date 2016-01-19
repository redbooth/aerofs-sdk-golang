package aerofsapi

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	GROUP_ROUTE = "groups"
)

func (c *Client) ListGroups(offset, results int) ([]byte, *http.Header, error) {
	query := url.Values{}
	query.Set("offset", strconv.Itoa(offset))
	query.Set("results", strconv.Itoa(results))
	link := c.getURL(GROUP_ROUTE, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) CreateGroup(groupName string) ([]byte, *http.Header, error) {
	link := c.getURL(GROUP_ROUTE, "")
	// TODO : Is this preferred to constructing a map, then marshalling?
	// robust vs. bootstrap
	newGroup := []byte(fmt.Sprintf(`{"name" : %s}`, groupName))

	res, err := c.post(link, bytes.NewBuffer(newGroup))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetGroup(groupId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"request", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) DeleteGroup(groupId string) error {
	route := strings.Join([]string{GROUP_ROUTE, groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
