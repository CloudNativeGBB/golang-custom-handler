package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"identity/settings"
	"identity/structs"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

func init() {
	fmt.Println("package: graph - initialized")
}

// AuthenticateUser function
func AuthenticateUser(trustFrameworkPolicy string, username string, password string) (*structs.UserAuth, error) {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("POST", fmt.Sprintf("https://%v/%v/%v/oauth2/v2.0/token", settings.TenantB2CLogin, settings.Tenant, trustFrameworkPolicy), nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	query := newRequest.URL.Query()
	query.Add("username", username)
	query.Add("password", password)
	query.Add("grant_type", "password")
	query.Add("scope", fmt.Sprintf("openid %v offline_access", settings.ClientID))
	query.Add("client_id", settings.ClientID)
	query.Add("response_type", "token id_token")
	newRequest.URL.RawQuery = query.Encode()

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		userAuth := structs.UserAuth{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &userAuth)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		return &userAuth, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}

// CreateUser function
func CreateUser(displayName string, givenName string, surname, signInType string, issuerAssignedID string, password string) error {

	identities := []structs.ObjectIdentity{}
	identity := structs.ObjectIdentity{}
	identity.SignInType = signInType
	identity.Issuer = settings.Tenant
	identity.IssuerAssignedID = issuerAssignedID
	identities = append(identities, identity)

	passwordProfile := structs.PasswordProfile{}
	passwordProfile.ForceChangePasswordNextSignInWithMfa = false
	passwordProfile.ForceChangePasswordNextSignIn = false
	passwordProfile.Password = password

	user := structs.User{}
	user.AccountEnabled = true
	user.DisplayName = displayName
	user.GivenName = givenName
	user.Surname = surname
	user.Identities = identities
	user.PasswordProfile = passwordProfile
	user.PasswordPolicies = "DisablePasswordExpiration,DisableStrongPassword"

	payload, _ := json.Marshal(user)
	payloadBuff := bytes.NewBuffer(payload)

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/users", payloadBuff)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	customErr := errors.New(string(bodyBytes))
	return customErr
}

// GetUser function
func GetUser(issuerAssignedID string) (*structs.User, error) {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users", nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	query := newRequest.URL.Query()
	query.Add("$filter", fmt.Sprintf("identities/any(id:id/issuer eq '%v' and id/issuerAssignedId eq '%v')", settings.Tenant, issuerAssignedID))
	query.Add("$select", "id,identities,displayname,givenname,surname,passwordpolicies")
	newRequest.URL.RawQuery = query.Encode()

	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		oDataContextUser := structs.ODataContextUser{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &oDataContextUser)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		if len(oDataContextUser.Value) > 0 {
			return &oDataContextUser.Value[0], nil
		}

		return nil, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}

// UpdateUser function
func UpdateUser(user *structs.User) error {

	payload, _ := json.Marshal(user)
	payloadBuff := bytes.NewBuffer(payload)

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("PATCH", fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%v", user.ID), payloadBuff)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	customErr := errors.New(string(bodyBytes))
	return customErr
}

// GetUserByObjectID function
func GetUserByObjectID(objectID string) (*structs.User, error) {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users", nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	query := newRequest.URL.Query()
	query.Add("$filter", fmt.Sprintf("id eq '%v'", objectID))
	query.Add("$select", "id,identities,displayname,givenname,surname,passwordpolicies")
	newRequest.URL.RawQuery = query.Encode()

	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		oDataContextUser := structs.ODataContextUser{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &oDataContextUser)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		if len(oDataContextUser.Value) > 0 {
			return &oDataContextUser.Value[0], nil
		}

		return nil, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}

// CreateGroup function
func CreateGroup(groupName string, groupDescription string) (*structs.Group, error) {

	group := structs.Group{}
	group.DisplayName = groupName
	group.Description = groupDescription
	group.MailEnabled = false
	group.MailNickname = uuid.New().String()
	group.SecurityEnabled = true

	payload, _ := json.Marshal(group)
	payloadBuff := bytes.NewBuffer(payload)

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/groups", payloadBuff)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		group := structs.Group{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &group)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		if &group != nil {
			return &group, nil
		}

		return nil, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}

// AddMember function
func AddMember(groupObjectID string, userObjectID string) error {

	directoryObject := structs.DirectoryObject{}
	directoryObject.ID = fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%v", userObjectID)

	payload, _ := json.Marshal(directoryObject)
	payloadBuff := bytes.NewBuffer(payload)

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%v/members/$ref", groupObjectID), payloadBuff)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	customErr := errors.New(string(bodyBytes))
	return customErr
}

// DeleteMember function
func DeleteMember(groupObjectID string, userObjectID string) error {

	directoryObject := structs.DirectoryObject{}
	directoryObject.ID = fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%v", userObjectID)

	payload, _ := json.Marshal(directoryObject)
	payloadBuff := bytes.NewBuffer(payload)

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("DELETE", fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%v/members/%v/$ref", groupObjectID, userObjectID), payloadBuff)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	customErr := errors.New(string(bodyBytes))
	return customErr
}

// ListGroup function
func ListGroup() (*[]structs.Group, error) {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/groups"), nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		oDataContextGroup := structs.ODataContextGroup{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &oDataContextGroup)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		if len(oDataContextGroup.Value) > 0 {
			return &oDataContextGroup.Value, nil
		}

		return nil, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}

// DeleteGroup function
func DeleteGroup(groupObjectID string) error {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("DELETE", fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%v", groupObjectID), nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return newRequestErr
	}

	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	customErr := errors.New(string(bodyBytes))
	return customErr
}

// ListMemberOf function
func ListMemberOf(userObjectID string) (*[]structs.Group, error) {

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%v/memberOf/microsoft.graph.group", userObjectID), nil)

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return nil, newRequestErr
	}

	newRequest.Header.Add("Authorization", settings.GlobalAuthToken)

	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return nil, responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return nil, bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		oDataContextGroup := structs.ODataContextGroup{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &oDataContextGroup)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return nil, bodyUnmarshalErr
		}

		if len(oDataContextGroup.Value) > 0 {
			return &oDataContextGroup.Value, nil
		}

		return nil, nil
	}

	customErr := errors.New(string(bodyBytes))
	return nil, customErr
}
