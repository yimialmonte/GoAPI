package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/yimialmonte/GoAPI/user"
	"gopkg.in/mgo.v2/bson"
)

func bodyToUser(r *http.Request, u *user.User) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	if u == nil {
		return errors.New("a user is required")
	}
	db, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(db, u)
}

func userGetAll(w http.ResponseWriter, r *http.Request) {
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
	}

	postBodyResponse(w, http.StatusOK, jsonResponse{"users": users})
}

func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}

	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)
}
