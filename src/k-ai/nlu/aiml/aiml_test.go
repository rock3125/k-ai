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

package aiml

import (
	"testing"
	"github.com/gocql/gocql"
	"k-ai/nlu/tokenizer"
	"k-ai/nlu/model"
	"k-ai/db/db_model"
	"encoding/json"
	"k-ai/db"
	"strconv"
	"time"
	"k-ai/util_ut"
)

const name_field = `[{"tokenList":[{"index":0,"list":[0],"tag":"NNP","text":"Peter de Vocht","dep":"compound","synid":-1,"semantic":"person"}]}]`
const loc_field = `[{"tokenList":[{"index":0,"list":[],"tag":"NNP","text":"Spain","dep":"ROOT","synid":-1,"semantic":""}]}]`

const whereIsPeter = `[{"index":0,"list":[1],"tag":"WRB","text":"where","dep":"advmod","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`
const activateTheRobot = `[{"index":0,"list":[],"tag":"VB","text":"Activate","dep":"ROOT","synid":-1,"semantic":""},{"index":1,"list":[2,0],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":2,"list":[0],"tag":"NNP","text":"Robot","dep":"dobj","synid":-1,"semantic":""},{"index":3,"list":[0],"tag":".","text":"!","dep":"punct","synid":-1,"semantic":""}]`
const whenWillYouHaveABody = `[{"index":0,"list":[3],"tag":"WRB","text":"When","dep":"advmod","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":"man"},{"index":2,"list":[3],"tag":"PRP","text":"you","dep":"nsubj","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"VB","text":"have","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[5,3],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":5,"list":[3],"tag":"NN","text":"body","dep":"dobj","synid":-1,"semantic":""},{"index":6,"list":[3],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`
const aisAreStupid = `[{"index":0,"list":[1],"tag":"NNS","text":"AIs","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBP","text":"are","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[4,1],"tag":"JJ","text":"stupid","dep":"amod","synid":-1,"semantic":""},{"index":3,"list":[4,1],"tag":"JJ","text":"little","dep":"amod","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":"NNS","text":"things","dep":"attr","synid":-1,"semantic":""},{"index":5,"list":[1],"tag":".","text":"!","dep":"punct","synid":-1,"semantic":""}]`
const aisAre = `[{"index":0,"list":[1],"tag":"NNS","text":"AIs","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBP","text":"are","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":".","text":"!","dep":"punct","synid":-1,"semantic":""}]`
const whatIsTheYear = `[{"index":0,"list":[1],"tag":"WP","text":"What","dep":"attr","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[3,1],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"NN","text":"year","dep":"nsubj","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`
const whatIsTheTime = `[{"index":0,"list":[1],"tag":"WP","text":"what","dep":"attr","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[3,1],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"NN","text":"time","dep":"nsubj","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`

// from json str back to a token list
func jsonToTokenList(t *testing.T, str string) []model.Token {
	var token_list []model.Token
	err := json.Unmarshal([]byte(str), &token_list)
	util_ut.Check(t, err)
	return token_list
}

// parse a json structure and append to sentence list
func parseText(t *testing.T, json_str string) []model.Sentence {
	var sl []model.Sentence
	err := json.Unmarshal([]byte(json_str),&sl)
	util_ut.Check(t, err)
	return sl
}

// index a KB entry's data fields for indexing
// set all sentence ids of these fields and items equal to the KB's id
func indexKBEntry(t *testing.T, entry *db_model.KBEntry) error {
	var data_map map[string]interface{}
	err := json.Unmarshal([]byte(entry.Json_data), &data_map)
	if err == nil {
		final_sentence := model.Sentence{Id: entry.Id, TokenList: make([]model.Token,0)}
		for key, value := range data_map {
			if key != "id" {
				valueStr := db.TypeToString(value)
				if len(valueStr) > 0 {
					if key == "name" {
						sentence_list := parseText(t, name_field)
						for _, sentence := range sentence_list {
							for _, token := range sentence.TokenList {
								final_sentence.TokenList = append(final_sentence.TokenList, token)
							}
						}
					} else {
						sentence_list := parseText(t, loc_field)
						for _, sentence := range sentence_list {
							for _, token := range sentence.TokenList {
								final_sentence.TokenList = append(final_sentence.TokenList, token)
							}
						}
					}
				}
			} // if not id field
		} // for each key value

		final_sentence_list := make([]model.Sentence,0)
		final_sentence_list = append(final_sentence_list, final_sentence)
		err = db_model.IndexText(entry.Topic, 0, final_sentence_list, 1.0)
		if err != nil { return err }

	} // if valid json
	return err
}


// test a simple rule with name substitution (from ai.aiml, line 82)
func TestAiml1(t *testing.T) {
	Aiml.initFromFile()

	result := Aiml.MatchTokenList(jsonToTokenList(t, activateTheRobot))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	}
	str := "AI activated. Awaiting your command."
	if result[0].Text != str {
		t.Errorf("expected %s but got %s", str, result[0].Text)
	}
}

// test "WHEN WILL YOU * BODY"  (ai.aiml, line 173)
func TestAiml2(t *testing.T) {
	
	Aiml.initFromFile()

	result := Aiml.MatchTokenList(jsonToTokenList(t, whenWillYouHaveABody))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	}
	str := "I will finish the robot body as soon as I can raise the funds for it."
	if result[0].Text != str {
		t.Errorf("expected %s but got %s", str, result[0].Text)
	}
	if len(result[0].TokenList) != 2 {  // check "have a" got bound to "*"
		t.Error("expected 'have a' to be bound")
	}
	tokenStr := tokenizer.ToString(result[0].TokenList)
	if tokenStr != "have a" {
		t.Error("expected 'have a' to be bound")
	}
}

