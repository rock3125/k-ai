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
	"strings"
	"github.com/gocql/gocql"
	"encoding/json"
	"time"
	"strconv"
	"k-ai/nlu/model"
	"k-ai/db/db_model"
	"k-ai/nlu/tokenizer"
	"k-ai/nlu/parser"
	"k-ai/util"
)

// replace magic {} values for AIML
func replaceMagicValues(binding model.AimlBinding, email string) string {
	text := binding.Text
	if strings.Contains(text, "{stats}") {
		statsStr := "system statistics: todo"
		text = strings.Replace(text, "{stats}", statsStr, -1)
	}
	if strings.Contains(text, "{year}") {
		text = strings.Replace(text, "{year}", strconv.Itoa(time.Now().Year()), -1)
	}
	if strings.Contains(text, "{day}") {
		text = strings.Replace(text, "{day}", time.Now().Weekday().String(), -1)
	}
	if strings.Contains(text, "{month}") {
		text = strings.Replace(text, "{month}", time.Now().Month().String(), -1)
	}
	if strings.Contains(text, "{time}") {
		text = strings.Replace(text, "{time}", time.Now().Format(time.Kitchen), -1)
	}
	if strings.Contains(text, "{date}") {
		text = strings.Replace(text, "{date}", time.Now().Format(time.RFC850), -1)
	}
	if strings.Contains(text, "{email}") {
		text = strings.Replace(text, "{email}", "your email address is " + email, -1)
	}
	if strings.Contains(text, "{name}") || strings.Contains(text, "{fullname}") {
		text = strings.Replace(text, "{name}", "I can't tell you your name but your email address is " + email, -1)
	}
	if strings.Contains(text, "{star}") {
		tstr := tokenizer.ToString(binding.TokenList)
		text = strings.Replace(text, "{star}", tstr, -1)
	}
	return text
}

// turn a set of index results into a series of text results
func addIndexResults(result_map map[gocql.UUID][]model.IndexMatch, result_list *model.ATResultList) {
	// load each url's object
	for id, index_list := range result_map {
		if len(index_list) > 0 {
			sentence, err := db_model.GetText(&id)
			if err == nil {
				str := sentence.ToString()
				result_list.ResultList = append(result_list.ResultList, model.ATResult{Text: str, Topic: sentence.Topic,
							Sentence_id: sentence.Id, Timestamp: util.GetTimeNowSting()})
			}
		}
	}
}


// perform special match characters on aiml matches
// if appropriate
func (mrg *AimlManager) PerformSpecialOps(binding_results []model.AimlBinding, topic string) (*model.ATResultList, error) {
	rs := model.ATResultList{ResultList: make([]model.ATResult,0)}
	schema_map, err := db_model.GetSchemaMap()
	if err != nil { return nil, err }

	for _, binding := range binding_results {

		////////////////////////////////////////////////////////////////////
		// db search?   the special entities from the kb system

		if strings.HasPrefix(binding.Text, "db_search:") {
			parts := strings.Split(binding.Text, ":")
			if len(parts) == 3 { // db_search: entity name : field name
				topic_str := parts[1]
				result_map, err := db_model.ReadIndexesWithFilterForTokens(binding.TokenList, topic_str, 0)
				if err != nil {
					return nil, err
				}
				if schema, ok := schema_map[parts[1]]; ok {

					// load each url's object
					for url, _ := range result_map {
						kb_entry, err := db_model.GetKBEntryById(&url, parts[1]) // load obj from db
						if err == nil {
							if kb_entry != nil && len(kb_entry.Json_data) > 0 {
								var itemMap map[string]string
								err = json.Unmarshal([]byte(kb_entry.Json_data), &itemMap)
								if err == nil {
									str := ""
									for _, field := range schema.Field_list {
										if value, ok := itemMap[field.Name]; ok {
											if len(str) > 0 {
												str += ", "
											}
											str += field.Name + ": " + value
										}
									}
									rs.ResultList = append(rs.ResultList, model.ATResult{Text: str,
										KB_id: kb_entry.Id, Topic: schema.Name, Timestamp: util.GetTimeNowSting() })
								}
							}
						}
					}

				}
			}

		} else if strings.HasPrefix(binding.Text, "{search:") && strings.HasSuffix(binding.Text, "}") {

			////////////////////////////////////////////////////////////////////
			// this is a special case of an AIML query that corrects a search pattern - so we build a query string
			// for the specified patterns and then do a proper search in the factoid system using that data

			// build a search string
			parts := strings.Split(binding.Text[8:len(binding.Text)-1], ",")
			search_str := ""
			for _, part := range parts {
				if part == "*" {
					for _, t_token := range binding.TokenList {
						if len(search_str) > 0 {
							search_str += " "
						}
						search_str += t_token.Text
					}
				} else {
					if len(search_str) > 0 {
						search_str += " "
					}
					search_str += part
				}
			}
			// parse the new search pattern
			sentence_list, err := parser.ParseText(search_str)
			if err == nil && len(sentence_list) > 0 {
				result_map, err := db_model.ReadIndexesWithFilterForTokens(sentence_list[0].TokenList, topic, 0)
				if err == nil {
					// build the search results into rs using the index map
					addIndexResults(result_map, &rs)
				}
			}

		} else { // default primitive AIML behavior

			////////////////////////////////////////////////////////////////////

			// replace any magic values such as {time} etc.
			binding_text:= replaceMagicValues(binding, topic)

			rs.ResultList = append(rs.ResultList, model.ATResult{Text: binding_text,
							Topic: binding.Origin, Timestamp: util.GetTimeNowSting() })
		}

	}
	return &rs, nil
}

