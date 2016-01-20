package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	USERS_ROUTE = "users"
)

func (c *Client) ListUsers(limit int, after, before *string) ([]byte, *http.Header, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(limit))
	if before != nil {
		query.Set("before", *before)
	}
	if after != nil {
		query.Set("after", *after)
	}

	link := c.getURL(USERS_ROUTE, query.Encode())
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetUser(email string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{USERS_ROUTE, email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) CreateUser(email, firstName, lastName string) ([]byte,
	*http.Header, error) {
	link := c.getURL(USERS_ROUTE, "")

	user := map[string]string{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
	}
	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	return unpackageResponse(res)
}

func (c *Client) UpdateUser(email, firstName, lastName string) ([]byte,
	*http.Header, error) {
	route := strings.Join([]string{USERS_ROUTE, email}, "/")
	link := c.getURL(route, "")

	user := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	defer res.Body.Close()

	return unpackageResponse(res)
}

func (c *Client) DeleteUser(email string) error {
	route := strings.Join([]string{USERS_ROUTE, email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	res.Body.Close()

	return err
}

func (c *Client) ChangePassword(email, password string) error {
	route := strings.Join([]string{USERS_ROUTE, email, "password"}, "/")
	link := c.getURL(route, "")
	data := []byte(`"` + password + `"`)

	_, err := c.put(link, bytes.NewBuffer(data))
	return err
}

func (c *Client) DisablePassword(email string) error {
	route := strings.Join([]string{USERS_ROUTE, email, "password"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

func (c *Client) CheckTwoFactorAuth(email string) error {
	route := strings.Join([]string{USERS_ROUTE, email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.get(link)
	return err
}

func (c *Client) DisableTwoFactorAuth(email string) error {
	route := strings.Join([]string{USERS_ROUTE, email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}
