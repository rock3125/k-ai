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
	"k-ai/nlu/model"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"k-ai/util_ut"
)

// from json str back to sentence list, unit testing aid
func jsonToSentenceList(t *testing.T, str string) []model.Sentence {
	var sentence_list []model.Sentence
	err := json.Unmarshal([]byte(str), &sentence_list)
	util_ut.Check(t, err)
	// setup ids for each sentence
	for i, _ := range sentence_list {
		sentence_list[i].RandomId()
	}
	return sentence_list
}

// parsed strings for unit testing
const someTextToIndexNewYork = `[{"tokenList":[{"index":0,"list":[1],"tag":"DT","text":"Some","dep":"det","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"IN","text":"to","dep":"prep","synid":-1,"semantic":""},{"index":3,"list":[4,2,1],"tag":"NN","text":"index","dep":"compound","synid":-1,"semantic":""},{"index":4,"list":[4,2,1],"tag":"NNP","text":"New York","dep":"compound","synid":-1,"semantic":"state"},{"index":6,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
const someOtherText = `[{"tokenList":[{"index":0,"list":[2],"tag":"DT","text":"Some other","dep":"det","synid":-1,"semantic":""},{"index":2,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""}]}]`
const someTextToIndexNewYorkOrJapan = `[{"tokenList":[{"index":0,"list":[1],"tag":"DT","text":"Some","dep":"det","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"IN","text":"to","dep":"prep","synid":-1,"semantic":""},{"index":3,"list":[2,1],"tag":"NN","text":"index","dep":"pobj","synid":-1,"semantic":""},{"index":4,"list":[4,3,2,1],"tag":"NNP","text":"New York","dep":"compound","synid":-1,"semantic":"state"},{"index":6,"list":[4,3,2,1],"tag":"CC","text":"or","dep":"cc","synid":-1,"semantic":""},{"index":7,"list":[4,3,2,1],"tag":"NNP","text":"Japan","dep":"conj","synid":-1,"semantic":"country"},{"index":8,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
const someOtherTextInJapan = `[{"tokenList":[{"index":0,"list":[2],"tag":"DT","text":"Some other","dep":"det","synid":-1,"semantic":""},{"index":2,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""},{"index":3,"list":[2],"tag":"IN","text":"in","dep":"prep","synid":-1,"semantic":""},{"index":4,"list":[3,2],"tag":"NNP","text":"Japan","dep":"pobj","synid":-1,"semantic":"country"}]}]`
const tokenListText = `[{"index":0,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""}]`
const tokenListJapanText = `[{"index":0,"list":[1],"tag":"NNP","text":"Japan","dep":"compound","synid":-1,"semantic":"country"},{"index":1,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""}]`

// "what country indexes text?"
const whatCountryTokenList = `[{"index":0,"list":[3],"tag":"WP","text":"what","dep":"dobj","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"NN","text":"country","dep":"nsubj","synid":-1,"semantic":""},{"index":2,"list":[3],"tag":"NNS","text":"indexes","dep":"compound","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"NN","text":"text","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[3],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`

// Peter was working from home.
const peterWasWorkingFromHome = `[{"tokenList":[{"index":0,"list":[2],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[2],"tag":"VBD","text":"was","dep":"aux","synid":-1,"semantic":""},{"index":2,"list":[],"tag":"VBG","text":"working","dep":"ROOT","synid":-1,"semantic":"location"},{"index":3,"list":[2],"tag":"IN","text":"from","dep":"prep","synid":-1,"semantic":""},{"index":4,"list":[3,2],"tag":"NN","text":"home","dep":"pobj","synid":-1,"semantic":"location"},{"index":5,"list":[2],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`

// verb "be" in tokenlist
const beTokenList = `[{"index":0,"list":[],"tag":"VB","text":"is","dep":"","synid":-1,"semantic":""}]`
const workingTokenList = `[{"index":0,"list":[],"tag":"VB","text":"work","dep":"","synid":-1,"semantic":""}]`

// the word Japan
const tokenListJapan = `[{"index":0,"list":[],"tag":"NNP","text":"Japan","dep":"","synid":-1,"semantic":"country"}]`

// the word New York
const tokenListNewYork = `[{"index":0,"list":[],"tag":"NNP","text":"New York","dep":"","synid":-1,"semantic":"state"}]`


// from json str back to token list
func jsonToTokenList(t *testing.T, str string) []model.Token {
	var token_list []model.Token
	err := json.Unmarshal([]byte(str), &token_list)
	util_ut.Check(t, err)
	return token_list
}

// check a result set contains the set of urls as provided
func contains(t *testing.T, set map[gocql.UUID][]model.IndexMatch, url_list...gocql.UUID) {
	for _, url := range url_list {
		if _, ok := set[url]; !ok {
			err_str := fmt.Sprintf("error: map does not contain url(%s)", url)
			t.Error(err_str)
			panic(err_str)
		}
	}
}

// perform index tests
func TestIndexer1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test", 1)

	// index some text
	err := IndexText("topic1", 0, jsonToSentenceList(t, someTextToIndexNewYork), 1.0)
	util_ut.Check(t, err)

	err = IndexText("topic2", 0, jsonToSentenceList(t, someOtherText), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected
	index_list, err := readIndexes("text", "topic1",  0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_list) == 1)
	util_ut.IsTrue(t, index_list[0].Word == "text")
	util_ut.IsTrue(t, index_list[0].Topic == "topic1")
	util_ut.IsTrue(t, index_list[0].Offset == 1)

	index_list, err = readIndexes("York", "topic1", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_list) == 1)
	util_ut.IsTrue(t, index_list[0].Topic == "topic1")
	// important: "New York" is treated as one entity with lexicon loaded
	util_ut.IsTrue(t, index_list[0].Offset == 4)

	index_list, err = readIndexes("text", "topic2", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_list) == 1)
	util_ut.IsTrue(t, index_list[0].Topic == "topic2")
	util_ut.IsTrue(t, index_list[0].Offset == 1)

	db.DropKeyspace("localhost", "kai_ai_index_test")
}

// perform further index tests
func TestIndexer2(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_2")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_2", 1)

	// index some text
	err := IndexText("topic2", 0, jsonToSentenceList(t, someTextToIndexNewYork), 1.0)
	util_ut.Check(t, err)

	err = IndexText("topic2", 0, jsonToSentenceList(t, someOtherText), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected
	index_list, err := readIndexes("text", "topic2", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_list) == 2)
	unique_urls := make(map[gocql.UUID]bool,0)
	for _, index := range index_list{
	  	util_ut.IsTrue(t, index.Word == "text")
		util_ut.IsTrue(t, index.Topic == "topic2")
		unique_urls[index.Sentence_id] = true
	}
	util_ut.IsTrue(t, len(unique_urls) == 2)

	db.DropKeyspace("localhost", "kai_ai_index_test_2")
}


