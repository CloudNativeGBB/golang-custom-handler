package httphandlers

import (
	"fmt"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.http.updateuser - initialized")
}

// UpdateUser function
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}
	fmt.Fprint(w, message)
}
