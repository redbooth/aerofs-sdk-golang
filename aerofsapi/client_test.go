package aerofsapi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

// These end-to-end tests run against a local AeroFS Test Appliance instance.
// To execute these tests, tokens for OAuth 2.0 authentication must be provided
// This can be done manually, by creating a 3rd-Party Application, and using the
// AuthClient to generate corresponding tokens. These constants are exported for
// the SDK tests

var UserToken string
var AdminToken string
var AppHost string

// Teardown Functions

// Remove all users
func removeUsers() error {
	c, err := NewClient(AdminToken, AppHost)
	body, _, err := c.ListUsers(1000, nil, nil)
	userResp := userListResponse{}
	err = json.Unmarshal(body, &userResp)
	if err != nil {
		fmt.Println("Failed to retrieve a list of users")
		return err
	}

	// Assume <handle>@aerofs.com users are not deleted
	// This ensures we do not delete users whose tokens we have
	// TODO : Should this be true for deployments too?
	for _, u := range userResp.Users {
		if !strings.Contains(u.Email, "aerofs.com") {
			err = c.DeleteUser(u.Email)
			if err != nil {
				fmt.Println("Unable to delete user %s", u.Email)
				return err
			}
		}
	}

	return nil
}

// Perform test teardown and setup
func TestMain(m *testing.M) {
	UserToken = os.Getenv("USERTOKEN")
	AdminToken = os.Getenv("ADMINTOKEN")
	AppHost = os.Getenv("APPHOST")

	//teardown
	err := removeUsers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rand.Seed(int64(os.Getpid()))
	os.Exit(m.Run())
}

// Create a new APIClient
func TestAPICreateClient(t *testing.T) {
	_, err := NewClient(AdminToken, AppHost)
	if err != nil {
		t.Fatal("Unable to create API client for testing")
	}
}

// Create a new User
func TestAPI_CreateUser(t *testing.T) {
	c, _ := NewClient(AdminToken, AppHost)
	email := fmt.Sprintf("test_email%d@moria.com", rand.Intn(10000))
	firstName := "Gimli"
	lastName := "Son of Gloin"

	b, _, e := c.CreateUser(email, firstName, lastName)
	if e != nil {
		t.Log("Error when attempting to create a user")
		t.Fatal(e)
	}

	t.Log("Successfully created the following new user")
	desc := User{}
	json.Unmarshal(b, &desc)
	t.Log(desc)
}

// List a set of Users
func TestAPI_ListUsers(t *testing.T) {
	c, _ := NewClient(AdminToken, AppHost)
	b, _, e := c.ListUsers(100, nil, nil)
	if e != nil {
		t.Log("Error when attempting to list users")
		t.Fatal(e)
	}

	t.Log("Successfully listed a set of users")
	desc := userListResponse{}
	json.Unmarshal(b, &desc)
	t.Log(desc.Users)
}

// Update an existing user
// Create a user, update their credentials and ensure they match
func TestAPI_UpdateUser(t *testing.T) {
	c, _ := NewClient(AdminToken, AppHost)

	email := fmt.Sprintf("test_email%d@moria.com", rand.Intn(10000))
	origUser := User{email, "Gimli", "Son of Gloin", []SharedFolder{}, []Invitation{}}
	new_firstName := "Sarumon"
	new_lastName := "Of Isengard"

	_, _, e := c.CreateUser(email, origUser.FirstName, origUser.LastName)
	if e != nil {
		t.Log("Error when attempting to create a user")
		t.Fatal(e)
	}

	b, _, e := c.UpdateUser(email, new_firstName, new_lastName)
	if e != nil {
		t.Log("Error when attempting to update a user")
		t.Fatal(e)
	}

	newUser := User{}
	e = json.Unmarshal(b, &newUser)
	if e != nil {
		t.Log("Error when attempting to unmarshal UserDescriptor")
		t.Fatal(e)
	}

	if reflect.DeepEqual(origUser, newUser) {
		t.Fatalf("New user %v is same from %v", newUser, origUser)
	}
	t.Log("New user %v is different from %v", newUser, origUser)
}

// Retrieve an uploadId, fileSize for an existing File
func TestAPI_GetUploadId(t *testing.T) {
	c, _ := NewClient(UserToken, AppHost)
	data, _, err := c.GetFolderChildren("root")
	if err != nil {
		t.Fatal("Error retrieving list of root Children")
	}

	children := Children{}
	err = json.Unmarshal(data, &children)
	if err != nil {
		t.Fatal("Error unmarshalling the children of root folder")
	}

	var fileId string
	var etag string
	for _, f := range children.Files {
		fileId = f.Id
		etag = f.Etag
		break
	}

	t.Logf("FileId,Etag are %s:%s", fileId, etag)
	uploadId, err := c.GetFileUploadId(fileId, []string{etag})
	if err != nil {
		t.Logf("Unable to get file upload id")
		t.Fatal(err)
	}
	t.Logf("UploadId for appconfig.json is %s", uploadId)
}

// Create a new user group
func TestAPI_CreateGroup(t *testing.T) {
	c, _ := NewClient(AdminToken, AppHost)
	groupName := fmt.Sprintf("testGroup_%d", rand.Intn(10000))

	body, _, err := c.CreateGroup(groupName)
	if err != nil {
		t.Logf("Unable to create new group %s", groupName)
		t.Fatal(err)
	}
	t.Log(string(body))
}
