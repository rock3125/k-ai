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
	"strings"
	"k-ai/nlu/model"
	"bytes"
)


/**
 * take a string apart into tokens
 * @param str the stirng to take apart
 * @return a list of tokens that makes the string
 */
func Tokenize(str string) []model.Token {
	tokenList := make([]model.Token,0)
	if ( len(str) > 0 ) {
		length := len(str)

		for i := 0; i < length; {
			tokenHandled := false

			// whitespace scanner
			ch := str[i:i+1]
			for isWhiteSpace(ch) && i < length {
				tokenHandled = true
				i = i + 1
				if ( i < length ) { ch = str[i:i+1] }
			}
			if ( tokenHandled ) {
				tokenList = append(tokenList, model.Token{Text: " "})
			}

			// add full-stops?
			for isFullStop(ch) && i < length {
				tokenHandled = true
				tokenList = append(tokenList, model.Token{Text: "."})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add hyphens?
			for isHyphen(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "-"})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add single quotes?
			for isSingleQuote(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "'"})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add double quotes?
			for isDoubleQuote(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "\""})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add special characters ( ) etc.
			for isSpecialCharacter(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: ch})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add punctuation ! ? etc.
			for isPunctuation(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: ch})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// numeric processor
			helper := ""
			for IsNumeric(ch) && i < length {
				tokenHandled = true;
				helper += ch;
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}
			if len(helper) > 0 {
				tokenList = append(tokenList, model.Token{Text: helper})
			}

			// text processor
			t_helper := ""
			for IsABC(ch) && i < length {
				tokenHandled = true;
				t_helper += ch;
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}
			if len(t_helper) > 0 {
				tokenList = append(tokenList, model.Token{Text: t_helper})
			}

			// discard unknown token?
			if ( !tokenHandled ) {
				i++; // skip
			}
		}
	}
	return handleContractions(tokenList)
}

/**
 * take a string apart into tokens with %1..n parameters
 * @param str the string to take apart
 * @return a list of tokens that makes the string
 */
func TokenizeWithParameter(str string) []model.Token {
	tokenList := make([]model.Token,0)
	if ( len(str) > 0 ) {
		length := len(str)

		for i := 0; i < length; {
			tokenHandled := false

			// whitespace scanner
			ch := str[i:i+1]
			for isWhiteSpace(ch) && i < length {
				tokenHandled = true
				i = i + 1
				if ( i < length ) { ch = str[i:i+1] }
			}
			if ( tokenHandled ) {
				tokenList = append(tokenList, model.Token{Text: " "})
			}

			// add full-stops?
			for isFullStop(ch) && i < length {
				tokenHandled = true
				tokenList = append(tokenList, model.Token{Text: "."})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add hyphens?
			for isHyphen(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "-"})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add single quotes?
			for isSingleQuote(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "'"})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add double quotes?
			for isDoubleQuote(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: "\""})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add special characters ( ) etc. but not %
			for ch != "%" && isSpecialCharacter(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: ch})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// add punctuation ! ? etc.
			for isPunctuation(ch) && i < length {
				tokenHandled = true;
				tokenList = append(tokenList, model.Token{Text: ch})
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}

			// numeric processor
			helper := ""
			for IsNumeric(ch) && i < length {
				tokenHandled = true;
				helper += ch;
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}
			if len(helper) > 0 {
				tokenList = append(tokenList, model.Token{Text: helper})
			}

			// text processor
			t_helper := ""
			for (IsABC(ch) || ch == "%") && i < length {
				tokenHandled = true;
				t_helper += ch;
				i = i + 1;
				if ( i < length ) { ch = str[i:i+1] }
			}
			if len(t_helper) > 0 {
				tokenList = append(tokenList, model.Token{Text: t_helper})
			}

			// discard unknown token?
			if ( !tokenHandled ) {
				i++; // skip
			}
		}
	}
	return handleContractions(tokenList)
}

/**
 * re-tokenize a spacey sentence because its tokenizer sucks - return null if there was no change
 * @param token a single token
 * @return a proper list of tokens
 */
func Retokenize(token model.Token) []model.Token {
	tokenList := Tokenize(token.Text);
	for i, _ := range tokenList {  // copy required fields
		tokenList[i].Tag = token.Tag
		tokenList[i].AncestorList = token.AncestorList
		tokenList[i].Dep = token.Dep
		tokenList[i].Index = token.Index
		tokenList[i].Semantic = token.Semantic
		tokenList[i].SynId = token.SynId
	}
	return tokenList
}

