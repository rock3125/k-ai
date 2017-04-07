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

package aiml

import (
	"k-ai/nlu/tokenizer"
	"strings"
	"k-ai/nlu/model"
)


/**
 * raw recursive string matcher - the first token must always match one of our constants
 * for this to be able to succeed - so no wildcards for the first character
 * @param str the string to match
 * @param manager the AIML manager library
 * @return null if failed, or a an AI/ML set of matching templates
 */
func (mgr *AimlManager) MatchTokenList(tokenList []model.Token) []model.AimlBinding {
	matchList := make([]model.AimlBinding,0)
	if len(tokenList) > 0 {
		tokenList := tokenizer.FilterOutPunctuation(tokenizer.FilterOutSpaces(tokenList))
		currentStr := strings.ToLower(tokenList[0].Text)
		if current, ok := mgr.NodeSet[currentStr]; ok {
			string_list := make([]model.AimlBinding, 0)

			// TODO: due to reloading of manager ability - we need to lock it for now to match
			mgr.Lock()
			defer mgr.Unlock()
			matchList = match(tokenList, 1, current, matchList, string_list)

		}
	}
	return matchList
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
func match(tokenList []model.Token, index int, rule *model.Aiml,
			matchList []model.AimlBinding,
			bindingList []model.AimlBinding) []model.AimlBinding {

	if len(tokenList) > 0 && index < len(tokenList) && rule != nil && len(rule.NodeSet) > 0 {
		// move on to the next token
		currentStr := strings.ToLower(tokenList[index].Text)
		if current, ok := rule.NodeSet[currentStr]; ok {
			matchList = match(tokenList, index + 1, current, matchList, bindingList)
		}

		// is there a wildcard available per chance?
		if current, ok := rule.NodeSet["*"]; ok {
			starList := make([]model.Token,0)
			first_current := current
			// start eating text till we get to a character that is allowed after the
			// wildcard, or we run out of characters
			for index < len(tokenList) {
				currentStr := strings.ToLower(tokenList[index].Text)
				if _, ok := current.NodeSet[currentStr]; ok {
					// found the next match - stop here
					current = current.NodeSet[currentStr]
					index = index + 1
					break;
				}
				starList = append(starList, tokenList[index])
				index = index + 1
			}

			// set the star value into the environment
			bindingList = append(bindingList, model.AimlBinding{Offset: len(matchList),
									Origin:                             first_current.Origin, TokenList: starList})
			if index < len(tokenList) {
				matchList = match(tokenList, index, current, matchList, bindingList)
			} else {
				if current != nil {
					return finish_bind(matchList, current.Origin, current.TemplateList, bindingList);
				}
			}
		}
		return matchList

	} else if len(tokenList) > 0 && index == len(tokenList) && rule != nil {
		// finally succeeded?
		if len(rule.TemplateList) == 0 { // end - but nothing to match
			// hope for a match on a "*"
			pattern, ok := rule.NodeSet["*"]
			if !ok {
				bindingList = append(bindingList, model.AimlBinding{Offset: len(matchList), Origin: rule.Origin}) // add empty pattern marker
				return finish_bind(matchList, rule.Origin, rule.TemplateList, bindingList);
			} else {
				bindingList = append(bindingList, model.AimlBinding{Offset: len(matchList), Origin: pattern.Origin}) // add empty pattern marker
				return finish_bind(matchList, pattern.Origin, pattern.TemplateList, bindingList);
			}
		} else {
			bindingList = append(bindingList, model.AimlBinding{Offset: len(matchList), Origin: rule.Origin}) // add empty pattern marker
			return finish_bind(matchList, rule.Origin, rule.TemplateList, bindingList);
		}
	}
	return matchList
}


/**
 * resolve / assign the bindings to the matches as required
 * @param matchList the list of matches store
 * @param list the items to add to the match store
 * @param bindings the bindings made along the way
 */
func finish_bind(matchList []model.AimlBinding, origin string, list []string, bindings []model.AimlBinding) []model.AimlBinding {
	final_match_list := make([]model.AimlBinding,0)
	for _, str := range list {
		matchList = append(matchList, model.AimlBinding{Text: str, Origin: origin})
	}
	if len(bindings) > 0 {
		newMatchList := make([]model.AimlBinding,0)
		stackIndex := bindings[0].Offset
		tokenList := bindings[0].TokenList
		nextStack := -1
		nextIndex := 1
		if len(bindings) > 1 {
			nextStack = bindings[nextIndex].Offset
		}
		for i := 0; i < len(matchList); i++ {
			match := matchList[i]
			// have we reached the next item in the binding stack?
			for nextStack >= 0 && i >= nextStack && nextIndex < len(bindings) {
				tokenList = bindings[nextIndex].TokenList
				stackIndex = bindings[nextIndex].Offset
				nextIndex += 1;
				if nextIndex < len(bindings) {
					nextStack = bindings[nextIndex].Offset
				} else {
					nextStack = -1; // there is no more next
				}
			}
			if i >= stackIndex && tokenList != nil && len(tokenList) > 0 {
				newMatch := model.AimlBinding{Text:     match.Text, Offset: match.Offset,
												Origin: origin, TokenList: match.TokenList}
				newMatch.TokenList = tokenList
				newMatchList = append(newMatchList, newMatch)
			} else {
				newMatchList = append(newMatchList, match)
			}
		}
		for _, item := range newMatchList {
			final_match_list = append(final_match_list, item)
		}
	}
	return final_match_list
}

