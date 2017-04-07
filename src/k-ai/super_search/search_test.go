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

package super_search

import (
	"testing"
	"k-ai/nlu/model"
	"encoding/json"
	"k-ai/db/db_model"
	"fmt"
	"github.com/gocql/gocql"
	"k-ai/util_ut"
)


// Peter and Davis Patrick went to the beach at 12:45 to view the boats coming in.
const str1 = `[{"tokenList":[{"index":0,"list":[4],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[0,4],"tag":"CC","text":"and","dep":"cc","synid":-1,"semantic":""},{"index":2,"list":[3,0,4],"tag":"NNP","text":"Davis","dep":"compound","synid":-1,"semantic":"person"},{"index":3,"list":[0,4],"tag":"NNP","text":"Patrick","dep":"conj","synid":-1,"semantic":"person"},{"index":4,"list":[],"tag":"VBD","text":"went","dep":"ROOT","synid":-1,"semantic":""},{"index":5,"list":[4],"tag":"IN","text":"to","dep":"prep","synid":-1,"semantic":""},{"index":6,"list":[7,5,4],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":7,"list":[5,4],"tag":"NN","text":"beach","dep":"pobj","synid":-1,"semantic":"location"},{"index":8,"list":[4],"tag":"IN","text":"at","dep":"prep","synid":-1,"semantic":""},{"index":9,"list":[8,4],"tag":"CD","text":"12:45","dep":"pobj","synid":-1,"semantic":""},{"index":10,"list":[11,4],"tag":"TO","text":"to","dep":"aux","synid":-1,"semantic":""},{"index":11,"list":[4],"tag":"VB","text":"view","dep":"advcl","synid":-1,"semantic":""},{"index":12,"list":[13,11,4],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":13,"list":[11,4],"tag":"NNS","text":"boats","dep":"dobj","synid":-1,"semantic":"vehicle"},{"index":14,"list":[4],"tag":"VBG","text":"coming","dep":"advcl","synid":-1,"semantic":""},{"index":15,"list":[14,4],"tag":"RB","text":"in","dep":"advmod","synid":-1,"semantic":""},{"index":16,"list":[4],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`

// Peter and Sherry went to the mall to watch a movie.
const str2 = `[{"tokenList":[{"index":0,"list":[3],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[0,3],"tag":"CC","text":"and","dep":"cc","synid":-1,"semantic":""},{"index":2,"list":[0,3],"tag":"NNP","text":"Sherry","dep":"conj","synid":-1,"semantic":"woman"},{"index":3,"list":[],"tag":"VBD","text":"went","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[3],"tag":"IN","text":"to","dep":"prep","synid":-1,"semantic":""},{"index":5,"list":[6,4,3],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":6,"list":[4,3],"tag":"NN","text":"mall","dep":"pobj","synid":-1,"semantic":""},{"index":7,"list":[8,3],"tag":"TO","text":"to","dep":"aux","synid":-1,"semantic":""},{"index":8,"list":[3],"tag":"VB","text":"watch","dep":"advcl","synid":-1,"semantic":""},{"index":9,"list":[10,8,3],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":10,"list":[8,3],"tag":"NN","text":"movie","dep":"dobj","synid":-1,"semantic":""},{"index":11,"list":[3],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`

// from json str back to sentence list and setup an ID for this sentence
func jsonToSentenceList(t *testing.T, str string) []model.Sentence {
	var sentence_list []model.Sentence
	err := json.Unmarshal([]byte(str), &sentence_list)
	util_ut.Check(t, err)
	for i, sentence := range sentence_list {
		sentence.RandomId()
		sentence_list[i] = sentence
	}
	return sentence_list
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


// test the super search system
func TestSS1(t *testing.T) {

	db_model.Delete_and_create_keyspace_for_unit_test("kai_ai_ss1")

	origin := "peter@peter.co.nz"

	sentence_list_1 := jsonToSentenceList(t, str1)
	sentence_list_2 := jsonToSentenceList(t, str2)

	// index this knowledge
	db_model.IndexText(origin, 0, sentence_list_1, 1.0)
	db_model.IndexText(origin, 0, sentence_list_2, 1.0)

	// perform the super searches for testing
	rs1, err := SuperSearch("any(Peter)", origin)
	util_ut.Check(t, err)
	isTrue(t, rs1 != nil && len(rs1) == 2)
	contains(t, rs1, sentence_list_1[0].Id, sentence_list_2[0].Id)

	// Peter exists in boy, boat only in url1
	rs2, err := SuperSearch("any(Peter) and any(boat)", origin)
	util_ut.Check(t, err)
	isTrue(t, rs2 != nil && len(rs2) == 1)
	contains(t, rs1, sentence_list_1[0].Id)

	// union test
	rs3, err := SuperSearch("any(boat) or any(movie)", origin)
	util_ut.Check(t, err)
	isTrue(t, rs3 != nil && len(rs3) == 2)
	contains(t, rs1, sentence_list_1[0].Id, sentence_list_2[0].Id)

	// and not test
	rs4, err := SuperSearch("any(Peter) and not any(movie)", origin)
	util_ut.Check(t, err)
	isTrue(t, rs4 != nil && len(rs4) == 1)
	contains(t, rs1, sentence_list_1[0].Id)


	db_model.Delete_keyspace_after_unit_test("kai_ai_ss1")
}