/**
 * get the item @ index in tokenList a contraction?
 * a contraction is always 3 symbols in the tokenizer
 * @param tokenList the list to check
 * @param index the index into the list
 * @return the contracted token or null if there is none there
 */
func getContraction(tokenList []model.Token, index int) *model.Token {
	if ( index + 2 < len(tokenList) ) {
		t1 := tokenList[index];
		t2 := tokenList[index+1];
		t3 := tokenList[index+2];
		if isContractionPrefix(strings.ToLower(t1.Text)) && t2.Text == "'" &&
			isContractionSuffix(strings.ToLower(t3.Text)) {
			contractionStr := t1.Text + "'" + t3.Text
			return &model.Token{Text: contractionStr}
		} else if t2.Text == "'" && strings.ToLower(t3.Text) == "s" {
			contractionStr := t1.Text + "'" + t3.Text
			return &model.Token{Text: contractionStr}
		}
	}
	return nil;
}

/**
 * fix contractions - just put them back into one word and the possessive
 * @param tokenList the list of tokens to fix
 * @return a fixed list of tokens
 */
func handleContractions( tokenList []model.Token ) []model.Token {
	tokenListWithContractions := make([]model.Token,0)
	for index := 0; index < len(tokenList); {
		contraction := getContraction(tokenList, index);
		if contraction != nil {
			tokenListWithContractions = append(tokenListWithContractions, *contraction)
			index = index + 3;
		} else {
			tokenListWithContractions = append(tokenListWithContractions, tokenList[index])
			index = index + 1;
		}
	}
	return tokenListWithContractions;
}

/**
 * given a list of tokens, remove all the white list items
 * @param tokenList a list of tokens in
 * @return the modified list of tokens with all white spaces removed
 */
func FilterOutSpaces( tokenList []model.Token ) []model.Token {
	resultList := make([]model.Token,0)
	for _, t_token := range tokenList {
		if !isWhiteSpace(t_token.Text) {
			resultList = append(resultList, t_token)
		}
	}
	return resultList;
}

/**
 * given a list of tokens, remove all the white list items
 * @param tokenList a list of tokens in
 * @return the modified list of tokens with all white spaces removed
 */
func FilterOutSpacesForSentences( sentenceList []model.Sentence ) []model.Sentence {
	resultList := make([]model.Sentence,0)
	for _, s_sentence := range sentenceList {
		if len(s_sentence.TokenList) > 0 {
			tokenList := make([]model.Token, 0)
			for _, t_token := range s_sentence.TokenList {
				if !isWhiteSpace(t_token.Text) {
					tokenList = append(tokenList, t_token)
				}
			}
			if len(tokenList) > 0 {
				resultList = append(resultList, model.Sentence{TokenList: tokenList})
			}
		}
	}
	return resultList;
}

/**
 * given a list of tokens, remove all the punctuation marks
 * @param tokenList a list of tokens in
 * @return the modified list of tokens with all white spaces removed
 */
func FilterOutPunctuation( tokenList []model.Token ) []model.Token {
	resultList := make([]model.Token,0)
	for _, t_token := range tokenList {
		if !isPunctuation(t_token.Text) && !isFullStop(t_token.Text) {
			resultList = append(resultList, t_token)
		}
	}
	return resultList;
}


func isWhiteSpace( ch string ) bool {
	return ch == " " || ch ==  "\t" || ch ==  "\r" || ch ==  "\n" || ch == "\u0008" ||
		   ch == "\ufeff" || ch == "\u303f" || ch == "\u3000" || ch == "\u2420" || ch == "\u2408" || ch == "\u202f" || ch == "\u205f" ||
		   ch == "\u2000" || ch == "\u2002" || ch == "\u2003" || ch == "\u2004" || ch == "\u2005" || ch == "\u2006" || ch == "\u2007" ||
		   ch == "\u2008" || ch == "\u2009" || ch == "\u200a" || ch == "\u200b";
}

func isFullStop( ch string ) bool {
	return ch == "\u002e" || ch == "\u06d4" || ch == "\u0701" || ch == "\u0702" ||
		   ch == "\ufe12" || ch == "\ufe52" || ch == "\uff0e" || ch == "\uff61";
}

