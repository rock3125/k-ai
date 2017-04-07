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
	"k-ai/db"
	"k-ai/util"
	"errors"
)


// create table if not exists <ks>.knowledge_base (
//    id uuid, kb text, origin text, json_data text,
//    primary key((kb), id)
// );


// a knowledge base entry
type KBEntry struct {
	Id        gocql.UUID    `json:"id"` // unique id of this item
	Topic string            `json:"topic"`
	Json_data string        `json:"json_data"` // actual info carrier
}

// save the current item to Cassandra
func (k *KBEntry) Save() error {
	if len(k.Json_data) == 0 || util.IsEmpty(&k.Id) || len(k.Topic) == 0 {
		return errors.New("KBEntry.Save() invalid parameter(s)")
	}

	value_map := make(map[string]interface{})
	value_map["json_data"] = k.Json_data
	value_map["id"] = k.Id
	value_map["topic"] = k.Topic

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("knowledge_base", value_map))
}

// load a KB item from db
func (k *KBEntry) Get() error {
	if util.IsEmpty(&k.Id) || len(k.Topic) == 0 {
		return errors.New("KBEntry.Get() invalid parameters")
	}

	cols := []string{"json_data"}
	where_map := make(map[string]interface{})
	where_map["id"] = k.Id
	where_map["topic"] = k.Topic

	cql_str := db.Cassandra.SelectPaginated("knowledge_base", cols, where_map, "", 0, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	var json_data string
	if iter.Scan(&json_data) {
		k.Json_data = json_data
	} else {
		defer iter.Close()
		return errors.New("item not found")
	}
	return iter.Close()
}


// load a KB item from db
func GetKBEntryById(id *gocql.UUID, topic string) (*KBEntry,error) {
	if util.IsEmpty(id) || len(topic) == 0 {
		return nil, errors.New("KBEntry.GetKBEntryById() invalid parameter(s)")
	}

	cols := []string{"json_data"}
	where_map := make(map[string]interface{})
	where_map["id"] = *id
	where_map["topic"] = topic

	cql_str := db.Cassandra.SelectPaginated("knowledge_base", cols, where_map, "", 0, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	var json_data string
	if iter.Scan(&json_data) {
		return &KBEntry{Id: *id, Json_data: json_data, Topic: topic}, nil
	} else {
		defer iter.Close()
		return nil, errors.New("item not found")
	}
}


// remove a KB item from db
func (k *KBEntry) Delete() error {
	if util.IsEmpty(&k.Id) || len(k.Topic) == 0 {
		return errors.New("KBEntry.Delete() invalid parameter(s)")
	}

	where_map := make(map[string]interface{})
	where_map["topic"] = k.Topic
	where_map["id"] = k.Id

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("knowledge_base", where_map))
}


// load a list of items, paginated
func GetKBEntryList(topic string, prev *gocql.UUID, page_size int) ([]KBEntry, error) {
	if len(topic) == 0 || page_size <= 0 {
		return nil, errors.New("KBEntry.GetKBEntryList() invalid parameter(s)")
	}

	cols := []string{"id", "json_data"}
	where_map := make(map[string]interface{})
	where_map["topic"] = topic

	cql_str := db.Cassandra.SelectPaginated("knowledge_base", cols, where_map, "id", prev, page_size)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	list := make([]KBEntry,0)

	var id gocql.UUID
	var json_data string

	for iter.Scan(&id, &json_data) {
		list = append(list, KBEntry{Id: id, Json_data: json_data, Topic: topic})
	}
	return list, iter.Close()
}

