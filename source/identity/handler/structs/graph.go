package structs

import (
	"fmt"
)

func init() {
	fmt.Println("package: structs.graph - initialized")
}

// PasswordProfile struct
type PasswordProfile struct {
	ForceChangePasswordNextSignInWithMfa bool   `json:"forceChangePasswordNextSignInWithMfa"`
	ForceChangePasswordNextSignIn        bool   `json:"forceChangePasswordNextSignIn"`
	Password                             string `json:"password"`
}

// ObjectIdentity struct
type ObjectIdentity struct {
	SignInType       string `json:"signInType"`
	Issuer           string `json:"issuer"`
	IssuerAssignedID string `json:"issuerAssignedId"`
}

// User struct
type User struct {
	ID               string           `json:"Id"`
	AccountEnabled   bool             `json:"accountEnabled"`
	DisplayName      string           `json:"displayName"`
	GivenName        string           `json:"givenName"`
	Surname          string           `json:"surname"`
	Identities       []ObjectIdentity `json:"identities"`
	PasswordProfile  PasswordProfile  `json:"passwordProfile"`
	PasswordPolicies string           `json:"passwordPolicies"`
}

// ODataContextUser struct
type ODataContextUser struct {
	ODataContext string `json:"@odata.context"`
	Value        []User `json:"value"`
}

// NoUserAuthToken struct
type NoUserAuthToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

// Group struct
type Group struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayName"`
	Description     string `json:"description"`
	MailEnabled     bool   `json:"mailEnabled"`
	MailNickname    string `json:"mailNickname"`
	SecurityEnabled bool   `json:"securityEnabled"`
}

// DirectoryObject struct
type DirectoryObject struct {
	ID string `json:"@odata.id"`
}

// ODataContextGroup struct
type ODataContextGroup struct {
	ODataContext string  `json:"@odata.context"`
	Value        []Group `json:"value"`
}

// UserAuth struct
type UserAuth struct {
	IDToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
