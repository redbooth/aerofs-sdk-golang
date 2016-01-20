package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	INVITEE_ROUTE = "invitees"
)

func (c *Client) GetInvitee(email string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{INVITEE_ROUTE, email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) CreateInvitee(email_to, email_from string) ([]byte,
	*http.Header, error) {
	link := c.getURL(INVITEE_ROUTE, "")
	invitee := map[string]string{
		"email_to":   email_to,
		"email_from": email_from,
	}

	data, err := json.Marshal(invitee)
	if err != nil {
		return nil, nil, errors.New("Unable to serialize invitation request")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	route := strings.Join([]string{INVITEE_ROUTE, email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	defer res.Body.Close()
	_, _, err = unpackageResponse(res)
	return err
}
