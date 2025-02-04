package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// authenticate the user performing the action
	user, code, err := rt.AuthenticateUserFromParameter("uname", r, ps)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// get the user to be banned from the resource parameter
	bannedUser, code, err := rt.GetUserFromParameter("banned_uname", r, ps)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// check whether the user performing the ban and the user
	// to be banned are the same
	if user.Id == bannedUser.Id {
		http.Error(w, ErrSelfBan.Error(), http.StatusBadRequest)
		return
	}

	// insert the ban into the database
	err = rt.db.InsertBan(user.UserIntoDatabaseUser(), bannedUser.UserIntoDatabaseUser())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	// return the banned user
	_ = json.NewEncoder(w).Encode(bannedUser)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// authenticate the user performing the action
	user, code, err := rt.AuthenticateUserFromParameter("uname", r, ps)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// get the banned user from the resource parameter
	bannedUser, code, err := rt.GetUserFromParameter("banned_uname", r, ps)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// remove the ban from the database
	err = rt.db.DeleteBan(user.UserIntoDatabaseUser(), bannedUser.UserIntoDatabaseUser())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204
}
