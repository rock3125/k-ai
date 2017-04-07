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
)

// a user entity
type User struct {
	Email string            `json:"email"`
	First_name string       `json:"first_name"`
	Surname string       	`json:"surname"`
	Salt gocql.UUID    		`json:"salt"` // unique id of this item
	Password_hash string	`json:"password_hash"`
}

// save the current item to Cassandra
func (user *User) Save() error {
	if len(user.Email) == 0 || util.IsEmpty(&user.Salt) || len(user.Password_hash) == 0 ||
		len(user.Surname) == 0 || len(user.First_name) == 0 {
		return errors.New("User.Save() invalid parameter(s)")
	}

	value_map := make(map[string]interface{})
	value_map["email"] = user.Email
	value_map["first_name"] = user.First_name
	value_map["surname"] = user.Surname
	value_map["salt"] = user.Salt
	value_map["password_hash"] = user.Password_hash

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("user", value_map))
}

// load a KB item from db
func (user *User) Get() error {
	if len(user.Email) == 0 {
		return errors.New("User.Get() invalid parameters")
	}

	cols := []string{"first_name", "surname", "salt", "password_hash"}
	where_map := make(map[string]interface{})
	where_map["email"] = user.Email

	cql_str := db.Cassandra.SelectPaginated("user", cols, where_map, "", 0, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	var first_name, surname, password_hash string
	var salt gocql.UUID
	if iter.Scan(&first_name, &surname, &salt, &password_hash) {
		user.First_name = first_name
		user.Surname = surname
		user.Salt = salt
		user.Password_hash = password_hash
	} else {
		defer iter.Close()
		return errors.New("user does not exist")
	}
	return iter.Close()
}

