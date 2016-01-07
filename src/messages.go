package aerofs

// Structures used when communicating with an AeroFS Appliance

type NewName struct {
  Grant string
  Authcode string
  ID  string
  Secret  string
  RedirectURL string
}

type Authorization struct {
  Token string      `json:"access_token"`
  TokenType string  `json:"token_type"`
  ExpireTime string `json:"expires_in"`
  Scope string      `json:"scope"`
}

type SharedFolder struct {
  Id string                 `json:"id"`
  Name string               `json:"email"`
  External bool             `json:"is_external"`
  Members []SFMember        `json:"members"`
  Groups []SFGroupmember    `json:"groups"`
  Pending []SFPendingMember `json:"pending"`
  Permission string         `json:"caller_effective_permissions"`
}

type SFMember struct {
  Email string          `json:"email"`
  FirstName string      `json:"first_name"`
  LastName string       `json:"last_name"`
  Permissions string    `json:"permissions"`

type SFGroupMember struct {
  Id string           `json:"id"`
  Name string         `json:"name"`
  Permissions []string  `json:"permissions"`
}

type SFPendingMember struct {
  Email string        `json:"email"`
  FirstName string    `json:"first_name"`
  LastName  string    `json:"last_name"`
  Inviter string      `json:"invited_by"`
  Permissions string  `json:"permissions"`
  Note  string        `json:"note"`
}

type Group struct {
  Id string             `json:"id"`
  Name string           `json:"name"`
  Members []GroupMember `json:"members"`
}

type GroupMember struct {
  Email string      `json:"email"`
  FirstName string  `json:"first_name"`
  LastName string   `json:"last_name"`
}

type User struct {
  Email string              `json:"email"`
  FirstName string          `json:"first_name"`
  LastName string           `json:"last_name"`
  Shares []SharedFolder     `json:"shares"`
  Invitations []Invitation  `json:"invitations"`
}

type Invitee struct {
  EmailTo string    `json:"email_to"`
  EmailFrom string  `json:"email_from"`
  SignupCode string `json:"signup_code"`
}

type Invitation struct {
  Id string             `json:"shared_id"`
  Name string           `json:"shared_name"`
  Inviter string        `json:"invited_by"`
  Permissions []string  `json:"permissions"`
}
