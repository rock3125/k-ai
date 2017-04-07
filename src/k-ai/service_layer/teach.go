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
	"k-ai/nlu/parser"
	"k-ai/db/db_model"
	"strings"
	"github.com/gorilla/mux"
	"fmt"
)

// perform a teach
//
func Teach(w http.ResponseWriter, r *http.Request) {
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
		db_model.AddLogEntry(username, "teach:" + bodyStr)

		// 1. parse it
		sentence_list, err := parser.ParseText(bodyStr)
		if err != nil {
			ATJsonError(w, "Unexpected parser error:" + err.Error())
			return
		}
		if len(sentence_list) == 0 {
			ATJsonError(w, "Please teach me something, this looks like an empty sentence.")
		} else if len(sentence_list) > 1 {
			ATJsonError(w, "Teach me using simple single sentences please.")
		} else {

			// replace I, me, myself pronoun references with KAI
			db_model.ResolveFirstAndSecondPerson(username, "Kai", &sentence_list[0])
			sentence := sentence_list[0] // the sentence

			if sentence.IsImperative() {
				ATJsonError(w, "That looks like a request or a command rather than information.")
			} else if sentence.IsQuestion() {
				ATJsonError(w, "That looks like a question, please use \"Ask me something\" for questions.")
			} else if !sentence.HasVerb() {
				ATJsonError(w, "I don't understand your statement, can you please change it?")
			} else if db_model.GetNumSearchTokens(sentence.TokenList) <= 1 {  // not enough information in the sentence
				ATJsonError(w, "There is something wrong with this sentence, Please rephrase it.")
			} else {
				sentence_list[0].RandomId()  // setup a guid for the new fact (random)
				err = db_model.SaveText(sentence_list, username)  // save a text factoid, and remove previous indexes
				if err != nil {
					ATJsonError(w, err.Error())
				} else {

					// index the text factoid
					err = db_model.IndexText(username, 0, sentence_list, 1.0)
					if err != nil {
						ATJsonError(w, err.Error()+"("+username+" indexes)")
					} else {

						// index this item in the global system too for finding across all items
						err = db_model.IndexText("global", 0, sentence_list, 1.0)
						if err != nil {
							ATJsonError(w, err.Error()+" (global indexes)")
						} else {

							guid_str := sentence_list[0].Id.String()
							ATJsonMessage(w, http.StatusAccepted, fmt.Sprintf("ok, got that and stored \"%s\" away as factoid \"%s\".", bodyStr, guid_str))
						}
					}

				} // if save text

			} // else if valid sentence

		} // else if sentences == 1

	} else {
		ATJsonError(w, "Teach text empty or too large.")
	}

}

