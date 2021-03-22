package httphandlers

import (
	"encoding/json"
	"fmt"
	"identity/graph"
	"identity/settings"
	"identity/structs"
	"io/ioutil"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.http.authuser - initialized")
}

// AuthUser function
func AuthUser(w http.ResponseWriter, r *http.Request) {

	bodyBytes, bodyBytesErr := ioutil.ReadAll(r.Body)
	if bodyBytesErr != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al validar el cuerpo del mensaje"))
		return
	}

	var req structs.AuthUserRequest
	json.Unmarshal(bodyBytes, &req)

	userAuth, userAuthErr := graph.AuthenticateUser(settings.PrimaryTrustFrameworkPolicy, req.Username, req.Password)

	if userAuthErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", userAuthErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al autenticar el usuario"))
		return
	}

	payload, _ := json.Marshal(userAuth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