// perform further index tests single keyword
func TestIndexer3(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_3")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_3", 1)

	// index some text
	err := IndexText("topic3", 0,jsonToSentenceList(t, someTextToIndexNewYork), 1.0)
	util_ut.Check(t, err)

	err = IndexText("topic3", 0,jsonToSentenceList(t, someOtherText), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected
	index_map, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListText), "topic3",0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 2)

	db.DropKeyspace("localhost", "kai_ai_index_test_3")
}


// perform further index tests multiple keyword
func TestIndexer4(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_3")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_3", 1)

	// index some text
	err := IndexText("topic4",0, jsonToSentenceList(t, someTextToIndexNewYorkOrJapan), 1.0)
	util_ut.Check(t, err)

	err = IndexText("topic4",0, jsonToSentenceList(t, someOtherTextInJapan), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected
	index_map, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListJapanText), "topic4", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 2)
	for _, index_list := range index_map {
		if len(index_list) != 2 {
			t.Error("expected two indexes per match (two keywords)")
		}
	}

	db.DropKeyspace("localhost", "kai_ai_index_test_3")
}


// test origin isolation
func TestIndexer5(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_4")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_4", 1)

	// index some text
	err := IndexText( "topic5",0, jsonToSentenceList(t, someTextToIndexNewYorkOrJapan), 1.0)
	util_ut.Check(t, err)

	err = IndexText( "topic6",0, jsonToSentenceList(t, someOtherTextInJapan), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected
	index_map, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListJapanText), "topic5", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 1)
	for _, index_list := range index_map {
		if len(index_list) != 2 {
			t.Error("expected only one index")
		}
	}

	db.DropKeyspace("localhost", "kai_ai_index_test_4")
}


