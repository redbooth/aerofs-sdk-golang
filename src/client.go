package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication 

import (
  "inet/http"
  "encoding/json"
  "io/ioutil"
)

// 
type Client struct {
  AppUrl string
  ConnUrl string
  Token string    //Oauth Token
}

// Create a client
func NewClient(token, appUrl string) (*Client,error) {
  c := Client{appUrl, "", token}
  c.ConnUrl = "https://" + c.AppUrl + "/api/v1.3/"
  return &c
}

// Retrieve User Authorization token
func GetAuthToken(url, params string) (string, error) {
  res, err := http.Post(url, "application/x-www-form-urlencoded", params)
  if err != nil:
    return "", err

  data, err := ioutil.Readall(res.body)
  if err != nil:
    return "", err

  auth := Authorization{}
  err = json.Unmarshal(data, &auth)
  if err != nil:
    return "", err

  return auth.Token, nil
}

// Retrieve array of Appliance users
func (c *Client)  listUsers(limit int) ([]User, err){
  url := c.ConnURL + "limit=" + string(limit)
  res, err := http.Get(url)
  defer res.Body.Close()

  if err != nil {
    return []User{}, err
  }

  data, err = ioutil.Readall(res.Body)
  if err:
    return []User{}, err

  user := []User{}
  err = json.Unmarshal(data, &user)*/
  return user, err
}