func isHyphen( ch string ) bool {
	return ch == "\u002d" || ch == "\u207b" || ch == "\u208b" || ch == "\ufe63" || ch == "\uff0d" || ch =="\u2014";
}

func isSingleQuote( ch string ) bool {
	return ch == "'" || ch == "\u02bc" || ch == "\u055a" || ch == "\u07f4" || ch =="\u07f5" || ch == "\u2019" || ch =="\uff07" ||
		   ch == "\u2018" || ch == "\u201a" || ch == "\u201b" || ch == "\u275b" || ch == "\u275c";
}

// return true if ch is a double quote character
func isDoubleQuote( ch string ) bool {
	return ch == "\u0022" || ch == "\u00ab" || ch == "\u00bb" || ch == "\u07f4" || ch =="\u07f5" || ch == "\u2019" || ch == "\uff07" ||
		   ch == "\u201c" || ch == "\u201d" || ch == "\u201e" || ch == "\u201f" || ch =="\u2039" || ch == "\u203a" || ch == "\u275d" ||
		   ch == "\u276e" || ch == "\u2760" || ch == "\u276f";
}

func isPunctuation( ch string ) bool {
	return ch == "!" || ch == "?" || ch == "," || ch == ":" || ch == ";";
}

func IsNumeric( ch string ) bool {
	return ch >= "0" && ch <= "9";
}

func IsABC( ch string ) bool {
	return (ch >= "a" && ch <= "z") || (ch >= "A" && ch <= "Z");
}

func isSpecialCharacter( ch string ) bool {
	return strings.Index("_%$#@^&*()[]{}<>/\\=+|", ch) >= 0
}

func isContractionPrefix( str string ) bool {
	index := strings.Index("couldn didn doesn don hadn hasn haven he how i isn it " +
		"might mightn must mustn she we weren what when where who would wouldn you should shouldn won wont", str)
	return index >= 0
}

func isContractionSuffix( str string ) bool {
	index := strings.Index("ll d re s t ve m", str)
	return index >= 0
}

// word start symbols (preceded by a space but not followed)
func isWordStartSymbol( ch string ) bool {
	return ch == "\\" || ch == "/" || ch == "[" || ch == "{" || ch == "(";
}

// word end symbols (followed, but not preceded by a space)
func isWordEndSymbol( ch string ) bool {
	return ch == ":" || ch == "]" || ch == ")" || ch == "}" || ch == "." || ch == "!" || ch == "?" || ch == "," || ch == ";";
}

// words that shouldn't be preceded by a space
func noSpaceBefore( ch string ) bool {
	return ch == ";" || ch == "n't" || ch == "'s" || ch == "'ll" || ch == "." || ch == "," || ch == "?" ||
		ch == "!" || ch == ":" || ch == "'m" || ch == "'re" || ch == ")" || ch == "]" || ch == "}" || ch == "'ve";
}

// words that shouldn't be followed by a space
func noSpaceAfter( ch string ) bool {
	return ch == "(" || ch == "[" || ch == "{";
}

/**
 * format a sentence using rules for punctuation
 * @param tokenList a list of tokens to format with spaces between words
 * @return a pretty string resembling an ordinary readable sentence for humans
 */
func ToString( tokenList []model.Token ) string {
	list := make([]string,0)
	quote := 0 // quote counter
	for _, token := range tokenList {
		text := token.Text
		size := len(list)
		if noSpaceBefore(text) { // remove space before current item?
			if size > 0 && list[size - 1] == " " {
				list = list[:size-1]
			}
		}
		if text == "\"" { // count quotes
			quote += 1
		}
		if text == "\"" && quote % 2 == 0 { // end quotes
			if size > 0 && list[size - 1] == " " {
				list = list[:size-1]
			}
			list = append(list, text)
			list = append(list, " ")
		} else if text == "\"" { // start quote
			list = append(list, text)
		} else if noSpaceAfter(text) { // no spaces after this token
			list = append(list, text)
		} else { // all others
			list = append(list, text)
			list = append(list, " ")
		}
	}
	var buffer bytes.Buffer
	for _, str := range list {
		buffer.WriteString(str)
	}
	return strings.TrimSpace(buffer.String())
}

