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

package parser

import (
	"testing"
	"k-ai/nlu/lexicon"
	"fmt"
	"k-ai/util_ut"
)

// test we can parse
func TestParser1(t *testing.T) {
	parse_str := "Peter de Vocht was here.  He then moved to London."
	sentence_list, err := ParseText(parse_str)
	util_ut.IsTrue(t, err == nil)
	util_ut.IsTrue(t, len(sentence_list) == 2)

	sentence_1 := sentence_list[0]
	util_ut.IsTrue(t, len(sentence_1.TokenList) == 4)

	sentence_2 := sentence_list[1]
	util_ut.IsTrue(t, len(sentence_2.TokenList) == 6)

	//bank_text, err := util.LoadTextFile(util.GetDataPath() + "/test_data/kb/government.txt")
	//if err == nil {
	//	str, _ := ParseTextToJson(bank_text)
	//	fmt.Printf("%s\n", str)
	//}

	str, _ := ParseTextToJson("John said he likes dogs.")
	fmt.Printf("%s\n", str)
}

// test we can parse a compound name after it has been added to the lexicon
func TestParser2(t *testing.T) {
	lexicon.Lexi.AddSemantic("Davis Patrick", "person")

	parse_str := "Peter and Davis Patrick went to the beach at 12:45 to view the boats comming in."
	ttList, err := ParseTextToTupleTree(parse_str)
	util_ut.IsTrue(t, err == nil)
	util_ut.IsTrue(t, len(ttList) == 1)

	// find davis patrick in its rightful place in the tree
	t_left_1 := ttList[0].Left
	util_ut.IsTrue(t, t_left_1 != nil)
	t_rght_2 := t_left_1.Right
	util_ut.IsTrue(t, t_rght_2 != nil)
	t_rght_3 := t_rght_2.Right
	util_ut.IsTrue(t, t_rght_3 != nil && len(t_rght_3.Tokens.TokenList) == 1)
	dp_token := t_rght_3.Tokens.TokenList[0]
	util_ut.IsTrue(t, dp_token.Text == "Davis Patrick" && dp_token.Semantic == "person")
}

// test parsing: "I am KAI, an Artificial Intelligence."
func TestParser3(t *testing.T) {
	parse_str := "I am KAI, an Artificial Intelligence."
	ttList, err := ParseTextToTupleTree(parse_str)
	util_ut.IsTrue(t, err == nil)
	util_ut.IsTrue(t, len(ttList) == 1)
	str := ttList[0].ToStringIndent()
	util_ut.IsTrue(t,  str == "_am{VBP}_ (I{nsubj} | KAI{ai} , an Artificial Intelligence .)")
}
