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
	"encoding/json"
	"k-ai/util"
	"github.com/gocql/gocql"
	"errors"
	"k-ai/db"
)

// and two items that can wrap around json_data in kb_entry
type KBSchemaField struct {
	Name string					`json:"name"`		// field name
	Semantic string				`json:"semantic"`	// the field's semantic
	Aiml string                 `json:"aiml"`	    // aiml queries
}

type KBSchema struct {
	Id gocql.UUID           	`json:"id"`			// unique id of this item
	Name string					`json:"name"`		// the schema's name
	Origin string				`json:"origin"`		// who created the schema
	Field_list []KBSchemaField	`json:"field_list"`	// the fields of this item
}

// create a schema
func (schema *KBSchema) SaveSchema() error {
	if util.IsEmpty(&schema.Id) || len(schema.Name) == 0 || len(schema.Origin) == 0 {
		return errors.New("invalid parameters")
	}
	json_str, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	schemaEntry := KBEntry{Id: schema.Id, Json_data: string(json_str), Topic: "schema"}
	return schemaEntry.Save()
}

// delete a schema
func (schema *KBSchema) DeleteSchema() error {
	if util.IsEmpty(&schema.Id) {
		return errors.New("invalid parameter")
	}
	schemaEntry := KBEntry{Id: schema.Id, Topic: "schema"}
	return schemaEntry.Delete()
}

// load a KB item from db by id
func GetSchemaById(Id *gocql.UUID) (*KBSchema, error) {
	if util.IsEmpty(Id) {
		return nil, errors.New("KBSchema.Get() invalid parameters")
	}
	cols := []string{"json_data"}
	where_map := make(map[string]interface{})
	where_map["id"] = Id
	where_map["topic"] = "schema"

	cql_str := db.Cassandra.SelectPaginated("knowledge_base", cols, where_map, "", 0, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	var json_data string
	if iter.Scan(&json_data) {
		var schema KBSchema
		err := json.Unmarshal([]byte(json_data), &schema)
		if err != nil {
			return nil, err
		}
		return &schema, nil
	} else {
		defer iter.Close()
		return nil, errors.New("item not found")
	}
}

// return a list of all schemas
func GetSchemaList() ([]KBSchema,error) {
	return_list := make([]KBSchema,0)
	var prev *gocql.UUID = nil
	page_size := 100
	size := page_size
	for size == page_size {
		list, err := GetKBEntryList("schema", prev, page_size)
		if err != nil {
			return return_list, err
		}
		for _, item := range list {
			var schema_item KBSchema
			err := json.Unmarshal([]byte(item.Json_data), &schema_item)
			if err != nil {
				return return_list, err
			}
			return_list = append(return_list, schema_item)
			prev = &schema_item.Id
		}
		size = len(list)
	}
	return return_list, nil
}


// return a list of all by name -> schema item
func GetSchemaMap() (map[string]KBSchema, error) {
	return_map := make(map[string]KBSchema,0)
	var prev *gocql.UUID = nil
	page_size := 100
	size := page_size
	for size == page_size {
		list, err := GetKBEntryList("schema", prev, page_size)
		if err != nil {
			return return_map, err
		}
		for _, item := range list {
			var schema_item KBSchema
			err := json.Unmarshal([]byte(item.Json_data), &schema_item)
			if err != nil { return return_map, err }
			return_map[schema_item.Name] = schema_item
			prev = &schema_item.Id
		}
		size = len(list)
	}
	return return_map, nil
}

