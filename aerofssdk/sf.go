package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

// SharedFolder, client wrapper
type SharedFolderClient struct {
	APIClient *api.Client
	Desc      SharedFolder
	Etag      string
}

// SharedFolder descriptors
type SharedFolder api.SharedFolder

// Retrieve a list of SharedFolder member descriptors
// TODO : Should an Etag be return for each one?
func ListSharedFolders(c *api.Client, sid string, etags []string) ([]SharedFolder, error) {
	body, _, err := c.ListSharedFolders(sid, etags)
	if err != nil {
		return nil, err
	}
	sfs := []SharedFolder{}
	err = json.Unmarshal(body, &sfs)
	if err != nil {
		return nil, errors.New("Unable to demarshal the list of retrieved SharedFolders")
	}
	return sfs, nil
}

// Retrieve an existing shared folder
func GetSharedFolderClient(c *api.Client, sid string, etags []string) (*SharedFolderClient, error) {
	body, header, err := c.ListSharedFolderMetadata(sid, etags)
	if err != nil {
		return nil, err
	}
	sfClient := SharedFolderClient{APIClient: c}
	err = json.Unmarshal(body, &sfClient.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal retrieved Shared Folder")
	}
	sfClient.Etag = header.Get("ETag")
	return &sfClient, nil
}

// Create a new shared folder and return a client associated with it
func CreateSharedFolderClient(c *api.Client, name string) (*SharedFolderClient, error) {
	body, _, err := c.CreateSharedFolder(name)
	if err != nil {
		return nil, err
	}

	sfClient := SharedFolderClient{APIClient: c}
	err = json.Unmarshal(body, &sfClient.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal retrieved SharedFolder")
	}

	return &sfClient, nil
}

// Synchronize the shared folder fields with the backend
func (sfClient *SharedFolderClient) load() error {
	body, header, err := sfClient.APIClient.ListSharedFolderMetadata(sfClient.Desc.Id, []string{sfClient.Etag})
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &sfClient.Desc)
	if err != nil {
		return errors.New("Unable to demarshal the retrieved SharedFolder")
	}
	sfClient.Etag = header.Get("ETag")
	return nil
}
