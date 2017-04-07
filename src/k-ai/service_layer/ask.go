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
	"io/ioutil"
	"k-ai/util"
	"k-ai/nlu/parser"
	"k-ai/nlu/aiml"
	"k-ai/nlu/model"
	"encoding/json"
	"k-ai/db/db_model"
	"math/rand"
	"github.com/gorilla/mux"
	"strings"
	"k-ai/db/freebase"
)

// perform an ask
//
func Ask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// check session is valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ATJsonError(w, "read error:" + err.Error())
		return
	}
	if len(body) > 0 && len(body) < 255 {
		bodyStr := string(body)
		// log the event
		db_model.AddLogEntry(username, "query:" + bodyStr)

		// 1. parse it
		sentence_list, err := parser.ParseText(bodyStr)
		if err != nil {
			ATJsonError(w, "Unexpected parser error:" + err.Error())
			return
		}
		if len(sentence_list) == 0 {
			ATJsonError(w, "Please ask me a proper question")
		} else if len(sentence_list) > 1 {
			ATJsonError(w, "Question too complex, more than one sentence")
		} else {
			sentence := sentence_list[0]  // the sentence

			if !sentence.IsQuestion() && !sentence.IsImperative() {
				ATJsonError(w, "That does not look like a question.  Ask me question or give me a command please.")
				return
			}

			ask_teach_result := model.ATResultList{ ResultList: make([]model.ATResult,0) }

			//////////////////////////////////////////////////////////////////////
			// 1. perform an AIML query
			{
				binding_list := aiml.Aiml.MatchTokenList(sentence.TokenList)
				if len(binding_list) > 0 {

					// if it has more than one binding, pick a random one from the list
					if len(binding_list) > 1 {
						answer := binding_list[rand.Intn(len(binding_list))]
						binding_list = []model.AimlBinding{answer}
					}

					rs, err := aiml.Aiml.PerformSpecialOps(binding_list, username)
					if err != nil {
						ATJsonError(w, err.Error())
						return
					}
					// append results to return set
					for _, item := range rs.ResultList {
						ask_teach_result.ResultList = append(ask_teach_result.ResultList, item)
					}
				}
			}

			// replace you, your, yourself pronoun references with KAI
			db_model.ResolveFirstAndSecondPerson(username, "Kai", &sentence)

			//////////////////////////////////////////////////////////////////////
			// 2. only search if we can find enough tokens (more than one)
			{
				if db_model.GetNumSearchTokens(sentence.TokenList) > 1 {

					// 3. perform an index search in the factoid system
					rs, err := db_model.FindText(sentence.TokenList, username)
					if err != nil {
						ATJsonError(w, err.Error())
						return
					}
					// 4. if we cannot find any results for the user, go global
					if len(rs.ResultList) == 0 {
						rs, err = db_model.FindText(sentence.TokenList, "global")
						if err != nil {
							ATJsonError(w, err.Error())
							return
						}
					}
					// append results to return set
					for _, item := range rs.ResultList {
						ask_teach_result.ResultList = append(ask_teach_result.ResultList, item)
					}
				}
			}

			//////////////////////////////////////////////////////////////////////
			// 3. search Freebase
			{
				// todo: pagination, page / page_size put it through the uri
				tuple_set, err := freebase.FreebaseQueryBySearch(sentence.TokenList, 0, 10)
				if err == nil {
					for _, tuple := range tuple_set {
						item := model.ATResult{Text: tuple.String(), Topic: "K/AI", Timestamp: util.GetTimeNowSting()}
						ask_teach_result.ResultList = append(ask_teach_result.ResultList, item)
					}
				}
			}

			// write the return result
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			if len(ask_teach_result.ResultList) == 0 {  // nothing?
				ask_teach_result.ResultList = append(ask_teach_result.ResultList,
					model.ATResult{Text: "Sorry, I don't know.",
						Topic: "K/AI", Timestamp: util.GetTimeNowSting()})
			}
			json_bytes, _ := json.Marshal(ask_teach_result)
			w.Write(json_bytes)
		}

	} else {
		ATJsonError(w, "Query empty or too large, invalid")
	}
}

