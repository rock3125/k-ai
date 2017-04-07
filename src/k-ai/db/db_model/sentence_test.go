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
	"k-ai/nlu/lexicon"
	"k-ai/util_ut"
)

// perform further index tests multiple keyword
func TestText1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_text")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_text", 1)

	// This is some test text that will be indexed.
	sentenceListStr := `[{"tokenList":[{"index":0,"list":[1],"tag":"DT","text":"This","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[4,1],"tag":"DT","text":"some","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[4,1],"tag":"NN","text":"test","dep":"compound","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":"NN","text":"text","dep":"attr","synid":-1,"semantic":""},{"index":5,"list":[8,4,1],"tag":"WDT","text":"that","dep":"nsubjpass","synid":-1,"semantic":""},{"index":6,"list":[8,4,1],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":""},{"index":7,"list":[8,4,1],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":8,"list":[4,1],"tag":"VBN","text":"indexed","dep":"relcl","synid":-1,"semantic":""},{"index":9,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentenceList := jsonToSentenceList(t, sentenceListStr)
	err := SaveText(sentenceList, "unit test")
	util_ut.Check(t, err)

	text_item, err := GetText(&sentenceList[0].Id)
	if err != nil || len(text_item.TokenList) != len(sentenceList[0].TokenList) {
		t.Error("invalid return")
	} else {
		err = DeleteText(&sentenceList[0].Id, "unit test")
		if err != nil {
			t.Error("delete failed")
		} else {
			text_item2, err := GetText(&sentenceList[0].Id)
			if err != nil || (text_item2 != nil && len(text_item2.TokenList) > 0) {
				t.Error("delete failed, get works")
			}
		}
	}

	db.DropKeyspace("localhost", "kai_ai_text")
}


// perform further index tests multiple keyword
func TestText2(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_text_2")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_text_2", 1)

	// This is some test text that will be indexed.  But this time its two sentences.
	s_str := `[{"tokenList":[{"index":0,"list":[1],"tag":"DT","text":"This","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[4,1],"tag":"DT","text":"some","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[4,1],"tag":"NN","text":"test","dep":"compound","synid":-1,"semantic":""},{"index":4,"list":[1],"tag":"NN","text":"text","dep":"attr","synid":-1,"semantic":""},{"index":5,"list":[8,4,1],"tag":"WDT","text":"that","dep":"nsubjpass","synid":-1,"semantic":""},{"index":6,"list":[8,4,1],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":""},{"index":7,"list":[8,4,1],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":8,"list":[4,1],"tag":"VBN","text":"indexed","dep":"relcl","synid":-1,"semantic":""},{"index":9,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]},{"tokenList":[{"index":11,"list":[16],"tag":"CC","text":"But","dep":"cc","synid":-1,"semantic":""},{"index":12,"list":[13,16],"tag":"DT","text":"this","dep":"det","synid":-1,"semantic":""},{"index":13,"list":[16],"tag":"NN","text":"time","dep":"nsubj","synid":-1,"semantic":""},{"index":14,"list":[16],"tag":"PRP$","text":"its","dep":"poss","synid":-1,"semantic":""},{"index":15,"list":[16],"tag":"CD","text":"two","dep":"nummod","synid":-1,"semantic":""},{"index":16,"list":[],"tag":"NNS","text":"sentences","dep":"ROOT","synid":-1,"semantic":""},{"index":17,"list":[16],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentenceList := jsonToSentenceList(t, s_str)

	err := SaveText(sentenceList, "unit test")
	util_ut.Check(t, err)

	textItem2, err := GetText(&sentenceList[0].Id)
	if err != nil || len(textItem2.TokenList) != len(sentenceList[0].TokenList) {
		t.Error("invalid return")
	} else {
		err = DeleteText(&sentenceList[0].Id, "unit test")
		if err != nil {
			t.Error("delete failed")
		} else {
			textItem3, err := GetText(&sentenceList[0].Id)
			if err != nil || (textItem3 != nil && len(textItem3.TokenList) != 0) {
				t.Error("delete failed, get works")
			}
		}
	}

	db.DropKeyspace("localhost", "kai_ai_text_2")
}

// text semantic indexing
func TestTextSemantics(t *testing.T) {
	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_sem_1")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_sem_1", 1)

	// smurfette eta smurf
	lexicon.Lexi.AddSemantic("smurfette", "smurf")

	// The smurfette is blue.
	const smurfetteBlueSL = `[{"tokenList":[{"index":0,"list":[1,2],"tag":"DT","text":"The","dep":"det","synid":-1,"semantic":""},{"index":1,"list":[2],"tag":"NN","text":"smurfette","dep":"nsubj","synid":-1,"semantic":"smurf"},{"index":2,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":3,"list":[2],"tag":"JJ","text":"blue","dep":"acomp","synid":-1,"semantic":""},{"index":4,"list":[2],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentence_list := jsonToSentenceList(t, smurfetteBlueSL)
	util_ut.IsTrue(t, len(sentence_list) == 1)
	err := IndexText("topic1", 0,sentence_list,1.0)
	util_ut.Check(t, err)

	// query: is a smurf blue?
	const smurfBlueSL = `[{"tokenList":[{"index":0,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":1,"list":[3,0],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":2,"list":[3,0],"tag":"NN","text":"smurf","dep":"compound","synid":-1,"semantic":""},{"index":3,"list":[0],"tag":"NN","text":"blue","dep":"attr","synid":-1,"semantic":""},{"index":4,"list":[0],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]}]`
	query_list := jsonToSentenceList(t, smurfBlueSL)
	index_map, err := ReadIndexesWithFilterForTokens(query_list[0].TokenList, "topic1", 0)
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(index_map) == 1)

	db.DropKeyspace("localhost", "kai_ai_sem_1")
}

