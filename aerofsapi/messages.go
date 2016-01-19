package aerofsapi

// Structures used when communicating with an AeroFS Appliance

type User struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

type Folder struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Parent    string     `json:"parent"`
	IsShared  bool       `json:"is_shared"`
	Sid       string     `json:"sid"`
	Path      ParentPath `json:"path"`
	ChildList Children   `json:"children"`
	Etag      string
}

type File struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Parent       string     `json:"parent"`
	LastModified string     `json:"last_modified"`
	Size         int        `json:"size"`
	Mime         string     `json:"mime_type"`
	Etag         string     `json:"etag"`
	Path         ParentPath `json:"path"`
	ContentState string     `json:"content_state"`
}

type ParentPath struct {
	Folders []Folder `json:"folders"`
}

type Children struct {
	Folders []Folder `json:"folders"`
	Files   []File   `json:"files"`
}

type SharedFolder struct {
	Id         string            `json:"id,omitempty"`
	Name       string            `json:"name"`
	External   bool              `json:"is_external,omitempty"`
	Members    []SFMember        `json:"members,omitempty"`
	Groups     []SFGroupMember   `json:"groups,omitempty"`
	Pending    []SFPendingMember `json:"pending,omitempty"`
	Permission []string          `json:"caller_effective_permissions,omitempty"`
}

type SFMember struct {
	Sid         string   `json:"-"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Permissions []string `json:"permissions"`
}

type SFGroupMember struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type SFPendingMember struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name,omitempty"`
	LastName    string   `json:"last_name,omitempty"`
	Inviter     string   `json:"invited_by,omitempty"`
	Permissions []string `json:"permissions"`
	Note        string   `json:"note"`
}

type Group struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Members []GroupMember `json:"members"`
}

type GroupMember struct {
	GroupId   string `json:"-"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Invitee struct {
	EmailTo    string `json:"email_to"`
	EmailFrom  string `json:"email_from"`
	SignupCode string `json:"signup_code,omitempty"`
}

type Invitation struct {
	Sid         string   `json:"share_id"`
	Name        string   `json:"share_name"`
	Inviter     string   `json:"invited_by"`
	Permissions []string `json:"permissions"`
}

type Device struct {
	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	OSFamily    string `json:"os_family"`
	InstallDate string `json:"install_data"`
}

type PermissionList struct {
	Permissions []string `json:"permissions"`
}

type DeviceStatus struct {
	Online   bool   `json:"online"`
	LastSeen string `json:"last_seen"`
}

// Response specific structures

type userListResponse struct {
	HasMore bool   `json:"has_more"`
	Users   []User `json:"data"`
}
