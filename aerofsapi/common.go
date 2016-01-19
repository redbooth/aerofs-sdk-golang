package aerofsapi

const (
	// Different scope levels applications can levy from users
	FileRead          = "files.read"
	FileWrite         = "files.write"
	FileAppData       = "files.appdata"
	UserRead          = "user.read"
	UserWrite         = "user.write"
	UserPassword      = "user.password"
	AclRead           = "acl.read"
	AclWrite          = "acl.write"
	AclInvitations    = "acl.invitations"
	OrganizationAdmin = "organization.admin"
)
