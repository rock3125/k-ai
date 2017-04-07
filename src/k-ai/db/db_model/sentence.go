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
	"encoding/json"
	"k-ai/nlu/model"
	"k-ai/nlu/tokenizer"
)

// save a piece of text as sentences
func SaveText(sentence_list []model.Sentence, topic string) error {

	if len(sentence_list) == 0 || len(topic) == 0 {
		return errors.New("invalid parameters")
	}

	for _, sentence := range sentence_list {

		if util.IsEmpty(&sentence.Id) {  // make sure the sentence has a valid id
			return errors.New("invalid sentence id, null")
		}

		sentence.Topic = topic // set topic on the sentence too
		json_str, err := json.Marshal(sentence)
		if err != nil { return err }

		// topic grouping
		value_map := make(map[string]interface{})
		value_map["id"] = sentence.Id
		value_map["topic"] = topic
		err = db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("sentence_by_topic", value_map))
		if err != nil { return err }

		// sentence actual data save
		value_map_2 := make(map[string]interface{})
		value_map_2["id"] = sentence.Id
		value_map_2["topic"] = topic
		value_map_2["json_data"] = string(json_str)
		err = db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("sentence_by_id", value_map_2))
		if err != nil { return err }
	}
	return nil
}

// delete a sentence in both topic and by id
func DeleteText(id *gocql.UUID, topic string) error {
	if util.IsEmpty(id) {
		return errors.New("invalid parameter")
	}
	where_map := make(map[string]interface{})
	where_map["topic"] = topic

	err := db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("sentence_by_topic", where_map))
	if err != nil { return err }

	where_map2 := make(map[string]interface{})
	where_map2["id"] = id
	err = db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("sentence_by_id", where_map2))
	if err != nil { return err }

	return nil
}

// get text
func GetText(sentence_id *gocql.UUID) (*model.Sentence, error) {
	if util.IsEmpty(sentence_id) {
		return nil, errors.New("invalid parameter(s)")
	}
	where_map := make(map[string]interface{})
	where_map["id"] = sentence_id

	cols := []string{"topic", "json_data"}
	cql_str := db.Cassandra.SelectPaginated("sentence_by_id", cols, where_map, "", nil, 1)
	iter := db.Cassandra.Session.Query(cql_str).Iter()
	var topic, json_data string
	if iter.Scan(&topic, &json_data) {
		var text_item model.Sentence
		err := json.Unmarshal([]byte(json_data), &text_item)
		if err != nil { return nil, err }
		text_item.Topic = topic
		util.CopyUUID(&text_item.Id, sentence_id)
		return &text_item, iter.Close()
	}
	return nil, iter.Close()
}

// find a piece of text using the indexes
func FindText(tokenList []model.Token, topic string) (*model.ATResultList, error) {
	index_list, err := ReadIndexesWithFilterForTokens(tokenList, topic, 0)
	if err != nil { return nil, err }
	rs := model.ATResultList{ResultList: make([]model.ATResult,0)}

	// go through each index and get the associated text if possible
	for sentence_id, _ := range index_list {
		sentence, err := GetText(&sentence_id)
		if err == nil && sentence != nil && len(sentence.TokenList) > 0 {
			str := tokenizer.ToString(sentence.TokenList)
			rs.ResultList = append(rs.ResultList,
				model.ATResult{Text: str, Sentence_id: sentence_id, Topic: sentence.Topic})
		}
	}
	return &rs, nil
}