// test semantic indexes (Japan: country, New York: state)
func TestIndexer6(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_5")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_5", 1)

	// index: "Some text to index New York or Japan.
	err := IndexText("topic6",0, jsonToSentenceList(t, someTextToIndexNewYorkOrJapan), 1.0)
	util_ut.Check(t, err)

	// use the index system to find items as we'd expected: "what country indexes text?"
	index_map, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, whatCountryTokenList), "topic6", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 1)
	for url, index_list := range index_map {
		if len(index_list) != 3 {
			t.Error("expected three hits for sentence: " + url.String())
		}
	}

	db.DropKeyspace("localhost", "kai_ai_index_test_5")
}


// test auxiliary verb indexing (Peter was working from home.) don't index the auxiliary verb in this case
func TestAuxiliaryVerbs1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_index_test_6")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_index_test_6", 1)

	// index: "Some text to index New York or Japan.
	sl1 := jsonToSentenceList(t, peterWasWorkingFromHome)
	util_ut.IsTrue(t, len(sl1) == 1)
	err := IndexText("topic7",0, sl1, 1.0)
	util_ut.Check(t, err)

	// make sure this wasn't indexed under its AUX tag
	index_map, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, beTokenList), "topic7",0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 0)

	// but WAS indexed under the verb "work"
	index_map_2, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, workingTokenList), "topic7", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_2) == 1)
	// one offset matches "working" @ 2
	util_ut.IsTrue(t, index_map_2[sl1[0].Id][0].Offset == 2)

	db.DropKeyspace("localhost", "kai_ai_index_test_6")
}


// test unindexing works
func TestUnIndexer1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_unindex_test_1")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_unindex_test_1", 1)

	// index some text
	sl1 := jsonToSentenceList(t, someTextToIndexNewYorkOrJapan)
	util_ut.IsTrue(t, len(sl1) == 1)
	err := IndexText( "topic8", 0, sl1, 1.0)
	util_ut.Check(t, err)

	sl2 := jsonToSentenceList(t, someTextToIndexNewYorkOrJapan)
	util_ut.IsTrue(t, len(sl2) == 1)
	err = IndexText("topic8", 0, sl2, 1.0)
	util_ut.Check(t, err)

	// make sure there are two urls for Japan
	index_map_1, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListJapan), "topic8", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_1) == 2)
	contains(t, index_map_1, sl1[0].Id, sl2[0].Id)
	// make sure there are two urls for New York
	index_map_2, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListNewYork), "topic8", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_2) == 2)
	contains(t, index_map_2, sl1[0].Id, sl2[0].Id)

	// we remove url1, halving the number of hits
	err = RemoveIndexes(sl1[0].Id, "topic8")
	util_ut.Check(t, err)

	// make sure there are two urls for Japan
	index_map_3, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListJapan), "topic8", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_3) == 1)
	contains(t, index_map_3, sl2[0].Id)
	// make sure there are two urls for New York
	index_map_4, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListNewYork), "topic8",0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_4) == 1)
	contains(t, index_map_4, sl2[0].Id)

	// we remove url2 => no more hits
	err = RemoveIndexes(sl2[0].Id, "topic8")
	util_ut.Check(t, err)
	index_map_5, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListJapan), "topic8", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_5) == 0)
	index_map_6, err := ReadIndexesWithFilterForTokens(jsonToTokenList(t, tokenListNewYork), "topic8",0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map_6) == 0)

	db.DropKeyspace("localhost", "kai_ai_unindex_test_1")
}

