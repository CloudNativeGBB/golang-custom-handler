package httphandlers

import (
	"encoding/json"
	"fmt"
	"identity/graph"
	"identity/jwt"
	"identity/settings"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.http.getuser - initialized")
}

// GetUser function
func GetUser(w http.ResponseWriter, r *http.Request) {

	// verify token
	verifyTokenErr := jwt.ValidateToken(r, settings.PrimaryTrustFrameworkPolicy)

	if verifyTokenErr != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al validar el token del usuario"))
		return
	}

	username := r.URL.Query().Get("username")

	if len(username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al validar el parámetro de usuario"))
		return
	}

	user, userErr := graph.GetUser(username)

	if userErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", userErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al obtener el usuario"))
		return
	}

	payload, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
