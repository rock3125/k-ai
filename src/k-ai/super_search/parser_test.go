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
	"runtime/debug"
	"k-ai/util_ut"
)

// test fn
func isTrue(t *testing.T, cond bool) {
	if !cond {
		t.Error("condition failed")
		debug.PrintStack()
		t.FailNow()
	}
}

// test the tuple tree parsing content
func TestWord1(t *testing.T) {
	item, err := parse_string("any(hello there!)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "word")
	isTrue(t, item.Word == "hello there!")
	isTrue(t, len(item.Semantic) == 0)
	isTrue(t, len(item.Tag) == 0)
	isTrue(t, !item.Exact)
}


func TestWord2(t *testing.T) {
	item, err := parse_string("any(Peter,nnp)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "word")
	isTrue(t, item.Word == "Peter")
	isTrue(t, item.Tag == "NNP")
	isTrue(t, len(item.Semantic) == 0)
	isTrue(t, !item.Exact)
}

func TestWord3(t *testing.T) {
	item, err := parse_string("person(Peter)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "word")
	isTrue(t, item.Word == "Peter")
	isTrue(t, item.Tag == "NNP")
	isTrue(t, item.Semantic == "person")
	isTrue(t, !item.Exact)
}

func TestWord4(t *testing.T) {
	item, err := parse_string("exact person(Peter)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "word")
	isTrue(t, item.Word == "Peter")
	isTrue(t, item.Tag == "NNP")
	isTrue(t, item.Semantic == "person")
	isTrue(t, item.Exact)
}

func TestAnd1(t *testing.T) {
	item, err := parse_string("any(test) and any(Peter dearest)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "and")
	isTrue(t, item.Right != nil)
	isTrue(t, item.Left != nil)

	any_test := item.Left
	isTrue(t, any_test.TType == "word")
	isTrue(t, any_test.Word == "test")
	isTrue(t, !any_test.Exact)

	peter_dearest := item.Right
	isTrue(t, peter_dearest.TType == "word")
	isTrue(t, peter_dearest.Word == "Peter dearest")
	isTrue(t, !peter_dearest.Exact)
}


func TestAnd2(t *testing.T) {
	item, err := parse_string("  any(test)  and  person(Peter best)  and exact any(Markie Mark)")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "and")
	isTrue(t, item.Right != nil)
	isTrue(t, item.Left != nil)

	next_and := item.Right
	isTrue(t, next_and.TType == "and")

	left1 := item.Left
	isTrue(t, left1.TType == "word")
	isTrue(t, left1.Word == "test")
	isTrue(t, !left1.Exact)

	left2 := next_and.Left
	isTrue(t, left2.TType == "word")
	isTrue(t, left2.Word == "Peter best")
	isTrue(t, !left2.Exact)
	isTrue(t, left2.Semantic == "person")
	isTrue(t, left2.Tag == "NNP")

	right2 := next_and.Right
	isTrue(t, right2.TType == "word")
	isTrue(t, right2.Word == "Markie Mark")
	isTrue(t, right2.Exact)
}

func TestBrackets1(t *testing.T) {
	item, err := parse_string("(any(hello there!))")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "word")
	isTrue(t, item.Word == "hello there!")
	isTrue(t, !item.Exact)
}

func TestAndNot1(t *testing.T) {
	item, err := parse_string("  location(test)  and  not   exact  person (Peter)  ")
	util_ut.Check(t, err)
	isTrue(t, item != nil)
	isTrue(t, item.TType == "and not")
	isTrue(t, item.Right != nil)
	isTrue(t, item.Left != nil)

	left1 := item.Left
	isTrue(t, left1.TType == "word")
	isTrue(t, left1.Word == "test")
	isTrue(t, !left1.Exact)
	isTrue(t, left1.Semantic == "location")

	right1 := item.Right
	isTrue(t, right1.TType == "word")
	isTrue(t, right1.Word == "Peter")
	isTrue(t, right1.Exact)
	isTrue(t, right1.Semantic == "person")
}


