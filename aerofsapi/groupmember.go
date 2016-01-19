package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	GROUPMEMBER_ROUTE = "groups"
)

func (c *Client) ListGroupMembers(groupId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{GROUPMEMBER_ROUTE, groupId, "members"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) AddGroupMember(groupId, name string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{GROUPMEMBER_ROUTE, groupId, "members"}, "/")
	link := c.getURL(route, "")
	newMember := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal provided group member")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetGroupMember(groupId, email string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{GROUPMEMBER_ROUTE, groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) RemoveMember(groupId, email string) error {
	route := strings.Join([]string{GROUPMEMBER_ROUTE, groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	_, _, err = unpackageResponse(res)
	return err
}
