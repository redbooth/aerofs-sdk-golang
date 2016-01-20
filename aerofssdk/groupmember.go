package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

// GroupMember, client wrapper
type GroupMemberClient struct {
	APIClient *api.Client
	Desc      GroupMember
}

// GroupMember descriptor
type GroupMember api.GroupMember

func ListGroupMembers(c *api.Client, groupId string) ([]GroupMember, error) {
	var groupMembers []GroupMember
	body, _, err := c.ListGroupMembers(groupId)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &groupMembers)
	if err != nil {
		return nil, errors.New("Unable to unmarshal list of group members")
	}

	return groupMembers, nil
}

func NewGroupMember(c *api.Client, groupId, memberEmail string) (*GroupMemberClient, error) {
	body, _, err := c.GetGroupMember(groupId, memberEmail)
	if err != nil {
		return nil, err
	}

	g := GroupMemberClient{APIClient: c, Desc: GroupMember{GroupId: groupId}}
	err = json.Unmarshal(body, &g.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal group member")
	}

	return &g, nil
}

// Update the groupMember information
func (g *GroupMemberClient) Load() error {
	body, _, err := g.APIClient.GetGroupMember(g.Desc.GroupId, g.Desc.Email)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &g.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal group member")
	}

	return nil
}
