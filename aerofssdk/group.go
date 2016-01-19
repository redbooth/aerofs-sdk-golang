package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

// Group, client wrapper
type GroupClient struct {
	APIClient *api.Client
	Desc      Group
}

// Group descriptor
type Group api.Group

// List all groups
func ListGroups(c *api.Client, offset, results int) (*[]Group, error) {
	body, _, err := c.ListGroups(offset, results)
	if err != nil {
		return nil, err
	}
	groups := []Group{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, errors.New("Unable to unmarshal list of groups")
	}

	return &groups, nil
}

// Retrieve an existing group
func NewGroupClient(c *api.Client, groupId string) (*GroupClient, error) {
	body, _, err := c.GetGroup(groupId)
	if err != nil {
		return nil, err
	}

	g := GroupClient{APIClient: c}
	err = json.Unmarshal(body, &g.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal existing group")
	}

	return &g, nil
}

// Create a group
func CreateGroupClient(c *api.Client, groupName string) (*GroupClient, error) {
	body, _, err := c.CreateGroup(groupName)
	if err != nil {
		return nil, err
	}

	g := GroupClient{APIClient: c}
	err = json.Unmarshal(body, &g.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal created group")
	}

	return &g, nil
}

// Update a group client
func (g *GroupClient) Load() error {
	body, _, err := g.APIClient.GetGroup(g.Desc.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &g.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal existing group")
	}

	return nil
}

// Delete the group
func (g *GroupClient) Delete() error {
	return g.APIClient.DeleteGroup(g.Desc.Id)
}

// Add a group member to the group
func (g *GroupClient) AddGroupMember(email string) error {
	_, _, err := g.APIClient.AddGroupMember(g.Desc.Id, email)
	if err != nil {
		return nil
	}

	g.Load()
	return nil
}
