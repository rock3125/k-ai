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

package lexicon

import (
	"testing"
	"fmt"
	"k-ai/util_ut"
	"runtime/debug"
)

// test fn
func contains(t *testing.T, word_list []string, str...string) {
	for _, item := range str {
		found := false
		for _, word := range word_list {
			if word == item {
				found = true
			}
		}
		if !found {
			err_str := fmt.Sprintf("word %s missing from list", item)
			t.Error(err_str)
			debug.PrintStack()
			t.FailNow()
		}
	}
}

// test a simple rule with name substitution (from ai.aiml, line 82)
func TestStemming1(t *testing.T) {
	util_ut.IsTrue(t, Lexi.GetStem("gasses") == "gas")
	util_ut.IsTrue(t, Lexi.GetStem("swimming") == "swim")
	util_ut.IsTrue(t, Lexi.GetStem("wrote") == "write")
}

// test the undesirables work
func TestUndesirables(t *testing.T) {
	util_ut.IsTrue(t, Lexi.IsUndesirable("the"))
	util_ut.IsTrue(t, Lexi.IsUndesirable("a"))
	util_ut.IsTrue(t, Lexi.IsUndesirable("an"))
	util_ut.IsTrue(t, Lexi.IsUndesirable("in"))
	util_ut.IsTrue(t, Lexi.IsUndesirable("my"))
}

// test related words
func TestStemWords(t *testing.T) {
	word_list_1 := Lexi.GetStemList("swim")
	contains(t, word_list_1, "swam", "swimming", "swims", "swum")

	word_list_2 := Lexi.GetStemList("gas")
	contains(t, word_list_2, "gasses", "gassed", "gassing")

	word_list_3 := Lexi.GetStemList("car")
	contains(t, word_list_3, "cars", "car's")
}

// test the auxilliaries are present and stemming
func TestStemAux(t *testing.T) {
	word_list_1 := Lexi.GetStemList("be")
	contains(t, word_list_1, "was", "is", "am", "are", "being", "been", "were")

	word_list_2 := Lexi.GetStemList("do")
	contains(t, word_list_2, "did", "does", "done", "doing")

	word_list_3 := Lexi.GetStemList("have")
	contains(t, word_list_3, "had", "has", "having")
}

// test synonyms
func TestSynonyms(t *testing.T) {
	word_list_1 := Lexi.GetSynonymList("car")
	contains(t, word_list_1, "automobile", "auto")
}

// test semantics
func TestSemantics(t *testing.T) {
	util_ut.IsTrue(t, Lexi.GetSemantic("John") == "male")
	util_ut.IsTrue(t, Lexi.GetSemantic("john") == "location")
}

