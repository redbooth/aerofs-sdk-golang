package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	FILE_ROUTE = "files"
)

func (c *Client) GetFileMetadata(fileId string, fields []string) ([]byte,
	*http.Header, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId}, "/")
	query := url.Values{"fields": fields}
	link := c.getURL(route, query.Encode())

	newHeader := http.Header{}
	newHeader.Set("Content-Type", "application/octet-stream")

	res, err := c.request("GET", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetFilePath(fileId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId, "path"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

//
func (c *Client) GetFileContent(fileId, rangeEtag string, startIndex, endIndex int, matchEtags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId, "content"}, "/")
	link := c.getURL(route, "")

	// Construct header
	byteRange := fmt.Sprintf("bytes=%d-%d", startIndex, endIndex)
	newHeader := http.Header{}
	newHeader.Add("Range", byteRange)

	if rangeEtag != "" {
		newHeader.Set("If-Range", rangeEtag)
	}

	if len(matchEtags) > 0 {
		for _, v := range matchEtags {
			newHeader.Add("If-None-Match", v)
		}
	}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Instantiate a newfile
func (c *Client) CreateFile(parentId, fileName string) ([]byte, *http.Header,
	error) {
	link := c.getURL(FILE_ROUTE, "")

	newFile := map[string]string{
		"parent": parentId,
		"name":   fileName,
	}
	data, err := json.Marshal(newFile)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Functions that upload content to a file

// Retrieve a files UploadID to be used for future content uploads
// Upload Identifiers are only valid for ~24 hours
func (c *Client) GetFileUploadId(fileId string, etags []string) (string, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	newHeader.Set("Content-Range", "bytes */*")
	newHeader.Set("Content-Length", "0")
	for _, v := range etags {
		newHeader.Add("If-Match", v)
	}

	res, err := c.request("PUT", link, &newHeader, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	_, h, err := unpackageResponse(res)
	if err != nil {
		return "", err
	}

	fmt.Println(h)
	return h.Get("Upload-ID"), nil
}

// Retrieve the list of bytes already transferred by an unfinished upload
func (c *Client) GetUploadBytesSize(fileId, uploadId string, etags []string) (int, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	newHeader.Set("Content-Range", "bytes /*/")
	newHeader.Set("Upload-ID", uploadId)
	newHeader.Set("Content-Length", "0")
	if len(etags) > 0 {
		for _, v := range etags {
			newHeader.Add("If-Match", v)
		}
	}

	res, err := c.request("PUT", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return 0, nil
	}

	bytesUploaded, err := strconv.Atoi(res.Header.Get("Range"))
	if err != nil {
		return 0, errors.New("Unable to parse value of bytes transferred from HTTP-Header")
	}
	return bytesUploaded, err
}

// Upload a single file chunk
func (c *Client) UploadFileChunk(fileId, uploadId string, chunks []byte, startIndex, lastIndex int) (*http.Header, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId, "content"}, "/")
	link := c.getURL(route, "")
	byteRange := fmt.Sprintf("bytes %d-%d/*", startIndex, lastIndex)
	newHeader := http.Header{}
	newHeader.Set("Content-Range", byteRange)
	newHeader.Set("Upload-ID", uploadId)
	newHeader.Set("Content-Length", "0")

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(chunks))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	_, h, e := unpackageResponse(res)
	return h, e
}

// Upload a file
func (c *Client) UploadFile(fileId, uploadId string, file io.Reader, etags []string) error {
	route := strings.Join([]string{FILE_ROUTE, fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{
		"If-Match":  etags,
		"Upload-ID": []string{uploadId},
	}

	err := c.uploadFileChunks(link, &newHeader, file)
	return err
}

// Helper function to upload sequential chunks of a file
func (c *Client) uploadFileChunks(link string, header *http.Header, file io.Reader) error {
	// Indices for file byte-ranges
	startIndex := 0
	endIndex := 0

	chunk := make([]byte, CHUNKSIZE)

	// Iterate over the file in CHUNKSIZE pieces until EOF occurs
	// Do not set "Instance-Length" of "Content-Range" until we hit EOF and use
	// the format "bytes <startIndex>-<endIndex>/* for intermediary uploads

	for {
		size, fileErr := file.Read(chunk)
		endIndex += size - 1

		switch {
		// If we have read all chunks, set Content-Range instance-Length
		case fileErr == io.EOF:
			byteRange := fmt.Sprintf("bytes %d-%d/%d", startIndex, endIndex, endIndex)
			header.Set("Content-Range", byteRange)
		case fileErr != nil:
			return fileErr
		case fileErr == nil:
			byteRange := fmt.Sprintf("bytes %d-%d/*", startIndex, endIndex)
			header.Set("Content-Range", byteRange)
		}

		res, httpErr := c.request("PUT", link, header, bytes.NewBuffer(chunk))
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		if fileErr == io.EOF {
			break
		}

		startIndex += size
	}
	return nil
}
func (c *Client) MoveFile(fileId, parentId, name string, etags []string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{FILE_ROUTE, fileId}, "/")
	link := c.getURL(route, "")

	newHeader := http.Header{}
	if len(etags) > 0 {
		for _, v := range etags {
			newHeader.Add("If-Match", v)
		}
	}

	newFile := map[string]string{
		"parent": parentId,
		"name":   name,
	}

	data, err := json.Marshal(newFile)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	return unpackageResponse(res)
}

// There must be at least one etag present
func (c *Client) DeleteFile(fileid string, etags []string) error {
	if len(etags) == 0 {
		return errors.New("At least 1 ETag must be present when deleting a file")
	}

	route := strings.Join([]string{FILE_ROUTE, fileid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-Match": etags}

	res, err := c.request("DEL", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	_, _, err = unpackageResponse(res)
	return err
}