// test "you" gets replaced with Kai and its semantic from the lexicon (ai)
func TestPronounReplacement1(t *testing.T) {

	{
		// what are you?
		const str1= `[{"tokenList":[{"index":0,"list":[1],"tag":"WP","text":"What","dep":"attr","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBP","text":"are","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"PRP","text":"you","dep":"nsubj","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]}]`
		sentence_list1 := jsonToSentenceList(t, str1)
		util_ut.IsTrue(t, len(sentence_list1) == 1)
		ResolveFirstAndSecondPerson("Peter", "Kai", &sentence_list1[0])
		token_list1 := sentence_list1[0].TokenList
		util_ut.IsTrue(t, token_list1[4].Text == "Kai")
		util_ut.IsTrue(t, token_list1[4].Tag == "NNP")
		util_ut.IsTrue(t, token_list1[4].Semantic == "ai")
		util_ut.IsTrue(t, token_list1[4].Index == 2)
	}

	{
		// I am Kai, an artificial intelligence.
		const str2= `[{"tokenList":[{"index":0,"list":[1],"tag":"PRP","text":"I","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBP","text":"am","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"NNP","text":"Kai","dep":"attr","synid":-1,"semantic":"ai"},{"index":3,"list":[2,1],"tag":",","text":",","dep":"punct","synid":-1,"semantic":""},{"index":4,"list":[5,2,1],"tag":"DT","text":"an","dep":"det","synid":-1,"semantic":""},{"index":5,"list":[5,2,1],"tag":"JJ","text":"artificial intelligence","dep":"amod","synid":-1,"semantic":""},{"index":7,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
		sentence_list2 := jsonToSentenceList(t, str2)
		util_ut.IsTrue(t, len(sentence_list2) == 1)
		ResolveFirstAndSecondPerson("Kai", "Peter", &sentence_list2[0])
		token_list2 := sentence_list2[0].TokenList
		util_ut.IsTrue(t, token_list2[2].Text == "Kai")
		util_ut.IsTrue(t, token_list2[2].Tag == "NNP")
		util_ut.IsTrue(t, token_list2[2].Semantic == "ai")
		util_ut.IsTrue(t, token_list2[2].Index == 0)
	}

	{
		// Peter likes point (i) the best.
		const str3= `[{"topic":"","tokenList":[{"index":0,"list":[1],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"likes","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"NN","text":"point","dep":"dobj","synid":-1,"semantic":"location"},{"index":3,"list":[2,1],"tag":"-LRB-","text":"(","dep":"punct","synid":-1,"semantic":""},{"index":4,"list":[2,1],"tag":"NN","text":"i","dep":"appos","synid":-1,"semantic":""},{"index":5,"list":[2,1],"tag":"-RRB-","text":")","dep":"punct","synid":-1,"semantic":""},{"index":6,"list":[7,1],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":7,"list":[1],"tag":"JJS","text":"best","dep":"npadvmod","synid":-1,"semantic":"person"},{"index":8,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
		sentence_list3 := jsonToSentenceList(t, str3)
		util_ut.IsTrue(t, len(sentence_list3) == 1)
		ResolveFirstAndSecondPerson("Kai", "Peter", &sentence_list3[0])
		token_list3 := sentence_list3[0].TokenList
		util_ut.IsTrue(t, token_list3[4].Text == "i")
		util_ut.IsTrue(t, token_list3[4].Tag == "NN")
		util_ut.IsTrue(t, token_list3[5].Text == ")")
		util_ut.IsTrue(t, token_list3[5].Tag == "-RRB-")
	}

}

