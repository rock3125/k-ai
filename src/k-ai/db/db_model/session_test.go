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
)

// test user CRUD
func TestSession1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_session")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_session", 1)

	// create a user and save it
	sess, _ := gocql.RandomUUID()
	session := Session{Session: sess, Email: "peter@peter.co.nz", Surname: "de Vocht", First_name: "Peter"}
	err := session.Save()
	util_ut.IsTrue(t, err == nil)

	// retrieve this user
	session2 := Session{Session: sess}
	err = session2.Get()
	util_ut.IsTrue(t, err == nil)

	util_ut.IsTrue(t, session2.First_name == "Peter")
	util_ut.IsTrue(t, session2.Surname == "de Vocht")
	util_ut.IsTrue(t, session2.Email == "peter@peter.co.nz")

	db.DropKeyspace("localhost", "kai_ai_session")
}

