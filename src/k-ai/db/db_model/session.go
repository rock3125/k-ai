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

package db_model

import (
	"github.com/gocql/gocql"
	"k-ai/util"
	"k-ai/db"
	"errors"
	"strings"
)

// a user entity
type Session struct {
	Email string            `json:"email"`
	First_name string       `json:"first_name"`
	Surname string       	`json:"surname"`
	Session gocql.UUID    	`json:"session"`
}

// save the current session to Cassandra
func (session *Session) Save() error {
	if len(session.Email) == 0 || util.IsEmpty(&session.Session) || len(session.Surname) == 0 ||
		len(session.First_name) == 0 {
		return errors.New("Session.Save() invalid parameter(s)")
	}

	value_map := make(map[string]interface{})
	value_map["email"] = session.Email
	value_map["first_name"] = session.First_name
	value_map["surname"] = session.Surname
	value_map["session"] = session.Session

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("session", value_map))
}

// delete current session object
func (session *Session) Delete() error {
	if util.IsEmpty(&session.Session) {
		return errors.New("Session.Delete() invalid parameter(s)")
	}

	value_map := make(map[string]interface{})
	value_map["session"] = session.Session

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("session", value_map))
}

// load a session using its id
func (session *Session) Get() error {
	if util.IsEmpty(&session.Session) {
		return errors.New("Session.Get() invalid parameters")
	}

	cols := []string{"first_name", "surname", "email"}
	where_map := make(map[string]interface{})
	where_map["session"] = session.Session

	cql_str := db.Cassandra.SelectPaginated("session", cols, where_map, "", 0, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	var first_name, surname, email string
	if iter.Scan(&first_name, &surname, &email) {
		session.First_name = first_name
		session.Surname = surname
		session.Email = email
	} else {
		defer iter.Close()
		return errors.New("session does not exist")
	}
	return iter.Close()
}

// get the username for a user object
func (session Session) GetUserName() string {
	return strings.TrimSpace(session.First_name + " " + session.Surname)
}

// validate a session object and return it if valid
func ValidateSession(session string) (*Session, error) {
	sid, err := gocql.ParseUUID(session)
	if err != nil {
		return nil, errors.New("invalid session id")
	}
	session_obj := Session{Session: sid}
	err = session_obj.Get()
	if err != nil {
		return nil, err
	}
	return &session_obj, nil
}

