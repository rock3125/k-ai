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

package model

import (
	"testing"
	"encoding/json"
	"fmt"
	"k-ai/util_ut"
)

// from json str back to a single sentence (the first)
func jsonToSentence(t *testing.T, str string) Sentence {
	var sentence_list []Sentence
	err := json.Unmarshal([]byte(str), &sentence_list)
	util_ut.Check(t, err)
	return sentence_list[0]
}

// test fn
func isTrue(t *testing.T, cond bool) {
	if !cond {
		t.Error("condition failed")
		panic("condition failed")
	}
}

// test the tuple tree parsing content
func TestValidPath(t *testing.T) {
	// pre-parsed string: Craig has a boat in the harbour.
	sentenceListStr := `[{"tokenList":[{"index":0,"list":[1],"tag":"NNP","text":"Craig","dep":"nsubj","synid":-1,"semantic":"person"},{"index":1,"list":[],"tag":"VBZ","text":"has","dep":"ROOT","synid":-1,"semantic":"person"},{"index":2,"list":[3,1],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[1],"tag":"NN","text":"boat","dep":"dobj","synid":-1,"semantic":"vehicle"},{"index":4,"list":[1],"tag":"IN","text":"in","dep":"prep","synid":-1,"semantic":""},{"index":5,"list":[6,4,1],"tag":"DT","text":"the","dep":"det","synid":-1,"semantic":""},{"index":6,"list":[4,1],"tag":"NN","text":"harbour","dep":"pobj","synid":-1,"semantic":"location"},{"index":7,"list":[1],"tag":".","text":".","dep":"punct","synid":-1,"semantic":""}]}]`
	var sentence_list []Sentence
	err := json.Unmarshal([]byte(sentenceListStr), &sentence_list)
	util_ut.Check(t, err)
	isTrue(t, len(sentence_list)==1)
	ttree := SentenceToTuple(sentence_list[0])
	fmt.Printf("%s\n", ttree.ToString())
	isTrue(t, ttree.ToStringIndent() == "_has{VBZ}_ (Craig{person} | a boat{vehicle} in the harbour{location} .)")
}


// test sentence.IsQuestion()
func TestIsQuestion1(t *testing.T) {

	// sentence lists:

	// what will be indexed?
	const str1 = `[{"tokenList":[{"index":0,"list":[3],"tag":"WP","text":"What","dep":"nsubjpass","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":"man"},{"index":2,"list":[3],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"VBN","text":"indexed","dep":"ROOT","synid":-1,"semantic":""},{"index":4,"list":[3],"tag":".","text":"?","dep":"punct","synid":-1,"semantic":""}]}]`
	// what will be indexed
	const str2 = `[{"tokenList":[{"index":0,"list":[3],"tag":"WP","text":"What","dep":"nsubjpass","synid":-1,"semantic":""},{"index":1,"list":[3],"tag":"MD","text":"will","dep":"aux","synid":-1,"semantic":"man"},{"index":2,"list":[3],"tag":"VB","text":"be","dep":"auxpass","synid":-1,"semantic":""},{"index":3,"list":[],"tag":"VBN","text":"indexed","dep":"ROOT","synid":-1,"semantic":""}]}]`
	// who is Peter
	const str3 = `[{"tokenList":[{"index":0,"list":[1],"tag":"WP","text":"who","dep":"nsubj","synid":-1,"semantic":""},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"NNP","text":"Peter","dep":"attr","synid":-1,"semantic":"person"}]}]`
	// Peter is who
	const str4 = `[{"tokenList":[{"index":0,"list":[1],"tag":"NNP","text":"Peter","dep":"nsubj","synid":-1,"semantic":"person"},{"index":1,"list":[],"tag":"VBZ","text":"is","dep":"ROOT","synid":-1,"semantic":""},{"index":2,"list":[1],"tag":"WP","text":"who","dep":"dobj","synid":-1,"semantic":""}]}]`
	// do you have any clue
	const str5 = `[{"tokenList":[{"index":0,"list":[2],"tag":"VBP","text":"do","dep":"aux","synid":-1,"semantic":""},{"index":1,"list":[2],"tag":"PRP","text":"you","dep":"nsubj","synid":-1,"semantic":""},{"index":2,"list":[],"tag":"VB","text":"have","dep":"ROOT","synid":-1,"semantic":"person"},{"index":3,"list":[4,2],"tag":"DT","text":"any","dep":"det","synid":-1,"semantic":""},{"index":4,"list":[2],"tag":"NN","text":"clue","dep":"dobj","synid":-1,"semantic":""}]}]`

	isTrue(t, jsonToSentence(t, str1).IsQuestion())
	isTrue(t, jsonToSentence(t, str2).IsQuestion())
	isTrue(t, jsonToSentence(t, str3).IsQuestion())
	isTrue(t, jsonToSentence(t, str4).IsQuestion())
	isTrue(t, jsonToSentence(t, str5).IsQuestion())

	isTrue(t, !jsonToSentence(t, str1).IsImperative())
	isTrue(t, !jsonToSentence(t, str2).IsImperative())
	isTrue(t, !jsonToSentence(t, str3).IsImperative())
	isTrue(t, !jsonToSentence(t, str4).IsImperative())
}

// test sentence.IsImperative()
func TestIsImperative1(t *testing.T) {
	// make me a sandwich
	const str1 = `[{"tokenList":[{"index":0,"list":[],"tag":"VB","text":"make","dep":"ROOT","synid":-1,"semantic":""},{"index":1,"list":[3,0],"tag":"PRP","text":"me","dep":"nsubj","synid":-1,"semantic":""},{"index":2,"list":[3,0],"tag":"DT","text":"a","dep":"det","synid":-1,"semantic":""},{"index":3,"list":[0],"tag":"NN","text":"sandwich","dep":"ccomp","synid":-1,"semantic":""}]}]`

	isTrue(t, jsonToSentence(t, str1).IsImperative())

	isTrue(t, !jsonToSentence(t, str1).IsQuestion())
}

