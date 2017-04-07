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
	"k-ai/nlu/model"
	"math"
	"strings"
	"bytes"
)

// hold return values of the getLargestMatching system
type matchingResult struct {
	newIndex int
	matchingToken model.Token
}

// maximum length of a word made up
const maxWordConstituentLength = 5

/**
 * Return the number of differences in case for two identical words
 * @param str1 word 1
 * @param str2 word 2
 * @return the number of differences in case
 */
func diffInCase( str1 string, str2 string ) int {
	if len(str1) == len(str2) {
		numDiff := 0
		for i := 0; i < len(str1); i++ {
			if str1[i] != str2[i] {
				numDiff++;
			}
		}
		return numDiff
	}
	return math.MaxInt32
}

/**
 * given a word, and items in a list from the lexicon, filter out those items
 * (if possible) that do not conform enough to the case wanted
 * @param list list of lexicon entries
 * @param word the word examined (from the original text)
 * @return an adjusted, or the original list
 */
func filterByCase(list []model.Token, word string, isStartOfSentence bool) []model.Token  {
	if len(word) > 0 && len(list) > 1 {
		returnList := make([]model.Token, 0)

		checkSize := 0 // one difference allowed for words @ start of sentenc
		if isStartOfSentence {
			checkSize = 1
		}
		if len(word) == 1 {
			checkSize = 0;
		}
		// are any of the words of the right case?
		rightCaseCount := 0
		for _, item := range list {
			if diffInCase(item.Text, word) <= checkSize {
				rightCaseCount++;
				returnList = append(returnList, item)
			}
		}
		// none of the words, or all of the words match
		if rightCaseCount == 0 || rightCaseCount == len(list) {
			return list
		} else {
			return returnList
		}
	}
	return list
}


// get the largest matching item from the lexicon
// remap: for remapping ancestors id_that_no_longer_exists -> correct_id
func (l *SLexicon) getLargestMatching(tokenList []model.Token, index int, remap map[int]int ) *matchingResult {
	size := maxWordConstituentLength;
	if index + size > len(tokenList) {
		size = len(tokenList) - index
	}
	resultList := make([]model.Token,0)
	resultListString := ""
	resultSize := 0

	sb := bytes.Buffer{}
	for i := 0; i < size; i++ {
		t_token := tokenList[index + i]
		wordStr := t_token.Text
		if len(wordStr) > 0 {
			sb.WriteString(wordStr)
			if tempList, ok := l.LWord[strings.ToLower(sb.String())]; ok {
				resultListString = sb.String()
				resultList = tempList
				resultSize = i + 1
			}
			sb.WriteString(" ")
		}
	}
	// return if we have a matching item
	if resultSize > 1 && len(resultList) > 0 {
		resultList = filterByCase(resultList, resultListString, index == 0);
		f_token := tokenList[index]
		// find the noun token if there is one in this set - avoid using adjectives as main markers
		offset := index
		for offset < (index + resultSize) && !strings.HasPrefix(f_token.Tag, "NN") {
			offset += 1
			f_token = tokenList[offset]
		}
		if !strings.HasPrefix(f_token.Tag, "NN") { // if we can't find it - just use the first token
			f_token = tokenList[index]
		}
		matchingToken := model.Token{
			Text: resultListString,
			Tag: f_token.Tag,
			AncestorList: f_token.AncestorList,
			Dep: f_token.Dep,
			Index: f_token.Index,
			Semantic: f_token.Semantic,
			SynId: f_token.SynId,
		}
		for i := index + 1; i < (index + resultSize); i++ {
			remap[i] = index
		}
		result := matchingResult{matchingToken: matchingToken, newIndex: index + resultSize}
		return &result
	}
	return nil
}


// find the longest sequence of words for a compound noun
func (l *SLexicon) GetLongestWordSequence(tokenList []model.Token) []model.Token {
	newTokenList := make([]model.Token,0)
	remap := make(map[int]int,0)
	for i := 0; i < len(tokenList); {
		wordStr := tokenList[i].Text
		if len(wordStr) > 0 {
			result := l.getLargestMatching(tokenList, i, remap)
			if result != nil {
				newTokenList = append(newTokenList, result.matchingToken)
				i = result.newIndex
			} else { // not in the lexicon
				newTokenList = append(newTokenList, tokenList[i])
				i++
			}
		} else {
			i++
		}
	}
	// remap the ancestor lists?
	if len(remap) > 0 {
		for _, t_token := range newTokenList {
			for j, value := range t_token.AncestorList {
				if offset, ok := remap[value]; ok {
					t_token.AncestorList[j] = offset
				}
			}
		}
	}
	return newTokenList
}


// find the longest sequence of words for a compound noun in a set of sentences
func (l *SLexicon) GetLongestWordSequenceForList(sentenceList []model.Sentence) []model.Sentence {
	newSentenceList := make([]model.Sentence,0)
	for _, sentence := range sentenceList {
		new_sentence := model.Sentence{TokenList: l.GetLongestWordSequence(sentence.TokenList)}
		newSentenceList = append(newSentenceList, new_sentence)
	}
	return newSentenceList
}

