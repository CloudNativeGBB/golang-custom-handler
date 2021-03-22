package httphandlers

import (
	"encoding/json"
	"fmt"
	"identity/graph"
	"net/http"
)

func init() {
	fmt.Println("package: handlers.http.testgraph - initialized")
}

// TestGetUser function
func TestGetUser(w http.ResponseWriter, r *http.Request) {

	user, userErr := graph.GetUser("robece@correo.com.mx")

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

// TestUpdateUser function
func TestUpdateUser(w http.ResponseWriter, r *http.Request) {

	user, userErr := graph.GetUser("robece@correo.com.mx")

	if userErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", userErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al obtener el usuario"))
		return
	}

	user.DisplayName = "Gobeto Cervantes"
	user.GivenName = "Gobeto"

	updateUserErr := graph.UpdateUser(user)

	if updateUserErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", updateUserErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al actualizar el usuario"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Se actualizó el usuario satisfactoriamente"))
}

// TestGetUserByObjectID function
func TestGetUserByObjectID(w http.ResponseWriter, r *http.Request) {

	user, userErr := graph.GetUserByObjectID("7668613a-0301-48e5-8276-2353ce3dfafa")

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

// TestCreateGroup function
func TestCreateGroup(w http.ResponseWriter, r *http.Request) {

	group, groupErr := graph.CreateGroup("Test", "Test")

	if groupErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", groupErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al crear el grupo"))
		return
	}

	payload, _ := json.Marshal(group)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// TestAddMember function
func TestAddMember(w http.ResponseWriter, r *http.Request) {

	addMemberErr := graph.AddMember("22aeaee7-c8c8-4285-ba0b-b722a01ea27c", "005e113e-e952-44c4-91c5-01f83083ea39")

	if addMemberErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", addMemberErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al agregar el miembro al grupo"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Se agregó el miembro al grupo"))
}

// TestDeleteMember function
func TestDeleteMember(w http.ResponseWriter, r *http.Request) {

	deleteMemberErr := graph.DeleteMember("b35d2c86-70cd-4d21-9004-656483810557", "7668613a-0301-48e5-8276-2353ce3dfafa")

	if deleteMemberErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", deleteMemberErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al borrar el miembro del grupo"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Se borró el miembro del grupo"))
}

// TestListGroup function
func TestListGroup(w http.ResponseWriter, r *http.Request) {

	groups, groupsErr := graph.ListGroup()

	if groupsErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", groupsErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al obtener el listado de grupos"))
		return
	}

	payload, _ := json.Marshal(groups)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// TestDeleteGroup function
func TestDeleteGroup(w http.ResponseWriter, r *http.Request) {

	addMemberErr := graph.DeleteGroup("ef48d21e-3410-4360-bc27-272734c95b49")

	if addMemberErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", addMemberErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al borar el grupo"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Se borró el grupo"))
}

// TestListMemberOf function
func TestListMemberOf(w http.ResponseWriter, r *http.Request) {

	listMemberOf, listMemberOfErr := graph.ListMemberOf("005e113e-e952-44c4-91c5-01f83083ea39")

	if listMemberOfErr != nil {
		fmt.Println(fmt.Sprintf("Error: %v", listMemberOfErr))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Se presentó un error al obtener el listado de los grupos asociados al usuario"))
		return
	}

	payload, _ := json.Marshal(listMemberOf)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