// test "(ROBOTS|ROBOT|AI|AIS|artificial intelligence) (ARE|IS) *", ai.aiml line 141
func TestAiml3(t *testing.T) {
	
	Aiml.initFromFile()

	result := Aiml.MatchTokenList(jsonToTokenList(t, aisAreStupid))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	}
	str := "I posses no emotions."
	if result[0].Text != str {
		t.Errorf("expected %s but got %s", str, result[0].Text)
	}
	expect_str := "stupid little things"
	if len(result[0].TokenList) != 3 {  // check "have a" got bound to "*"
		t.Errorf("expected '%s' to be bound", expect_str)
	}
	tokenStr := tokenizer.ToString(result[0].TokenList)
	if tokenStr != expect_str {
		t.Errorf("expected '%s' to be bound", expect_str)
	}
}


// test time (date.aiml, line 29)
func TestAimlDateTime(t *testing.T) {
	
	Aiml.initFromFile()

	result := Aiml.MatchTokenList(jsonToTokenList(t, whatIsTheTime))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	}
	str := "{time}"
	if result[0].Text != str {
		t.Errorf("expected %s but got %s", str, result[0].Text)
	}
}


// test empty pattern for "(ROBOTS|ROBOT|AI|AIS|artificial intelligence) (ARE|IS) *", ai.aiml line 141
func TestAiml4(t *testing.T) {
	
	Aiml.initFromFile()

	result := Aiml.MatchTokenList(jsonToTokenList(t, aisAre))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	} else {
		str := "I posses no emotions."
		if result[0].Text != str {
			t.Errorf("expected %s but got %s", str, result[0].Text)
		} else {
			if len(result[0].TokenList) != 0 { // check "have a" got bound to "*"
				t.Error("expected empty match")
			}
		}
	}
}


// test indexed database search for AIML with the ATResultList system
func TestAiml_db(t *testing.T) {
	db_model.Delete_and_create_keyspace_for_unit_test("kai_ai_schema_aiml_test")

	// field list
	field_list := make([]db_model.KBSchemaField,0)
	field_list = append(field_list, db_model.KBSchemaField{Name: "name", Aiml: "who is *?\nwhere is *?"})
	field_list = append(field_list, db_model.KBSchemaField{Name: "location"})

	// create a new schema
	schema_uuid, err := gocql.RandomUUID()
	util_ut.Check(t, err)
	schema_1 := db_model.KBSchema{Id: schema_uuid, Origin: "peter", Name: "address", Field_list: field_list}

	err = schema_1.SaveSchema()
	util_ut.Check(t, err)

	// setup the aiml system with the db
	err = Aiml.SetupDbSchema()
	util_ut.Check(t, err)

	// add a record to the database for our new schema
	kb_uuid, err := gocql.RandomUUID()
	util_ut.Check(t, err)

	// create a new kb entry in the system
	kbe1 := db_model.KBEntry{Id: kb_uuid, Topic: "address",
		Json_data:               "{\"name\": \"Peter de Vocht\", \"location\": \"Spain\"}"}

	err = kbe1.Save()  // save it to the db
	util_ut.Check(t, err)

	err = indexKBEntry(t, &kbe1) // make it find-able on its fields
	util_ut.Check(t, err)

	// use the db to find the kb entry
	result := Aiml.MatchTokenList(jsonToTokenList(t, whereIsPeter))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	} else {
		rs, err := Aiml.PerformSpecialOps(result, "address")
		util_ut.Check(t, err)
		if len(rs.ResultList) != 1 {
			t.Errorf("len(result) != 1 after special ops, but %d", len(rs.ResultList))
		} else {
			// get the url
			for _, item := range rs.ResultList {
				if item.Text != "name: Peter de Vocht, location: Spain" {
					t.Error("expected: \"name: Peter de Vocht, location: Spain\", got %s", item.Text)
				}
			}
		}
	}

	db_model.Delete_keyspace_after_unit_test("kai_ai_schema_aiml_test")
}


// test normal AIML template match with {} replacement for year
func TestAiml5(t *testing.T) {
	
	Aiml.initFromFile()

	db_model.Delete_and_create_keyspace_for_unit_test("kai_ai_schema_aiml_test_2")

	result := Aiml.MatchTokenList(jsonToTokenList(t, whatIsTheYear))
	if result == nil || len(result) != 1 {
		t.Errorf("len(result) != 1, but %d", len(result))
	}
	str := "{year}"
	if result[0].Text != str {
		t.Errorf("expected %s but got %s", str, result[0].Text)
	} else {
		rs, err := Aiml.PerformSpecialOps(result, "peter@peter.co.nz")
		util_ut.Check(t, err)
		if len(rs.ResultList) != 1 {
			t.Errorf("len(result) != 1 after special ops, but %d", len(rs.ResultList))
		} else {
			// get the url
			for _, item := range rs.ResultList {
				year_str := strconv.Itoa(time.Now().Year())
				if item.Text != year_str {
					t.Error("expected: \"%s\", got \"%s\"", year_str, item.Text)
				}
			}
		}
	}

	db_model.Delete_keyspace_after_unit_test("kai_ai_schema_aiml_test_2")
}

