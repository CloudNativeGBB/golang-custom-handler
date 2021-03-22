package timerhandlers

import (
	"encoding/json"
	"fmt"
	"identity/settings"
	"identity/structs"
	"io/ioutil"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.timer.publickeys - initialized")
	getPublicKeysInternal()
}

// GetPublicKeys function
func GetPublicKeys(w http.ResponseWriter, r *http.Request) {

	getPublicKeysErr := getPublicKeysInternal()

	if getPublicKeysErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", getPublicKeysErr))
		http.Error(w, getPublicKeysErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response":"success"}`))
	return
}

func getPublicKeysInternal() error {

	fmt.Println("Process of obtaining public keys started")

	keys, getPublicKeysErr := getPublicKeys(settings.PrimaryTrustFrameworkPolicy)

	if getPublicKeysErr != nil {
		return getPublicKeysErr
	}

	settings.GlobalPublicKeys = keys

	fmt.Println("Process of obtaining public keys finalized")

	return nil
}

func getPublicKeys(trustFrameworkPolicy string) (structs.Keys, error) {

	resp, getErr := http.Get(fmt.Sprintf("https://%v/%v/discovery/v2.0/keys?p=%v", settings.TenantB2CLogin, settings.Tenant, trustFrameworkPolicy))

	if getErr != nil {
		return structs.Keys{}, getErr
	}

	defer resp.Body.Close()
	body, readAllErr := ioutil.ReadAll(resp.Body)

	if readAllErr != nil {
		return structs.Keys{}, readAllErr
	}

	var keys structs.Keys
	unmarshalErr := json.Unmarshal(body, &keys)

	if unmarshalErr != nil {
		return keys, unmarshalErr
	}

	return keys, nil
}
