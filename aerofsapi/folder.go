package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	FOLDER_ROUTE = "folders"
)

func (c *Client) GetFolderMetadata(folderId string, fields []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FOLDER_ROUTE, folderId}, "/")
	query := ""
	if len(fields) > 0 {
		v := url.Values{"fields": fields}
		query = v.Encode()
	}
	link := c.getURL(route, query)

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetFolderPath(folderId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FOLDER_ROUTE, folderId, "path"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetFolderChildren(folderId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FOLDER_ROUTE, folderId, "children"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	return unpackageResponse(res)
}

func (c *Client) CreateFolder(parentId, name string) ([]byte, *http.Header, error) {
	link := c.getURL(FOLDER_ROUTE, "")

	newFolder := map[string]string{
		"parent": parentId,
		"name":   name,
	}
	data, err := json.Marshal(newFolder)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal JSON for new folder")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Move a folder given its existing unique ID, the ID of its new parent and its
// new folder Name
func (c *Client) MoveFolder(folderId, newParentId, newFolderName string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FOLDER_ROUTE, folderId}, "/")
	link := c.getURL(route, "")

	content := map[string]string{"parent": newParentId, "name": newFolderName}
	data, err := json.Marshal(content)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal JSON for moving folder")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) DeleteFolder(folderId string, etags []string) error {
	route := strings.Join([]string{FOLDER_ROUTE, folderId}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DELETE", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		_, _, err = unpackageResponse(res)
	}

	return err
}

func (c *Client) ShareFolder(folderId string) error {
	route := strings.Join([]string{FOLDER_ROUTE, folderId, "is_shared"}, "/")
	link := c.getURL(route, "")

	res, err := c.put(link, nil)
	defer res.Body.Close()
	if err != nil {
		_, _, err = unpackageResponse(res)
	}

	return err
}
