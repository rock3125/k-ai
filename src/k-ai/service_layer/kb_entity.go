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

package service_layer

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/gocql/gocql"
	"encoding/json"
	"io/ioutil"
	"strings"
	"k-ai/db"
	"k-ai/nlu/parser"
	"k-ai/db/db_model"
	"k-ai/nlu/aiml"
	"regexp"
	"k-ai/util"
	"os"
	"k-ai/nlu/lexicon"
	"errors"
	"k-ai/nlu/model"
)

/**
 * list a set of knowledge-base entries paginated
 * @param request the request object
 * @param sessionIDStr the session of the active user
 * @param type the type of the entry to search for
 * @param prev the previous uuid string (or "null" for first page)
 * @param page_size number of items per page
 * @param json_field the json field to search on (ignored if this is "null" or query_str is "null")
 * @param query_str the query to execute (ignored if this is "null" or json_field is "null")
 * @return a list of knowledge-base entries or empty list if not found
 */
// server main entry point
func ListEntities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	_, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}

	// {type}/{prev}/{page_size}/{json_field}/{query_str}
	topic := vars["topic"]
	prev_str := vars["prev"]
	page_size,err := strconv.Atoi(vars["page_size"])
	if err != nil || page_size < 0 {
		JsonError(w,"invalid page_size, must be a positive integer")
		return
	}
	json_field := vars["json_field"]
	query_str := vars["query_str"]

	// set previous to empty if its null
	if prev_str == "null" {
		prev_str = ""
	}

	entity_list := make([]db_model.KBEntry,0)
	if json_field == "null" || query_str == "null" {
		if len(prev_str) > 0 {
			uuid, err := gocql.ParseUUID(prev_str)
			if err != nil || page_size < 0 {
				JsonError(w, "invalid previous id, 'prev' is not a valid guid")
				return
			}
			entity_list, err = db_model.GetKBEntryList(topic, &uuid, page_size)
			if err != nil {
				JsonError(w, err.Error())
				return
			}
		} else {
			entity_list, err = db_model.GetKBEntryList(topic, nil, page_size)
			if err != nil {
				JsonError(w, err.Error())
				return
			}
		}
	} else {
		//List<KBEntry> entryList = kbService.findPaginated(user.getOrganisation_id(), json_field, topic, query_str, prev_uuid, page_size);
		JsonError(w, "filter not implemented")
		return
	}

	json_str, err := json.Marshal(entity_list)
	if err != nil {
		JsonError(w,"invalid json, entity_list")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_str)
}


/**
 * delete a knowledge-base entry
 * @param kb the knowledge base system
 * @param id the ids
 * @return returns 200 on success
 */
func DeleteEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	// {kb}/{prev}/{page_size}/{json_field}/{query_str}
	topic := vars["topic"]
	id_str := vars["id"]

	uuid, err := gocql.ParseUUID(id_str)
	if err != nil {
		JsonError(w, "invalid, 'id' is not a valid guid")
	} else {
		// log the event
		db_model.AddLogEntry(username, "delete entity " + uuid.String() + "," +topic)

		kb_entry := db_model.KBEntry{Id: uuid, Topic: topic }
		err := kb_entry.Delete()
		if err != nil {
			JsonError(w,err.Error())
		} else {
			JsonMessage(w, http.StatusOK,"ok")
		}
	}
}

/**
 * save a knowledge-base entry
 * @param type the types
 * @return returns 200 on success
 */
func SaveEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var schema_item db_model.KBEntry
	err = decoder.Decode(&schema_item)
	if err != nil {
		JsonError(w,err.Error())
	} else {
		err = schema_item.Save()
		if err != nil {
			JsonError(w, err.Error())
		} else {
			// log the event
			db_model.AddLogEntry(username, "save entity " + schema_item.Id.String() + "," + schema_item.Topic)

			// tell the AIML system to reload - new information
			err = aiml.Aiml.Reload()
			if err != nil {
				JsonError(w, err.Error())
			} else {
				JsonMessage(w, http.StatusOK,"ok")
			}
		}
	}
}

// index a KB entry's data fields for indexing
func indexKBEntry(entry *db_model.KBEntry, schema_semantic_map map[string]string) error {
	var data_map map[string]interface{}
	err := json.Unmarshal([]byte(entry.Json_data), &data_map)
	if err == nil {
		kb_sentence := model.Sentence{Id: entry.Id, TokenList: make([]model.Token,0) }
		for key, value := range data_map {
			if key != "id" {
				// must be in the schema
				if _, ok := schema_semantic_map[key]; ok {
					valueStr := db.TypeToString(value)
					if len(valueStr) > 0 {
						sentence_list, err := parser.ParseText(valueStr)
						if err != nil { return err }
						for _, sentence := range sentence_list {
							for _, t_token := range sentence.TokenList {
								kb_sentence.TokenList = append(kb_sentence.TokenList, t_token)
							}
						}
					}
				}
			} // if not id field
		} // for each key value

		// index this whole lot as one sentence
		if len(kb_sentence.TokenList) > 0 {
			sentence_list := make([]model.Sentence, 0)
			sentence_list = append(sentence_list, kb_sentence)
			err = db_model.IndexText(entry.Topic, 0, sentence_list, 1.0)
			if err != nil { return err }
		}

	} // if valid json
	return err
}

