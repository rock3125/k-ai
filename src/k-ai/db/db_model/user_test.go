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
	"testing"
	"k-ai/db"
	"k-ai/util_ut"
	"github.com/gocql/gocql"
	"k-ai/util"
)

// test user CRUD
func TestUser1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_user")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_user", 1)

	// create a user and save it
	salt, _ := gocql.RandomUUID()
	user := User{Email: "peter@peter.co.nz", Password_hash: "password hash",
					First_name: "Peter", Surname: "de Vocht", Salt: salt}
	err := user.Save()
	util_ut.IsTrue(t, err == nil)

	// retrieve this user
	user2 := User{Email: "peter@peter.co.nz"}
	err = user2.Get()
	util_ut.IsTrue(t, err == nil)

	util_ut.IsTrue(t, user2.First_name == "Peter")
	util_ut.IsTrue(t, user2.Surname == "de Vocht")
	util_ut.IsTrue(t, user2.Password_hash == "password hash")
	util_ut.IsTrue(t, util.EqualUUID(&user2.Salt, &salt))

	db.DropKeyspace("localhost", "kai_ai_user")
}

