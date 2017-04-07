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
	"k-ai/db/db_model"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"k-ai/nlu/parser"
)

// return a pagianted list of topics
func GetTopicList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	_, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}

	page := strings.ToLower(strings.TrimSpace(vars["page"]))
	page_size, err := strconv.Atoi(vars["page_size"])
	if err != nil {
		JsonError(w, "Invalid page-size: " + err.Error())
		return
	}

	topic_list, err := db_model.GetTopicList(page, page_size)
	if err != nil {
		JsonError(w, "GetTopicList error: " + err.Error())
		return
	}

	// return the json
	w.Header().Set("Content-Type", "application/json")
	json_bytes, _ := json.Marshal(topic_list)
	w.Write(json_bytes)
}

// remove an existing topic by name
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	topic_name := strings.ToLower(strings.TrimSpace(vars["topic_name"]))
	if len(topic_name) <= 2 || len(topic_name) > 50 {
		JsonError(w, "invalid topic name, a minimum of three characters and a maximum of 50 characters.")
		return
	}

	// log the event
	db_model.AddLogEntry(username, "remove a topic: " + topic_name)

	err = db_model.DeleteTopic(topic_name)
	if err != nil {
		JsonError(w, "DeleteTopic error: " + err.Error())
	}
	JsonMessage(w, http.StatusOK, "ok")
}

// save a topic to the db (which indexes too)
func SaveTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is a valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	topic_name := strings.ToLower(strings.TrimSpace(vars["topic_name"]))
	if len(topic_name) <= 2 || len(topic_name) > 50 {
		JsonError(w, "SaveTopic: invalid topic name, a minimum of three characters and a maximum of 50 characters.")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ATJsonError(w, "read error:" + err.Error())
		return
	}
	if len(body) > 10 && len(body) < 65536 {

		// remove older versions of this topic if possible
		db_model.DeleteTopic(topic_name)

		bodyStr := string(body)

		// log the event
		db_model.AddLogEntry(username, "save topic: " + topic_name)

		// body to sentence list
		sentence_list, err := parser.ParseText(bodyStr)
		if err != nil {
			ATJsonError(w, "Unexpected parser error:" + err.Error())
			return
		}

		if len(sentence_list) > 0 {
			// setup random ids for this set
			for i, _ := range sentence_list {
				sentence_list[i].RandomId()
			}
			err = db_model.SaveTopic(topic_name, bodyStr, sentence_list)
			if err != nil {
				JsonError(w, "SaveTopic: error: " + err.Error())
			} else {
				JsonMessage(w, http.StatusOK, "saved")
			}

		} else {
			JsonError(w, "SaveTopic: cannot parse text for topic: " + topic_name)
		}

	} else {
		JsonError(w, "SaveTopic: topic text length must be between 10 and 65536 bytes maximum")
	}
}


