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
	"fmt"
	"k-ai/util"
	"k-ai/util_ut"
)

// I need a loan from the bank for a mortgage.
const bankTokenList = `[{"index":0,"list":[1],"tag":"PRP","text":"I","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBP","text":"need","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[3,1],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"NN","text":"loan","dep":"dobj","synid":-1,"semantic":"woman"},{"index":4,"list":[3,1],"tag":"IN","text":"from","dep":"prep","synid":-1,"semantic":""},{"index":5,"list":[6,4,3,1],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":6,"list":[4,3,1],"tag":"NN","text":"bank","dep":"pobj","synid":-1,"semantic":"container"},{"index":7,"list":[3,1],"tag":"IN","text":"for","dep":"prep","synid":-1,"semantic":""},{"index":8,"list":[9,7,3,1],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":9,"list":[7,3,1],"tag":"NN","text":"mortgage","dep":"pobj","synid":-1,"semantic":""},{"index":10,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]`
// My government is helping the community with money.
const govTokenList = `[{"index":0,"list":[1,3],"tag":"PRP$","text":"My","dep":"poss","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"NN","text":"government","dep":"nsubj","synid":-1,"semantic":""},{"index":2,"list":[3],"tag":"VBZ","text":"is","dep":"aux","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"VBG","text":"helping","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[5,3],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":5,"list":[3],"tag":"NN","text":"community","dep":"dobj","synid":-1,"semantic":""},{"index":6,"list":[3],"tag":"IN","text":"with","dep":"prep","synid":-1,"semantic":""},{"index":7,"list":[6,3],"tag":"NN","text":"money","dep":"pobj","synid":-1,"semantic":""},{"index":8,"list":[3],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]`

// test the topics are in order of the list
func checkMainTopic(t *testing.T, topicList TopicScores, topics...string) {
	if len(topicList) != len(topics) {
		err_str := "topics do not all match"
		t.Error(err_str)
		panic(err_str)
	}
	// 1. check the topics all exist as expected
	for i, topic := range topics {
		if topic != topicList[i].Topic {
			err_str := fmt.Sprintf("topic score mismatch for '%s'", topic)
			t.Error(err_str)
			panic(err_str)
		}
	}
}

// test topic picker
func TestTopicker1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_topic_test_1")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_topic_test_1", 1)

	// create some topic indexes first from two different sources
	bank_json, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/bank.json")
	util_ut.Check(t, err)
	bank_sentence_list := jsonToSentenceList(t, bank_json)
	if len(bank_sentence_list) == 0 {
		t.Error("bank json empty")
	}
	err = SaveTopic("bank", "some bank text", bank_sentence_list)
	util_ut.Check(t, err)

	government_json, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/government.json")
	util_ut.Check(t, err)
	gov_sentence_list := jsonToSentenceList(t, government_json)
	if len(gov_sentence_list) == 0 {
		t.Error("government json empty")
	}
	err = SaveTopic("government", "some government text", gov_sentence_list)
	util_ut.Check(t, err)

	// then see how well we match each topic
	bankTopics, err := GetTopTopics(jsonToTokenList(t, bankTokenList))
	util_ut.Check(t, err)
	checkMainTopic(t, bankTopics, "bank")

	govTopics, err := GetTopTopics(jsonToTokenList(t, govTokenList))
	util_ut.Check(t, err)
	checkMainTopic(t, govTopics, "government", "bank") // in order of importance


	db.DropKeyspace("localhost", "kai_ai_topic_test_1")
}

// test CRUD works
func TestTopicCRUD1(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_topic_test_2")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_topic_test_2", 1)

	// create two topics
	// create some topic indexes first from two different sources
	bank_json, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/bank.json")
	util_ut.Check(t, err)
	bank_text, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/bank.txt")
	util_ut.Check(t, err)
	bank_sentence_list := jsonToSentenceList(t, bank_json)
	if len(bank_sentence_list) == 0 {
		t.Error("bank json empty")
	}
	err = SaveTopic("banking", bank_text, bank_sentence_list)
	util_ut.Check(t, err)

	government_json, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/government.json")
	util_ut.Check(t, err)
	government_text, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/government.txt")
	util_ut.Check(t, err)
	gov_sentence_list := jsonToSentenceList(t, government_json)
	if len(gov_sentence_list) == 0 {
		t.Error("government json empty")
	}
	err = SaveTopic("government", government_text, gov_sentence_list)
	util_ut.Check(t, err)

	// get a list of the topics and check it
	topic_list, err := GetTopicList("", 10)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(topic_list) == 2)

	// delete the banking topic
	err = DeleteTopic("banking")
	util_ut.Check(t, err)

	// re-get list
	topic_list, err = GetTopicList("", 10)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(topic_list) == 1)

	// delete the banking topic
	err = DeleteTopic("government")
	util_ut.Check(t, err)

	// re-get list
	topic_list, err = GetTopicList("", 10)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(topic_list) == 0)

	db.DropKeyspace("localhost", "kai_ai_topic_test_2")
}

