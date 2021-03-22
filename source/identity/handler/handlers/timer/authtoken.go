package timerhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"identity/settings"
	"identity/structs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	fmt.Println("package: handlers.timer.authtoken - initialized")
	getAuthTokenInternal()
}

// GetAuthToken function
func GetAuthToken(w http.ResponseWriter, r *http.Request) {

	authTokenErr := getAuthTokenInternal()

	if authTokenErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", authTokenErr))
		http.Error(w, authTokenErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response":"success"}`))
	return
}

func getAuthTokenInternal() error {

	fmt.Println("Process of obtaining authorization token started")

	authToken, authTokenErr := getAuthToken()

	if authTokenErr != nil {
		return authTokenErr
	}

	settings.GlobalAuthToken = authToken

	fmt.Println("Process of obtaining authorization token finalized")

	return nil
}

func getAuthToken() (string, error) {

	data := url.Values{}
	data.Set("client_id", settings.ClientID)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", settings.ClientSecret)
	data.Set("grant_type", "client_credentials")

	client := &http.Client{}
	newRequest, newRequestErr := http.NewRequest("POST", fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token", settings.Tenant), strings.NewReader(data.Encode()))

	if newRequestErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", newRequestErr))
		return "", newRequestErr
	}

	newRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, responseErr := client.Do(newRequest)

	if responseErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", responseErr))
		return "", responseErr
	}

	defer response.Body.Close()

	bodyBytes, bodyBytesErr := ioutil.ReadAll(response.Body)

	if bodyBytesErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", bodyBytesErr))
		return "", bodyBytesErr
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		authToken := structs.NoUserAuthToken{}
		bodyUnmarshalErr := json.Unmarshal(bodyBytes, &authToken)

		if bodyUnmarshalErr != nil {
			fmt.Println(fmt.Sprintf("Error: %v", bodyUnmarshalErr))
			return "", bodyUnmarshalErr
		}

		return authToken.AccessToken, nil
	}

	customErr := errors.New(string(bodyBytes))
	return "", customErr
}
