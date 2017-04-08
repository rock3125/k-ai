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

package freebase

import (
	"k-ai/nlu/tokenizer"
	"strings"
	"k-ai/nlu/model"
)


/**
 * raw recursive string/token matcher
 * @param tokenList - list of tokens to match
 * @param nodeSet - the node set tree
 * @return null if failed, or an executable item that matches the pattern
 */
func MatchTokenList(tokenList []model.Token, nodeSet map[string]*FBPatternTree, vars map[string][]model.Token) []model.Token {
	if len(tokenList) > 0 {
		tokenList := tokenizer.FilterOutPunctuation(tokenizer.FilterOutSpaces(tokenList))
		currentStr := strings.ToLower(tokenList[0].Text)
		if current, ok := nodeSet[currentStr]; ok {
			return match(tokenList, 1, current, vars)
		}
	}
	return make([]model.Token,0)
}

// get % parameter
func get_parameter(node *FBPatternTree) (string, *FBPatternTree) {
	// is there a wildcard available per chance?
	pattern := ""
	var pattern_node *FBPatternTree
	for key, value := range node.Patterns {
		if strings.Contains(key, "%") {
			pattern = key
			pattern_node = value
		}
	}
	return pattern, pattern_node
}

/**
 * recursive matcher helper
 *
 * must become more sophisticated - on a no match - recurse to potential star match patterns
 * and keep multiple possible matches
 *
 * @param tokenList the liste of tokens to match
 * @param index the index into the list of tokens to match
 * @param rule the current rule to use for matching
 */
func match(tokenList []model.Token, index int, node *FBPatternTree,
			vars map[string][]model.Token) []model.Token {

	if len(tokenList) > 0 && index < len(tokenList) && node != nil {
		// move on to the next token
		currentStr := strings.ToLower(tokenList[index].Text)
		if current, ok := node.Patterns[currentStr]; ok {
			executable := match(tokenList, index + 1, current, vars)
			if len(executable) > 0 {
				return executable
			}
		}

		// is there a wildcard available per chance?
		pattern, pattern_node := get_parameter(node)
		if len(pattern) > 0 {
			starList := make([]model.Token,0)
			first_current := currentStr
			// start eating text till we get to a character that is allowed after the
			// wildcard, or we run out of characters
			for index < len(tokenList) {
				currentStr := strings.ToLower(tokenList[index].Text)
				if _, ok := pattern_node.Patterns[currentStr]; ok {
					// found the next match - stop here
					pattern_node = pattern_node.Patterns[currentStr]
					index = index + 1
					break;
				}
				starList = append(starList, tokenList[index])
				index = index + 1
			}

			// set the star value into the environment
			vars[first_current] = starList
			if index < len(tokenList) {
				return match(tokenList, index, pattern_node, vars)
			} else {
				if pattern_node != nil {
					return finish_bind(pattern_node, vars);
				}
			}
		}
		return make([]model.Token, 0)

	} else if len(tokenList) > 0 && index == len(tokenList) && node != nil {
		// finally succeeded?
		if len(node.Patterns) == 0 { // end - but nothing to match
			// hope for a match on a "*"
			pattern, pattern_node := get_parameter(node)
			if len(pattern) == 0 {
				return make([]model.Token, 0)
			} else {
				return finish_bind(pattern_node, vars);
			}
		} else {
			return finish_bind(node, vars);
		}
	}
	return make([]model.Token, 0)
}


/**
 * resolve / assign the bindings to the matches as required
 * @param matchList the list of matches store
 * @param list the items to add to the match store
 * @param bindings the bindings made along the way
 */
func finish_bind(node *FBPatternTree, vars map[string][]model.Token) []model.Token {
	return node.Executable
}