// create a new lexicon for a given word and permanent storage
func writeToLexicon(lexicon_id string, semantic_map map[string][]string) error {
	if len(lexicon_id) > 0 {
		filename := util.GetDataPath() + "/lexicon/semantics/" + lexicon_id + ".txt"
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err == nil {
			defer f.Close()
			for key, word_list := range semantic_map {
				for _, word := range word_list {
					_, err = f.WriteString(word + ":" + key + "\n")
					if err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}
	} else {
		return errors.New("invalid filename")
	}
	return nil
}

// create a semantic lexicon for these items - skip numbers for now
func updateSemanticLexicon(entry *db_model.KBEntry,
							schema_semantic_map map[string]string,
							word_semantic_map map[string][]string) error {
	var data_map map[string]interface{}
	err := json.Unmarshal([]byte(entry.Json_data), &data_map)
	var isAscii = regexp.MustCompile(`^([a-z]|[A-Z]| |-)+$`)
	if err == nil {
		for key, value := range data_map {
			if key != "id" {
				// must be in the schema
				if semantic, ok := schema_semantic_map[key]; ok { // part of the schema?
					valueStr := strings.TrimSpace(db.TypeToString(value))
					// ascii only for now
					if len(semantic) > 0 && len(valueStr) > 0 && isAscii.MatchString(valueStr) {
						if word_list, ok := word_semantic_map[semantic]; ok {
							word_semantic_map[semantic] = append(word_list, valueStr)
						} else {
							word_semantic_map[semantic] = make([]string,0)
							word_semantic_map[semantic] = append(word_semantic_map[semantic], valueStr)
						}
						lexicon.Lexi.AddSemantic(valueStr, semantic)
					}
				}

			} // if not id field
		} // for each key value
	} // if valid json
	return err
}

// take all "wrong" characters out of an email address to use it as a filename
func sanitizeForFilename(username string) string {
	str := make([]byte,0)
	for i := 0; i < len(username); i++ {
		ch := username[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			str = append(str, username[i])
		} else {
			str = append(str, '_')
		}
	}
	return string(str)
}

/**
 * upload a set of instances
 * @return returns 200 on success
 */
func UploadInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	id_str := vars["id"]
	uuid, err := gocql.ParseUUID(id_str)
	if err != nil {
		JsonError(w, "invalid, 'id' is not a valid guid")
	} else {
		// get schema
		schema, err := db_model.GetSchemaById(&uuid)
		if err != nil {
			JsonError(w, "invalid, 'id' is not a valid schema guid")
			return
		}
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			JsonError(w, err.Error())
		} else {

			// log the event
			db_model.AddLogEntry(username, "upload instances for " + schema.Name + ", file string size: " + strconv.Itoa(len(contents)))

			// get existing fields from the schema proto-type (the field_list)
			field_map := make(map[string]string,0)
			for _, field := range schema.Field_list {
				field_map[field.Name] = field.Semantic
			}

			// record all words and their semantics
			word_semantic_map := make(map[string][]string,0)

			// process all lines
			for _, line := range strings.Split(string(contents), "\n") {
				if strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}") {

					var item_uuid gocql.UUID
					valid := false

					// create new json for saving entity
					str := "{"
					var data_map map[string]interface{}
					err := json.Unmarshal([]byte(line), &data_map)
					if err == nil {
						for key, value := range data_map {

							if key == "id" { // special case for KBEntry
								item_uuid, err = gocql.ParseUUID(db.TypeToString(value))
								valid = (err == nil)
							}
							if _, ok := field_map[key]; ok {
								if len(str) > 2 {
									str += ","
								}
								str += "\"" + key + "\": \"" + db.TypeToString(value) + "\""
							}

						}
					}
					str += "}"

					// save entity
					if !valid {
						item_uuid, _ = gocql.RandomUUID()
					}

					kb_entry := db_model.KBEntry{Id: item_uuid, Topic: schema.Name, Json_data: str}

					// save the item to the db
					kb_entry.Save()

					// index the item for retrieval
					indexKBEntry(&kb_entry, field_map)

					// update the lexicon with these new entities
					updateSemanticLexicon(&kb_entry, field_map, word_semantic_map)

				} // valid line?
			}

			// write this data to a new lexicon
			writeToLexicon(sanitizeForFilename(username), word_semantic_map)

			JsonMessage(w, http.StatusCreated,"ok")
		}
	}
}

