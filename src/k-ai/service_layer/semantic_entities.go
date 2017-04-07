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
	"net/http"
	"github.com/gorilla/mux"
	"strings"
	"regexp"
	"k-ai/nlu/lexicon"
	"k-ai/util"
	"os"
	"encoding/json"
	"sort"
	"k-ai/db/db_model"
)

//////////////////////////////////////////////////////////////////////////////////////////
// todo: semantic entities ARE saved with who dunnit, but aren't filtered as yet
// todo: this means that any user can overwrite another user's semantic entities
//////////////////////////////////////////////////////////////////////////////////////////

type SEResult struct {
	Name string 		`json:"name"`
	Semantic string     `json:"semantic"`
}

type SEResultList []SEResult

func (s SEResultList) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s SEResultList) Len() int      { return len(s) }
func (s SEResultList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }


// add an operation to the end of the lexicon updates file
func appendToLexiconUpdates(operation string, origin string, word string, semantic string) error {
	filename := util.GetDataPath() + "/lexicon/lexicon_updates.txt"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err == nil {
		defer f.Close()
		_, err = f.WriteString(operation + "|" + origin + "|" + word + ":" + semantic + "\n")
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}


// lexicon based semantics
// create, save, and find

// save a semantic entity /entities/save/{username}/{name}/{semantic}
func SaveSemanticEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	var isAscii = regexp.MustCompile(`^([a-z]|[A-Z]| |-)+$`)

	// {type}/{prev}/{page_size}/{json_field}/{query_str}
	name := strings.TrimSpace(vars["name"])
	semantic := strings.ToLower(strings.TrimSpace(vars["semantic"]))
	if len(name) == 0 || len(name) > 255 || len(semantic) == 0 || len(semantic) > 30 {
		JsonError(w, "invalid data, too small or large")
		return
	}

	// log the event
	db_model.AddLogEntry(username, "save semantic entity " + name + "=" + semantic)

	// check they're ASCII only for now
	if len(name) == 0 || !isAscii.MatchString(name) {
		JsonError(w, "invalid name value")
		return
	}
	if len(semantic) == 0 || !isAscii.MatchString(semantic) {
		JsonError(w, "invalid semantic value")
		return
	}
	lexicon.Lexi.Semantic[name] = semantic

	// and save it to make a note of it
	err = appendToLexiconUpdates("save", username, name, semantic)
	if err != nil {
		JsonError(w, err.Error())
	} else {
		JsonMessage(w, http.StatusAccepted,"ok")
	}
}

// remove an existing semantic entity /entities/delete/{name}
func DeleteSemanticEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	// {type}/{prev}/{page_size}/{json_field}/{query_str}
	name := strings.TrimSpace(vars["name"])

	// check they're ASCII only for now
	if len(name) == 0 {
		JsonError(w, "invalid name value")
		return
	}

	// log the event
	db_model.AddLogEntry(username, "delete semantic entity " + name)

	if _, ok := lexicon.Lexi.Semantic[name]; ok {
		// remove it from the map
		delete(lexicon.Lexi.Semantic, name)

		// and save it to make a note of it
		err := appendToLexiconUpdates("del", username, name, "")
		if err != nil {
			JsonError(w, err.Error())
		}
	}
	JsonMessage(w, http.StatusOK,"ok")
}

// find all items matching: /entities/find/{name}
func FindSemanticEntities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	// log the event
	db_model.AddLogEntry(username, "search for semantic entity " + vars["name"])

	find_name := strings.ToLower(strings.TrimSpace(vars["name"]))
	if len(find_name) <= 2 {
		JsonError(w, "invalid name value, too short, a minimum of three characters is required")
		return
	}

	result_list := make(SEResultList,0)
	for name, semantic := range lexicon.Lexi.Semantic {
		if strings.Contains(strings.ToLower(name), find_name) {
			result_list = append(result_list, SEResult{Name: name, Semantic: semantic})
		}
	}

	// sort list by name
	sort.Sort(result_list)

	// return the json
	w.Header().Set("Content-Type", "application/json")
	json_bytes, _ := json.Marshal(result_list)
	w.Write(json_bytes)
}

