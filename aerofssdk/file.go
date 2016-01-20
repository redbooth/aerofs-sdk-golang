package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	"io"
)

// File, client wrapper
type FileClient struct {
	APIClient *api.Client
	Desc      File
	OnDemand  []string
}

// File descriptor
type File api.File

// Construct a FileClient given a file identifier and APIClient
func NewFileClient(c *api.Client, fileId string, fields []string) (*FileClient, error) {
	body, header, err := c.GetFileMetadata(fileId, fields)
	if err != nil {
		return nil, err
	}

	f := FileClient{APIClient: c, OnDemand: fields}
	err = json.Unmarshal(body, &f.Desc)

	if err != nil {
		return nil, errors.New("Unable to unmarshal existing File")
	}
	f.Desc.Etag = header.Get("ETag")
	return &f, nil
}

// Reload the ParentPath of the File
func (f *FileClient) LoadPath() error {
	body, header, err := f.APIClient.GetFilePath(f.Desc.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &f.Desc.Path)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved file ParentPath")
	}

	f.Desc.Etag = header.Get("ETag")
	return nil
}

// Move the file to a new parent folder
func (f *FileClient) Move(newName, parentId string) error {
	body, header, err := f.APIClient.MoveFile(f.Desc.Id, parentId, newName,
		[]string{f.Desc.Etag})
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &f.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal the new File location")
	}
	f.Desc.Etag = header.Get("Etag")
	return nil
}

// Retrieve the file contents
func (f *FileClient) GetContent() ([]byte, error) {
	body, header, err := f.APIClient.GetFileContent(f.Desc.Id, f.Desc.Etag, 0,
		f.Desc.Size-1, []string{})
	if err != nil {
		return nil, err
	}

	f.Desc.Etag = header.Get("ETag")
	return body, nil
}

// Update the existing content of a file
func (f *FileClient) UploadFile(file io.Reader) error {
	uploadId, err := f.APIClient.GetFileUploadId(f.Desc.Id, []string{f.Desc.Etag})
	if err != nil {
		return errors.New("Unable to retrieve UploadId for file")
	}

	return f.APIClient.UploadFile(f.Desc.Id, uploadId, file,
		[]string{f.Desc.Etag})
}
