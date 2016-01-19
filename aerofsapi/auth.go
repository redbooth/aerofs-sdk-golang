package aerofsapi

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Interface between an application and the AeroFS server
// Can be constructed from a generated AeroFS Config Json file
// Note that the Redirect URL must be the same as that stored when registering
// the third party application

type AuthClient struct {
	// The URL of an AeroFS Appliance instance
	AeroUrl string `json:"hostname"`

	// 3rd-Party specific
	Id     string `json:"client_id"`
	Secret string `json:"client_secret"`

	// The URL the user will be redirected to after accepting scope permissions
	Redirect string `json:"redirect"`

	// 3rd-Party specific Permission scopes
	Scopes []string

	// A unique identifier created by the 3rd-Party App and eventually passed back
	// to it after the user has confirmed scopes on the AeroFS Appliance end
	State string
}

// The response when receiving a token given an authorization code
type AccessResponse struct {
	// OAuth 2.0 Token
	Token string `json:"access_token"`

	// Default is "Bearer"
	TokenType string `json:"token_type"`

	// Number of seconds until token expiration
	ExpireTime int `json:"expires_in"`

	// User scopes associated with this token
	Scopes string `json:"scope"`
}

// Create a new AuthClient from an AeroFS appconfig.json file
// The appconfig.json only contains the AeroURL and client_{id,secret}
func NewAuthClient(fileName, redirectUri, state string, scopes []string) (*AuthClient, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Unable to read the defined file")
	}

	// Unmarshalling will not reset these fields
	authClient := AuthClient{
		Redirect: redirectUri,
		State:    state,
		Scopes:   scopes,
	}
	err = json.Unmarshal(data, &authClient)
	if err != nil {
		return nil, errors.New("Unable to unmarshal the requested appconfig.json")
	}

	return &authClient, nil
}

// Return a URL to the AeroFS Appliance so the user can authorize third-party
// access
func (auth *AuthClient) GetAuthorizationUrl() string {
	scopes := strings.Join(auth.Scopes, ",")
	v := make(url.Values)
	v.Set("response_type", "code")
	v.Set("client_id", auth.Id)
	v.Set("redirect_uri", auth.Redirect)
	v.Set("scope", scopes)
	if auth.State != "" {
		v.Set("state", auth.State)
	}

	route := "authorize"
	link := url.URL{
		Scheme:   "https",
		Host:     auth.AeroUrl,
		Path:     route,
		RawQuery: v.Encode(),
	}

	return link.String()
}

// Retrieve User OAuth token, granted scopes given an Authorization code
func (auth *AuthClient) GetAccessToken(code string) (string, []string, error) {
	v := make(url.Values)
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("client_id", auth.Id)
	v.Set("client_secret", auth.Secret)
	v.Set("redirect_uri", auth.Redirect)

	link := url.URL{Scheme: "https",
		Host: auth.AeroUrl,
		Path: strings.Join([]string{"auth", "token"}, "/"),
	}
	body := bytes.NewBuffer([]byte(v.Encode()))
	encoding := "application/x-www-form-urlencoded"

	res, err := http.Post(link.String(), encoding, body)
	defer res.Body.Close()
	if err != nil {
		return "", []string{}, err
	}

	accessResponse := AccessResponse{}
	err = GetEntity(res, &accessResponse)
	grantedScopes := strings.Split(accessResponse.Scopes, ",")
	return accessResponse.Token, grantedScopes, err
}
