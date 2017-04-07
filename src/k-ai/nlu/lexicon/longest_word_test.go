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
	"k-ai/util_ut"
	"encoding/json"
	"k-ai/nlu/model"
)

// from json str back to a single sentence (the first)
func jsonToTokenList(t *testing.T, str string) []model.Token {
	var token_list []model.Token
	err := json.Unmarshal([]byte(str), &token_list)
	util_ut.Check(t, err)
	return token_list
}

// test "the old lady" gets its proper longest word assignment
func TestLongestWord_1(t *testing.T) {
	// "the old lady" -> "the [old lady]" one concept, noun, with semantic female and ancestors adjusted
	const text = `[{"index":0,"list":[],"tag":"DT","text":"The","dep":"ROOT","synid":-1,"semantic":""},
					{"index":1,"list":[1],"tag":"JJ","text":"old","dep":"adj","synid":-1,"semantic":""},
					{"index":2,"list":[2],"tag":"NN","text":"lady","dep":"nsubj","synid":-1,"semantic":""}]`

	token_list := jsonToTokenList(t, text)
	util_ut.IsTrue(t, len(token_list) == 3)

	new_token_list := Lexi.GetLongestWordSequence(token_list)
	util_ut.IsTrue(t, len(new_token_list) == 2)

	old_lady := new_token_list[1]
	util_ut.IsTrue(t, old_lady.Tag == "NN")
	util_ut.IsTrue(t, old_lady.Text == "old lady")
	util_ut.IsTrue(t, old_lady.Dep == "nsubj")
	util_ut.IsTrue(t, len(old_lady.AncestorList) == 1 && old_lady.AncestorList[0] == 1)
}

