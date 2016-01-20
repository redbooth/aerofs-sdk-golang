# aerofs-sdk-golang
An AeroFS Private Cloud API SDK written in Golang. The AeroFS Golang SDK is
composed of two packages: 
* **aerofsapi** -  Map the AeroFS API spec to individual calls
  * Supports all routes documented by the AeroFS API v1.3 Specification
* **aerofssdk** - Higher-level interface to the API
  * Supports the creation of File, Folder, Group, GroupMember, SharedFolder, SharedFolderMember and User objects

### Installation
```sh
$ go get github.com/aerofs/aerofs-sdk-golang/aerofsapi
$ go get github.com/aerofs/aerofs-sdk-golang/aerofssdk
```

### Testing
The API, SDK unit tests test against a local AeroFS Appliance. 
#### Do not execute the tests against a product instance as the tests mutate state.
Only run if you have a setup test instance. The tests require the following three 
environment variables to be set.
* **USERTOKEN** - An OAuth token with all permissions but **organization.admin**
* **ADMINTOKEN** - An OAuth token for a user with all permission scopes
* **APPHOST** - The hostname of the local AeroFS Appliance
```sh
$ cd aerofsapi
$ go test -v
$ cd ../aerofssdk
$ go test -v
```

### Melkor
Melkor is a test app that uses the API,SDK to enumerate lists of files, folders
and number of users on an AeroFS deployment
