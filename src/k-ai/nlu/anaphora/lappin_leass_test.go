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

package anaphora

import (
	"testing"
	"strings"
	"runtime/debug"
	"k-ai/nlu/model"
	"encoding/json"
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


// helper: check the anaphora are resolved as expected
func checkPronounReference(t *testing.T, sentence_list model.SentenceList, anaphora string, words...string) {
	for _, sentence := range sentence_list {
		for _, t_token := range sentence.TokenList {
			found := false // any of the words we're looking for?
			for _, word := range words {
				if strings.ToLower(t_token.Text) == word {
					found = true
				}
			}
			if found {  // if so - check it matches the token index for anaphora resolution
				if t_token.Anaphora != anaphora {
					t.Errorf("pronoun reference incorrect for word %s, expected anaphora resolved to %s, but found %s",
										t_token.Text, anaphora, t_token.Anaphora)
					debug.PrintStack()
					t.FailNow()
				}
			}
		}
	}
}

// test Lappin Leass algorithm - simple single sentence resolution
func TestLL1(t *testing.T) {
	util_ut.IsTrue(t, LL.n_back > 0)  // setup correctly
	util_ut.IsTrue(t, len(LL.pronoun_set) > 0)

	// John said he likes dogs
	const s_list_1 = `[{"tokenList":[{"index":0,"list":[1],"tag":"NNP","text":"John","dep":"nsubj","synid":-1,"semantic":"male"},{"index":1,"list":[],"tag":"VBD","text":"said","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[3,1],"tag":"PRP","text":"he","dep":"nsubj","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"VBZ","text":"likes","dep":"ccomp","synid":-1,"semantic":""},{"index":4,"list":[3,1],"tag":"NNS","text":"dogs","dep":"dobj","synid":-1,"semantic":"animal"},{"index":5,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentence_list_1 := jsonToSentenceList(t, s_list_1)
	util_ut.IsTrue(t, len(sentence_list_1) == 1)

	// check that we can detect the pronoun
	util_ut.IsTrue(t, LL.HasPronoun(sentence_list_1[0]))

	// resolve the pronoun
	util_ut.IsTrue(t, LL.ResolvePronouns(sentence_list_1) == 1)

	// check its the right one
	util_ut.IsTrue(t, sentence_list_1[0].TokenList[2].Anaphora == "John")
}


// test non-match with invalid semantic
func TestLL2(t *testing.T) {
	util_ut.IsTrue(t, LL.n_back > 0)  // setup correctly
	util_ut.IsTrue(t, len(LL.pronoun_set) > 0)

	// John said he likes dogs
	const s_list_1 = `[{"tokenList":[{"index":0,"list":[1],"tag":"NNP","text":"John","dep":"nsubj","synid":-1,"semantic":"location"},{"index":1,"list":[],"tag":"VBD","text":"said","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[3,1],"tag":"PRP","text":"he","dep":"nsubj","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"VBZ","text":"likes","dep":"ccomp","synid":-1,"semantic":""},{"index":4,"list":[3,1],"tag":"NNS","text":"dogs","dep":"dobj","synid":-1,"semantic":"animal"},{"index":5,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentence_list_1 := jsonToSentenceList(t, s_list_1)
	util_ut.IsTrue(t, len(sentence_list_1) == 1)

	// check that we can detect the pronoun
	util_ut.IsTrue(t, LL.HasPronoun(sentence_list_1[0]))

	// resolve the pronoun
	util_ut.IsTrue(t, LL.ResolvePronouns(sentence_list_1) == 0)

	// check its the right one
	util_ut.IsTrue(t, sentence_list_1[0].TokenList[2].Anaphora == "?")
}


// test matching multi sentences with ambiguities
func TestLL3(t *testing.T) {
	util_ut.IsTrue(t, LL.n_back > 0)  // setup correctly
	util_ut.IsTrue(t, len(LL.pronoun_set) > 0)

	// The old lady pulled her spectacles down and looked over them  about the room; then she put them up and looked out under them.
	// She seldom or never looked THROUGH them for so small a thing  as a boy; they were her state pair, the pride of her heart, and were
	// built for ”style,” not service—she could have seen through a pair  of stove-lids just as well.
	const s_list_1 = `[{"tokenList":[{"index":0,"list":[1,3,18],"tag":"DT","text":"The","dep":"det","synid":-1,"semantic":""},{"index":2,"list":[3,18],"tag":"NN","text":"old lady","dep":"nsubj","synid":-1,"semantic":"female"},{"index":3,"list":[18],"tag":"VBD","text":"pulled","dep":"ccomp","synid":-1,"semantic":""},{"index":4,"list":[5,3,18],"tag":"PRP$","text":"her","dep":"poss","synid":-1,"semantic":""},{"index":5,"list":[3,18],"tag":"NNS","text":"spectacles","dep":"dobj","synid":-1,"semantic":""},{"index":6,"list":[3,18],"tag":"RB","text":"down","dep":"advmod","synid":-1,"semantic":""},{"index":7,"list":[3,18],"tag":"CC","text":"and","dep":"cc","synid":-1,"semantic":""},{"index":8,"list":[3,18],"tag":"VBD","text":"looked","dep":"conj","synid":-1,"semantic":""},{"index":9,"list":[8,3,18],"tag":"IN","text":"over","dep":"prep","synid":-1,"semantic":""},{"index":10,"list":[9,8,3,18],"tag":"PRP","text":"them","dep":"pobj","synid":-1,"semantic":""},{"index":12,"list":[8,3,18],"tag":"IN","text":"about","dep":"prep","synid":-1,"semantic":""},{"index":13,"list":[14,12,8,3,18],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":14,"list":[12,8,3,18],"tag":"NN","text":"room","dep":"pobj","synid":-1,"semantic":"location"},{"index":15,"list":[18],"tag":":","text":";","dep":"punct","synid":-1,"semantic":""},{"index":16,"list":[18],"tag":"RB","text":"then","dep":"advmod","synid":-1,"semantic":""},{"index":17,"list":[18],"tag":"PRP","text":"she","dep":"nsubj","synid":-1,"semantic":""},{"index":18,"list":[],"tag":"VBD","text":"put","dep":"ROOT","synid":-1,"semantic":""},{"index":19,"list":[18],"tag":"PRP","text":"them","dep":"dobj","synid":-1,"semantic":""},{"index":20,"list":[18],"tag":"RP","text":"up","dep":"prt","synid":-1,"semantic":""},{"index":21,"list":[18],"tag":"CC","text":"and","dep":"cc","synid":-1,"semantic":""},{"index":22,"list":[18],"tag":"VBD","text":"looked","dep":"conj","synid":-1,"semantic":""},{"index":23,"list":[22,18],"tag":"RP","text":"out","dep":"prt","synid":-1,"semantic":""},{"index":24,"list":[22,18],"tag":"IN","text":"under","dep":"prep","synid":-1,"semantic":""},{"index":25,"list":[24,22,18],"tag":"PRP","text":"them","dep":"pobj","synid":-1,"semantic":""},{"index":26,"list":[18],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""},{"index":27,"list":[26,18],"tag":"SP","text":"\n\t","dep":"","synid":-1,"semantic":""}]},{"tokenList":[{"index":28,"list":[32,46],"tag":"PRP","text":"She","dep":"nsubj","synid":-1,"semantic":""},{"index":29,"list":[32,46],"tag":"RB","text":"seldom","dep":"advmod","synid":-1,"semantic":""},{"index":30,"list":[29,32,46],"tag":"CC","text":"or","dep":"cc","synid":-1,"semantic":""},{"index":31,"list":[32,46],"tag":"RB","text":"never","dep":"neg","synid":-1,"semantic":""},{"index":32,"list":[46],"tag":"VBD","text":"looked","dep":"ccomp","synid":-1,"semantic":""},{"index":33,"list":[32,46],"tag":"IN","text":"THROUGH","dep":"prep","synid":-1,"semantic":""},{"index":34,"list":[33,32,46],"tag":"PRP","text":"them","dep":"pobj","synid":-1,"semantic":""},{"index":35,"list":[32,46],"tag":"IN","text":"for","dep":"prep","synid":-1,"semantic":""},{"index":36,"list":[37,39,35,32,46],"tag":"RB","text":"so","dep":"advmod","synid":-1,"semantic":""},{"index":37,"list":[39,35,32,46],"tag":"JJ","text":"small","dep":"amod","synid":-1,"semantic":""},{"index":38,"list":[39,35,32,46],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":39,"list":[35,32,46],"tag":"NN","text":"thing","dep":"pobj","synid":-1,"semantic":""},{"index":41,"list":[39,35,32,46],"tag":"IN","text":"as","dep":"prep","synid":-1,"semantic":""},{"index":42,"list":[43,41,39,35,32,46],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":43,"list":[41,39,35,32,46],"tag":"NN","text":"boy","dep":"pobj","synid":-1,"semantic":"male"},{"index":44,"list":[46],"tag":":","text":";","dep":"punct","synid":-1,"semantic":""},{"index":45,"list":[46],"tag":"PRP","text":"they","dep":"nsubj","synid":-1,"semantic":""},{"index":46,"list":[],"tag":"VBD","text":"were","dep":"ROOT","synid":-1,"semantic":""},{"index":47,"list":[49,46],"tag":"PRP$","text":"her","dep":"poss","synid":-1,"semantic":""},{"index":48,"list":[49,46],"tag":"NN","text":"state","dep":"compound","synid":-1,"semantic":""},{"index":49,"list":[46],"tag":"NN","text":"pair","dep":"attr","synid":-1,"semantic":""},{"index":50,"list":[49,46],"tag":",","text":",","dep":"punct","synid":-1,"semantic":""},{"index":51,"list":[52,49,46],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":52,"list":[49,46],"tag":"NN","text":"pride","dep":"conj","synid":-1,"semantic":""},{"index":53,"list":[52,49,46],"tag":"IN","text":"of","dep":"prep","synid":-1,"semantic":""},{"index":54,"list":[55,53,52,49,46],"tag":"PRP$","text":"her","dep":"poss","synid":-1,"semantic":""},{"index":55,"list":[53,52,49,46],"tag":"NN","text":"heart","dep":"pobj","synid":-1,"semantic":""},{"index":56,"list":[52,49,46],"tag":",","text":",","dep":"punct","synid":-1,"semantic":""},{"index":57,"list":[46],"tag":"CC","text":"and","dep":"cc","synid":-1,"semantic":""},{"index":58,"list":[60,46],"tag":"VBD","text":"were","dep":"auxpass","synid":-1,"semantic":""},{"index":59,"list":[58,60,46],"tag":"SP","text":"\n\t","dep":"","synid":-1,"semantic":""},{"index":60,"list":[46],"tag":"VBN","text":"built","dep":"conj","synid":-1,"semantic":""},{"index":61,"list":[60,46],"tag":"IN","text":"for","dep":"prep","synid":-1,"semantic":""},{"index":62,"list":[63,61,60,46],"tag":"JJ","text":"”","dep":"amod","synid":-1,"semantic":""},{"index":63,"list":[61,60,46],"tag":"NN","text":"style","dep":"pobj","synid":-1,"semantic":""},{"index":64,"list":[60,46],"tag":",","text":",","dep":"punct","synid":-1,"semantic":""},{"index":65,"list":[60,46],"tag":"NFP","text":"”","dep":"punct","synid":-1,"semantic":""},{"index":66,"list":[67,72,60,46],"tag":"RB","text":"not","dep":"neg","synid":-1,"semantic":""},{"index":67,"list":[72,60,46],"tag":"NN","text":"service","dep":"oprd","synid":-1,"semantic":""},{"index":68,"list":[72,60,46],"tag":"XX","text":"—","dep":"dep","synid":-1,"semantic":""},{"index":69,"list":[72,60,46],"tag":"PRP","text":"she","dep":"nsubj","synid":-1,"semantic":""},{"index":70,"list":[72,60,46],"tag":"MD","text":"could","dep":"aux","synid":-1,"semantic":""},{"index":71,"list":[72,60,46],"tag":"VB","text":"have","dep":"aux","synid":-1,"semantic":""},{"index":72,"list":[60,46],"tag":"VBN","text":"seen","dep":"conj","synid":-1,"semantic":""},{"index":73,"list":[72,60,46],"tag":"IN","text":"through","dep":"prep","synid":-1,"semantic":""},{"index":74,"list":[75,73,72,60,46],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":75,"list":[73,72,60,46],"tag":"NN","text":"pair","dep":"pobj","synid":-1,"semantic":""},{"index":77,"list":[75,73,72,60,46],"tag":"IN","text":"of","dep":"prep","synid":-1,"semantic":""},{"index":78,"list":[80,77,75,73,72,60,46],"tag":"NN","text":"stove","dep":"compound","synid":-1,"semantic":""},{"index":79,"list":[80,77,75,73,72,60,46],"tag":"HYPH","text":"-","dep":"punct","synid":-1,"semantic":""},{"index":80,"list":[77,75,73,72,60,46],"tag":"NNS","text":"lids","dep":"pobj","synid":-1,"semantic":""},{"index":81,"list":[83,72,60,46],"tag":"RB","text":"just as","dep":"advmod","synid":-1,"semantic":""},{"index":83,"list":[72,60,46],"tag":"RB","text":"well","dep":"advmod","synid":-1,"semantic":""},{"index":84,"list":[46],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	sentence_list_1 := jsonToSentenceList(t, s_list_1)
	util_ut.IsTrue(t, len(sentence_list_1) == 2)  // has two sentences

	// check that we can detect the pronoun
	util_ut.IsTrue(t, LL.HasPronoun(sentence_list_1[0]))
	util_ut.IsTrue(t, LL.HasPronoun(sentence_list_1[1]))

	// resolve the 11 pronouns
	pn_count := LL.ResolvePronouns(sentence_list_1)
	util_ut.IsTrue(t, pn_count == 11)

	// check the references are correct;  she and her refers to token 2 (old lady), and them and they refer to token 5 (spectacles)
	checkPronounReference(t, sentence_list_1, "old lady", "she", "her")
	checkPronounReference(t, sentence_list_1, "spectacles", "them", "they")
}


