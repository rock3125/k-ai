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

package db

import (
	"testing"
	"github.com/gocql/gocql"
	"strings"
	"k-ai/util_ut"
)

// test we can generate a good looking statement
func TestCqlGeneration1(t *testing.T) {

	cassandra := CCassandra{ keyspace: "test" }

	uuid := gocql.UUID{}
	for i, _ := range uuid {
		uuid[i] = byte(i)
	}
	columns := []string{"s1","s2"}
	nv := make(map[string]interface{},0)
	nv["intValue"] = 4
	nv["strValue"] = "Peter"
	nv["uuidValue"] = uuid

	// SELECT s1,s2 FROM test.cf WHERE strValue='Peter' AND uuidValue=00010203-0405-0607-0809-0a0b0c0d0e0f
	resultStr1 := cassandra.SelectPaginated("cf", columns, nv, "", 0, 1)
	util_ut.IsTrue(t, strings.Contains(resultStr1, "SELECT s1,s2 FROM test.cf WHERE "))
	util_ut.IsTrue(t, strings.Contains(resultStr1, "uuidValue=00010203-0405-0607-0809-0a0b0c0d0e0f"))
	util_ut.IsTrue(t, strings.Contains(resultStr1, "intValue=4"))
	util_ut.IsTrue(t, strings.Contains(resultStr1, "strValue='Peter'"))
	util_ut.IsTrue(t, strings.Contains(resultStr1, " LIMIT 1;"))

	resultStr2 := cassandra.SelectPaginated("cf", make([]string,0), nv, "", 0, 1)
	util_ut.IsTrue(t, strings.Contains(resultStr2, "SELECT * FROM test.cf WHERE "))

	resultStr3 := cassandra.SelectPaginated("cf", columns, nv, "uuidValue", uuid, 10)
	util_ut.IsTrue(t, strings.Contains(resultStr3, "SELECT s1,s2 FROM test.cf WHERE "))
	util_ut.IsTrue(t, strings.Contains(resultStr3, "intValue=4"))
	util_ut.IsTrue(t, strings.Contains(resultStr3, "strValue='Peter'"))
	util_ut.IsTrue(t, strings.Contains(resultStr3, "uuidValue>00010203-0405-0607-0809-0a0b0c0d0e0f"))
	util_ut.IsTrue(t, strings.Contains(resultStr3, " LIMIT 10;"))
}

// test delete
func TestCqlGeneration2(t *testing.T) {

	cassandra := CCassandra{ keyspace: "test" }

	uuid := gocql.UUID{}
	for i, _ := range uuid {
		uuid[i] = byte(i)
	}
	nv := make(map[string]interface{},0)
	nv["intValue"] = 4
	nv["strValue"] = "Peter"
	nv["uuidValue"] = uuid

	// SELECT s1,s2 FROM test.cf WHERE strValue='Peter' AND uuidValue=00010203-0405-0607-0809-0a0b0c0d0e0f
	resultStr := cassandra.Delete("cf", nv)
	util_ut.IsTrue(t, strings.Contains(resultStr, "DELETE FROM test.cf WHERE "))
	util_ut.IsTrue(t, strings.Contains(resultStr, "uuidValue=00010203-0405-0607-0809-0a0b0c0d0e0f"))
	util_ut.IsTrue(t, strings.Contains(resultStr, "intValue=4"))
	util_ut.IsTrue(t, strings.Contains(resultStr, "strValue='Peter'"))
}

// test insert
func TestCqlGeneration3(t *testing.T) {

	cassandra := CCassandra{ keyspace: "test" }

	uuid := gocql.UUID{}
	for i, _ := range uuid {
		uuid[i] = byte(i)
	}

	iv := make(map[string]interface{},0)
	iv["intValue"] = 4
	iv["strValue"] = "Peter"

	// SELECT s1,s2 FROM test.cf WHERE strValue='Peter' AND uuidValue=00010203-0405-0607-0809-0a0b0c0d0e0f
	resultStr := cassandra.Insert("cf", iv)
	util_ut.IsTrue(t, strings.Contains(resultStr, "INSERT INTO test.cf ("))
	util_ut.IsTrue(t, strings.Contains(resultStr, "(intValue,strValue)") || strings.Contains(resultStr, "(strValue,intValue)"))
	util_ut.IsTrue(t, strings.Contains(resultStr, "VALUES (4,'Peter');") || strings.Contains(resultStr, "VALUES ('Peter',4);"))
}

