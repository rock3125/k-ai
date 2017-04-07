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
	"testing"
	"k-ai/db"
	"k-ai/nlu/model"
	"encoding/json"
	"k-ai/util"
	"github.com/gocql/gocql"
	"k-ai/db/db_model"
	"k-ai/util_ut"
)

// test fn
func isTrue(t *testing.T, cond bool) {
	if !cond {
		t.Error("condition failed")
		panic("condition failed")
	}
}

// test the tuple tree parsing content
func TestAsk1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ask_text")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ask_text", 1)

	// This is some test text that will be indexed.test
	sentenceListStr := `[{"tokenList":[{"index":0,"list":[1],"tag":"DT","text":"This","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[4,1],"tag":"DT","text":"some","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[4,1],"tag":"NN","text":"test","dep":"compound","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":"NN","text":"text","dep":"attr","synid":-1,"semantic":""},{"index":5,"list":[8,4,1],"tag":"WDT","text":"that","dep":"nsubjpass","synid":-1,"semantic":""},{"index":6,"list":[8,4,1],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":""},{"index":7,"list":[8,4,1],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":8,"list":[4,1],"tag":"VBN","text":"indexed","dep":"relcl","synid":-1,"semantic":""},{"index":9,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	var sentenceList []model.Sentence
	err := json.Unmarshal([]byte(sentenceListStr), &sentenceList)
	util_ut.Check(t, err)

	id, err := gocql.RandomUUID()
	util_ut.Check(t, err)
	for i, _ := range sentenceList {
		util.CopyUUID(&sentenceList[i].Id, &id)
	}

	err = db_model.SaveText(sentenceList, "unit test")
	util_ut.Check(t, err)

	err = db_model.IndexText("unit test",0, sentenceList, 1.0)
	util_ut.Check(t, err)


	// do an index search:  What will be indexed?
	q_str := `[{"index":0,"list":[3],"tag":"WP","text":"What","dep":"nsubjpass","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":"man"},{"index":2,"list":[3],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"VBN","text":"indexed","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[3],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]`
	var qTokenList []model.Token
	err = json.Unmarshal([]byte(q_str), &qTokenList)
	util_ut.Check(t, err)

	rs, err := db_model.FindText(qTokenList, "unit test")
	util_ut.Check(t, err)
	isTrue(t, len(rs.ResultList) == 1)


	db.DropKeyspace("localhost", "kai_ask_text")
}

