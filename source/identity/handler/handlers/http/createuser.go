package httphandlers

import (
	"encoding/json"
	"fmt"
	"identity/graph"
	"identity/structs"
	"io/ioutil"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.http.createuser - initialized")
}

// CreateUser function
func CreateUser(w http.ResponseWriter, r *http.Request) {

	bodyBytes, bodyBytesErr := ioutil.ReadAll(r.Body)
	if bodyBytesErr != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al validar el cuerpo del mensaje"))
		return
	}

	var req structs.CreateUserRequest
	json.Unmarshal(bodyBytes, &req)

	createUserErr := graph.CreateUser(
		fmt.Sprintf("%v %v", req.GivenName, req.Surname),
		req.GivenName,
		req.Surname,
		"emailAddress",
		req.Email,
		req.Password)

	if createUserErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", createUserErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al crear el usuario"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Usuario creado satisfactoriamente"))
}
