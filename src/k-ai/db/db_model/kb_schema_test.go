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
	"github.com/gocql/gocql"
	"k-ai/util_ut"
)

// perform KBSchema tests
func TestKBSchema1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_schema_test")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_schema_test", 1)

	// field list
	field_list := make([]KBSchemaField,0)
	field_list = append(field_list, KBSchemaField{Name: "field 1", Semantic: "s1", Aiml: "who is *?"})
	field_list = append(field_list, KBSchemaField{Name: "field 2", Semantic: "s1", Aiml: "what is *?"})

	// create a new schema
	uuid, err := gocql.RandomUUID()
	util_ut.Check(t, err)
	schema_1 := KBSchema{Id: uuid, Origin: "peter", Name: "schema 1", Field_list: field_list}

	err = schema_1.SaveSchema()
	util_ut.Check(t, err)

	// reload the schema(s)
	schema_list, err := GetSchemaList()
	util_ut.Check(t, err)

	// must be one and test its fields etc.
	if len(schema_list) != 1 {
		t.Error("expected list size 1")
		return
	}

	schema_2 := schema_list[0]
	if schema_2.Name != "schema 1" || len(schema_2.Field_list) != 2 ||
	   len(schema_2.Field_list) != 2 {
		t.Error("invalid schema item re-read")
		return
	} else if schema_2.Field_list[0].Semantic != "s1" ||
		schema_2.Field_list[1].Semantic != "s1" {
		t.Error("invalid schema item re-read, semantic not set")
		return
	}

	// delete the schema
	schema_2.DeleteSchema()

	// make sure it's gone
	schema_list_2, err := GetSchemaList()
	util_ut.Check(t, err)
	// must be one and test its fields etc.
	if len(schema_list_2) != 0 {
		t.Error("expected list size 0 after delete")
	}

	db.DropKeyspace("localhost", "kai_ai_schema_test")
}

