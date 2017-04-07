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
	"fmt"
	"k-ai/db"
	"testing"
	"github.com/gocql/gocql"
	"strconv"
	"k-ai/util_ut"
	"k-ai/util"
)

// perform CRUD on KBEntry
func TestKBEntry1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_test")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_test", 1)

	// see if it works
	for i := 0; i < 100; i++ {
		uuid, err := gocql.RandomUUID()
		util_ut.Check(t, err)
		e1 := KBEntry{Id: uuid, Json_data: "test json data " +strconv.Itoa(i), Topic: "topic1"}
		err = e1.Save()
		util_ut.Check(t, err)
	}

	// get paginated list after inserts
	var prev *gocql.UUID = nil  // pagination
	unique_map := make(map[string]int,0)
	total_size := 0
	size := 10
	var first_id gocql.UUID		// record UUID for delete
	var second_id gocql.UUID	// and for a get
	for size == 10 {
		list, err := GetKBEntryList("topic1", prev, 10)
		for _, item := range list {
			if util.IsEmpty(&first_id) {
				util.CopyUUID(&first_id, &item.Id)
			} else if util.IsEmpty(&second_id) {
				util.CopyUUID(&second_id, &item.Id)
			}
			total_size += 1
			unique_map[item.Json_data] = 1
		}
		util_ut.Check(t, err)
		size = len(list)
		if size == 10 {
			prev = &list[9].Id
		}
	}
	if total_size != 100 {
		t.Errorf("KBEntry: size %d, unique %d\n", total_size, len(unique_map))
		return
	} else {
		fmt.Print("KBEntry: insert 100 item test success\n")
	}

	// delete an item
	del_entry := KBEntry{Id: first_id, Topic: "topic1"}
	del_err := del_entry.Delete()
	util_ut.Check(t, del_err)

	// record and make sure it is gone
	size = 10
	prev = nil
	second_total_size := 0
	second_unique_map := make(map[string]int,0)
	for size == 10 {
		list, err := GetKBEntryList( "topic1", prev, 10)
		for _, item := range list {
			second_total_size += 1
			second_unique_map[item.Json_data] = 1
		}
		util_ut.Check(t, err)
		size = len(list)
		if size == 10 {
			prev = &list[9].Id
		}
	}
	if second_total_size != 99 {
		t.Errorf("KBEntry: size %d, unique %d\n", total_size, len(unique_map))
		return
	} else {
		fmt.Print("KBEntry: delete item test success\n")
	}

	// get a single entry
	get_entry := KBEntry{Id: second_id, Topic: "topic1"}
	get_err := get_entry.Get()
	util_ut.Check(t, get_err)
	if len(get_entry.Json_data) == 0 {
		t.Error("KBEntry: get failed")
	} else {
		fmt.Print("KBEntry: get item test success\n")
	}

	db.DropKeyspace("localhost", "kai_ai_test")
}

