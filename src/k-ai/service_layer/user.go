/*
 * Copyright (c) 2017 by Peter de Vocht
 *
 * All rights reserved. No part of this publication may be reproduced, distributed, or
 * transmitted in any form or by any means, including photocopying, recording, or other
 * electronic or mechanical methods, without the prior written permission of the publisher,
 * except in the case of brief quotations embodied in critical reviews and certain other
 * noncommercial uses permitted by copyright law.
 *
 */

package service_layer

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"k-ai/db/db_model"
	"github.com/gocql/gocql"
	"k-ai/util"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gorilla/mux"
	"strings"
)

// create a new user
//
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// get body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ATJsonError(w, "read error:" + err.Error())
		return
	}

	// serialise item into object
	var user db_model.User
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		ATJsonError(w, "json error:" + err.Error())
		return
	}

	if len(user.Email) == 0 || len(user.Password_hash) == 0 || len(user.Surname) == 0 || len(user.First_name) == 0 {
		ATJsonError(w, "user save invalid object")
		return
	}

	// make sure it doesn't exist yet
	existing_user := db_model.User{Email: user.Email}
	err = existing_user.Get()
	if err == nil {
		ATJsonError(w, "a user with that email address already exists:" + existing_user.Email)
		return
	}

	// set salt and password hash
	rid, _ := gocql.RandomUUID()
	util.CopyUUID(&user.Salt, &rid)

	// turn password into a hash
	pwd_str := rid.String() + "-" + user.Password_hash
	pwd_hash := sha256.Sum256([]byte(pwd_str))
	user.Password_hash = hex.EncodeToString(pwd_hash[:])

	err = user.Save()
	if err != nil {
		ATJsonError(w, "save error:" + err.Error())
		return
	}

	// create a session
	session := db_model.Session{Email: user.Email, First_name: user.First_name, Surname: user.Surname}
	session_id, _ := gocql.RandomUUID()
	util.CopyUUID(&session.Session, &session_id)
	session.Save()

	// write session back to user
	json_data, _ := json.Marshal(&session)
	w.Write(json_data)
}


// check an existing session - return nil for error if exists
func ConfirmSession(session *db_model.Session) error {
	err := session.Get()
	if err != nil { return err }
	return nil
}


// Signin an existing user
//
func Signin(w http.ResponseWriter, r *http.Request) {
	// get body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ATJsonError(w, "read error:"+err.Error())
		return
	}
	// serialise item into object
	var user db_model.User
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		ATJsonError(w, "json error:" + err.Error())
		return
	}
	if len(user.Email) == 0 || len(user.Password_hash) == 0 {
		ATJsonError(w, "login invalid object")
		return
	}

	// convert passed in pwd string into a hash
	existing_password := user.Password_hash

	// turn password into a hash
	err = user.Get()
	if err != nil {
		ATJsonError(w, "login: " + err.Error())
		return
	}
	pwd_str := user.Salt.String() + "-" + existing_password
	hashed := sha256.Sum256([]byte(pwd_str))
	pwd_hash := hex.EncodeToString(hashed[:])

	if pwd_hash != user.Password_hash {
		ATJsonError(w, "login: password incorrect")
		return
	}

	// create a session
	session := db_model.Session{Email: user.Email, First_name: user.First_name, Surname: user.Surname}
	session_id, _ := gocql.RandomUUID()
	util.CopyUUID(&session.Session, &session_id)
	session.Save()

	// write session back to user
	json_data, _ := json.Marshal(&session)
	w.Write(json_data)
}


// Sigout an existing user
//
func Signout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// get signout
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	sid, err := gocql.ParseUUID(session)
	if err != nil {
		ATJsonError(w, "invalid session id: "+err.Error())
		return
	} else {
		session_obj := db_model.Session{Session: sid}
		err = session_obj.Delete()
		if err != nil {
			ATJsonError(w, "invalid session id: "+err.Error())
		} else {
			ATJsonMessage(w, http.StatusOK, "session deleted")
		}
	}
}
