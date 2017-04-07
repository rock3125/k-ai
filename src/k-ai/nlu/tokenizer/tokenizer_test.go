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

package tokenizer

import (
	"testing"
	"fmt"
)

// test simple string to tokens and back to string
func TestTokenizer1(t *testing.T) {
	// test a simple string is split into the right tokens
	fmt.Print("TestTokenizer1, Tokenize()\n")
	tokenList := Tokenize("This, is a  test string.  Is what I mean   to talk about?")
	if len(tokenList) != 26 {  // includes spaces and punctuation
		t.Error("len(tokenList) != 26")
	}

	// remove the spaces
	fmt.Print("TestTokenizer1, FilterOutSpaces()\n")
	tokenList = FilterOutSpaces(tokenList)
	if len(tokenList) != 15 {  // spaces removed
		t.Error("len(tokenList) != 15")
	}

	// test the "toString" method for pretty print
	fmt.Print("TestTokenizer1, ToString()\n")
	prettyStr := ToString(tokenList)
	if prettyStr != "This, is a test string. Is what I mean to talk about?" {
		t.Error("pretty print failed")
	}
}

// helper - test a contractions is as we'd expect it to be using the tokenizer
func checkContraction( t *testing.T, contractionStr string ) {
	fmt.Printf("checkContraction(%s)\n", contractionStr)
	tokenList := Tokenize(contractionStr)
	if len(tokenList) != 1 || tokenList[0].Text != contractionStr {
		if len(tokenList) == 1 {
			t.Errorf("contraction test failed: expected '%s' but got '%s", contractionStr, tokenList[0].Text)
		} else {
			t.Errorf("contraction test failed: size of token list(%d) != 1", len(tokenList))
		}
	}
}

// go through all the contractions
func TestContractions(t *testing.T) {
	checkContraction(t, "couldn't");
	checkContraction(t, "didn't");
	checkContraction(t, "don't");
	checkContraction(t, "doesn't");
	checkContraction(t, "he's");
	checkContraction(t, "how's");
	checkContraction(t, "I'd");
	checkContraction(t, "I'll");
	checkContraction(t, "I'm");
	checkContraction(t, "it'd");
	checkContraction(t, "isn't");
	checkContraction(t, "it'll");
	checkContraction(t, "it's");
	checkContraction(t, "might've");
	checkContraction(t, "mightn't");
	checkContraction(t, "must've");
	checkContraction(t, "mustn't");
	checkContraction(t, "she's");
	checkContraction(t, "she'll");
	checkContraction(t, "she's");
	checkContraction(t, "should've");
	checkContraction(t, "shouldn't");
	checkContraction(t, "we'd");
	checkContraction(t, "we'll");
	checkContraction(t, "we're");
	checkContraction(t, "we've");
	checkContraction(t, "weren't");
	checkContraction(t, "what're");
	checkContraction(t, "what've");
	checkContraction(t, "when's");
	checkContraction(t, "who'll");
	checkContraction(t, "who's");
	checkContraction(t, "won't");
	checkContraction(t, "would've");
	checkContraction(t, "wouldn't");
	checkContraction(t, "you'll");
	checkContraction(t, "you're");
}


// go through all the contractions
func TestContractionsCaseBased(t *testing.T) {
	checkContraction(t, "CouLdN'T")
	checkContraction(t,"Didn't")
}


// a sentence with a contraction
func TestContractionsSentence(t *testing.T) {
	fmt.Print("TestContractionsSentence, Tokenize()\n")
	tokenList := Tokenize("You shouldn't have.")
	if len(tokenList) != 6 {
		t.Error("TestContractionsSentence() len(tokenList) != 6")
	}
	if tokenList[2].Text != "shouldn't" {
		t.Error("TestContractionsSentence() \"shouldn't\" token incorrect @ position 2")
	}
}


// a sentence with a contraction
func TestSpecialCharacters(t *testing.T) {
	fmt.Print("TestSpecialCharacters, Tokenize()\n")
	tokenList := Tokenize("http://www.peter.co.nz")
	if len(tokenList) != 11 {
		t.Error("TestSpecialCharacters() len(tokenList) != 11")
	}
}




