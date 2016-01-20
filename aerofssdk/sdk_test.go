package aerofssdk

import (
	"fmt"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	"math/rand"
	"os"
	"strings"
	"testing"
)

// Oauth tokens with different permissions
var UserToken string
var AdminToken string

// The hostname of the AeroFS Appliance
var AppHost string

func TestMain(m *testing.M) {
	// Retrieve connection information
	UserToken = os.Getenv("USERTOKEN")
	AdminToken = os.Getenv("ADMINTOKEN")
	AppHost = os.Getenv("APPHOST")

	// Perform teardown
	err := rmUsers()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// Remove all of test-generated users
// Note that users with <email>@aerofs.com are persisted and not removed
func rmUsers() error {
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")
	users, e := ListUsers(c, 1000)
	if e != nil {
		return e
	}

	for _, user := range *users {
		uClient := UserClient{c, user}
		if !strings.Contains(uClient.Desc.Email, "aerofs.com") {
			err := uClient.Delete()
			if err != nil {
				fmt.Printf("Unable to remove users")
				return err
			}
		}
	}
	return nil
}

// Create a new user
func TestCreateUser(t *testing.T) {
	t.Logf("Creating new user")
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")

	t.Logf("Creating a new user")
	rand.Seed(int64(os.Getpid()))
	email := fmt.Sprintf("elrond.elf%d@middleearth.org", rand.Intn(100))
	firstName := "Melkor"
	lastName := "Bauglir"
	u, e := CreateUserClient(c, email, firstName, lastName)

	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to create new user")
	} else if u.Desc.Email == email && u.Desc.FirstName == firstName && u.Desc.LastName == lastName {
		t.Logf("Successfully created a new user")
		t.Log(*u)
	} else {
		t.Fatal("User created with incorrect fields")
		if u != nil {
			t.Fatal(*u)
		}
	}
}

// Update an already existing user
func TestUpdateUser(t *testing.T) {
	// Create new user
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")
	email := fmt.Sprintf("melkor.morgoth%d@gmail.com", rand.Intn(10000))
	firstName := "Melkor"
	lastName := "Bauglir"
	u, e := CreateUserClient(c, email, firstName, lastName)
	if e != nil {
		t.Fatalf("Unable to create new user : %s", e)
	}

	// Update created user
	t.Log(*u)
	e = u.Update("Eru", "Iluvatar")
	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to update user")
	} else {
		t.Logf("Successfully updated user")
		t.Log(*u)
	}

}

// Retrieve a list of backend users
func TestListUsers(t *testing.T) {
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")
	u, e := ListUsers(c, 1000)
	if e != nil {
		t.Fatalf("Unable to retrieve a list of users : %s", e)
	}
	if u != nil {
		t.Logf("There are %d users", len(*u))
		t.Log(*u)
	}
}

// Retrieve the root folder for a given user
func TestGetFolder(t *testing.T) {
	c, _ := api.NewClient(UserToken, "share.syncfs.com")
	f, e := NewFolderClient(c, "root", []string{"path", "children"})
	if e != nil {
		t.Fatalf("Unable to retrieve a FolderClient : %s", e)
	}

	f.LoadChildren()
	f.LoadMetadata()
	t.Log(*f)
}
