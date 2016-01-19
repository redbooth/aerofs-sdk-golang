package aerofsapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) ListSFInvitations(email string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) ViewPendingSFInvitation(email, sid string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) AcceptSFInvitation(email, sid string, external int) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	query := url.Values{}
	query.Set("external", strconv.Itoa(external))
	link := c.getURL(route, query.Encode())

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Ignore an existing invitation to a shared folder
func (c *Client) IgnoreSFInvitation(email, sid string) error {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
