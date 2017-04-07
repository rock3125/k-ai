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

package grammar

import (
	"testing"
	"k-ai/nlu/tokenizer"
	"runtime/debug"
)

func checkSingle(t *testing.T, glib *GrammaryLibrary, str string) {
	tokenList := glib.Parse(tokenizer.Tokenize(str))
	if len(tokenList) != 1 {
		t.Errorf("len(tokenList) != 1, but %d for %s", len(tokenList), str)
		debug.PrintStack()
		t.FailNow()
	}
	if tokenList[0].Text != str {
		t.Errorf("token-format incorrect (%s) != (%s)", tokenList[0].Text, str)
		debug.PrintStack()
		t.FailNow()
	}
}

//  your email address is test@test.com
func TestSpacesSurviveAndDontBecomeOneWord(t *testing.T) {

	parse_str := "My email address is peter@peter.co.nz"
	tokenList := Grammar.Parse(tokenizer.Tokenize(parse_str))
	if len(tokenList) == 1 {
		t.Error("error len(tokenList) == 1")
	}
}

// test simple string to tokens and back to string
func TestSimpleSinglePatterns(t *testing.T) {

	checkSingle(t, &Grammar, "1 January 2016")
	checkSingle(t, &Grammar, "25 Mar 2099")

	checkSingle(t, &Grammar, "2016-04-18 15:59:07")
	checkSingle(t, &Grammar, "2016-04-18 15:59:07.123")
	checkSingle(t, &Grammar, "2011-01-31")
	checkSingle(t, &Grammar, "2016-04-12 11:22 PM")
	checkSingle(t, &Grammar, "June 1, 2001")

	checkSingle(t, &Grammar, "mailto://Blair-l/customer___oneok/22.txt")
	checkSingle(t, &Grammar, "http://www.peter.co.nz")
	checkSingle(t, &Grammar, "www.peter.co.nz")

	checkSingle(t, &Grammar, "11:23:00 AM")
	checkSingle(t, &Grammar, "23:23")
	checkSingle(t, &Grammar, "23:23:00.000")
	checkSingle(t, &Grammar, "00:00:00 PM")
	checkSingle(t, &Grammar, "00:00:00.000")
	checkSingle(t, &Grammar, "23:59:59.999")

	checkSingle(t, &Grammar, "713-853-5660")
	checkSingle(t, &Grammar, "(713) 853-5660")
}

// test simple string to tokens and back to string
func TestSentencePatterns(t *testing.T) {
	str := "mailto://Blair-l/customer___oneok/22.txt,mailto://Blair-l/customer___oneok/23.txt"
	tokenList := Grammar.Parse(tokenizer.Tokenize(str))
	if len(tokenList) != 3 {
		t.Errorf("len(tokenList) != 3, but %d for %s", len(tokenList), str)
	}
	if tokenList[0].Text != "mailto://Blair-l/customer___oneok/22.txt" || tokenList[1].Text != "," ||
		tokenList[2].Text != "mailto://Blair-l/customer___oneok/23.txt" {
		t.Errorf("format incorrect (%s) != (%s)", tokenizer.ToString(tokenList), str)
	}
}

