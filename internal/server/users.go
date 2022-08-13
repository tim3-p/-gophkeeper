package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tim3-p/gophkeeper/internal/common"
	"github.com/tim3-p/gophkeeper/internal/store"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Print("createUser")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var user common.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}

	var resp common.AddUserResponse
	resp.Name = user.Name
	resp.Status = "OK"
	resp.ID, err = serverStore.AddUser(user)
	if err != nil {
		log.Printf("AddUser() error: %v", err)
		if errors.Is(err, store.ErrAlreadyExists) {
			resp.Status = "already exists"
			writeStatus(w,
				http.StatusBadRequest,
				"User Already Exists",
			)
			return
		}
		resp.Status = "error"
		writeStatus(w,
			http.StatusBadRequest,
			"Cannot Add User",
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("cannot encode AddUserResponse: %v", err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	log.Print("changePassword")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var userInfo common.User
	err = json.Unmarshal(body, &userInfo)

	var resp common.AddUserResponse
	resp.Name = user
	resp.Status = "OK"
	err = serverStore.ChangeUserPassword(user, userInfo.Password)
	if err != nil {
		log.Printf("cannot change user password: %v", err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("cannot encode AddUserResponse: %v", err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}
