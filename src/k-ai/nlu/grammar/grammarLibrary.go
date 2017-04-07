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
	"strings"
	"strconv"
	"k-ai/util"
	"k-ai/nlu/tokenizer"
	"k-ai/nlu/model"
	"k-ai/logger"
)

type Match struct {
	ruleName string
    resultList []model.Token
	index int
}

type GrammaryLibrary struct {
	initialised bool							// setup?
	GrammarMap map[string]*GrammarLhs			// the map for looking up items
	GrammarConversionMap map[string]string		// the map for looking up item's conversion patterns
	GrammarModificationMap map[string]string	// the map for looking up item's modification patterns
	StartPattern map[string][]GrammarLhs		// a start token map for matching possible rules
}

// setup maps
func (g *GrammaryLibrary) Init() {
	g.GrammarMap = make(map[string]*GrammarLhs,0)
	g.GrammarConversionMap = make(map[string]string,0)
	g.GrammarModificationMap = make(map[string]string,0)
	g.StartPattern = make(map[string][]GrammarLhs,0)
	g.initialised = true
}

// return true if ch is a number 0..9
func isNumeric(ch string) bool {
	return len(ch) > 0 && ch[0:1] >= "0" && ch[0:1] <= "9"
}

// return true if ch is a number a..zA..Z
func isABC(ch string) bool {
	return len(ch) > 0 && ((ch[0:1] >= "a" && ch[0:1] <= "z") || (ch[0:1] >= "A" && ch[0:1] <= "Z"))
}

// parse the .range() part of a potential query
func (g GrammaryLibrary) parseRange( str string, node *GrammarRhs ) string {
	if ( node != nil && strings.Contains(str, ".range(") ) {
		index := strings.Index(str, ".range(")
		returnStr := str[:index]
		rangeStr := str[index + 7:]
		index2 := strings.Index(rangeStr, ")")
		if ( index2 < 3 ) {
			panic(".range() pattern missing )")
		}
		rangeStr = rangeStr[:index2]
		parts := strings.Split(rangeStr, ",")
		if len(parts) != 2 {
			panic(".range() must have two comma separated items")
		}
		node.NumberRangeStart, _ = strconv.Atoi(parts[0])
		node.NumberRangeEnd, _ = strconv.Atoi(parts[1])
		return returnStr
	}
	return str
}


// process a single rhs rule
func (g GrammaryLibrary) parseGrammarRhs( rhs string ) []GrammarRhs {
	if len(rhs) == 0 {
		panic("Grammar rhs empty")
	}
	resultList := make([]GrammarRhs, 0)
	// or bag of words rule?
	if strings.HasPrefix(rhs,"[") || strings.HasSuffix(rhs, "]") {
		rhs := strings.TrimSpace(rhs[1:len(rhs)-1])
		bag := strings.Split(rhs, " ")
		node := GrammarRhs{}
		node.Init()
		for _, value := range bag {
			node.PatternSet[value] = true
		}
		resultList = append(resultList, node)
	} else {
		// ordinary ordered list
		stringList := strings.Split(rhs, " ")
		for _, str := range stringList {
			isRepeat := false
			if len(str) > 1 && strings.HasSuffix(str, "+") {
				isRepeat = true
				str = str[:len(str)-1]
			}
			node := GrammarRhs{}
			node.Init()
			node.Text = g.parseRange(str, &node)
			node.IsRepeat = isRepeat
			resultList = append(resultList, node)
		}
	}
	return resultList
}


// process a single line
func (g GrammaryLibrary) processPattern( line string ) *GrammarLhs {
	if len(line) > 0 && strings.Contains(line, "=") {
		index := strings.Index(line, "=")
		if ( index > 0 ) {
			lhs := strings.TrimSpace(line[:index])
			rhs := strings.TrimSpace(line[index+1:])

			lhsParts := strings.Split(lhs, " ")

			if len(lhsParts) != 2 {
				panic("Grammar pattern must have private/public/pattern name part")
			}
			if lhsParts[0] != "private" && lhsParts[0] != "public" && lhsParts[0] != "pattern" && lhsParts[0] != "modifier" {
				panic("Grammar pattern must start with 'public', 'private' or 'pattern'")
			}

			// special conversion patterns for rules
			if lhsParts[0] == "pattern" {
				lhs := GrammarLhs{ Name: strings.TrimSpace(lhsParts[1]), ConversionPattern: strings.TrimSpace(rhs) }
				lhs.Init()
				return &lhs

			} else if lhsParts[0] == "modifier" {
				lhs := GrammarLhs{ Name: strings.TrimSpace(lhsParts[1]), Modifier: strings.TrimSpace(rhs) }
				lhs.Init()
				return &lhs

			} else {

				grammarLhs := GrammarLhs{ IsPublic: lhsParts[0] == "public", Name: strings.TrimSpace(lhsParts[1]) }
				grammarLhs.RhsList = g.parseGrammarRhs(strings.TrimSpace(rhs))
				return &grammarLhs
			}
		}
	}
	return nil
}


// setup the lhs start non terminal lookup(s)
func (g GrammaryLibrary) setupFirstLetterLookup() {
	// resolve any reference
	for _, lhs := range g.GrammarMap {
		if lhs.IsPublic {
			tokenList := lhs.GetStartTokens()
			if len(tokenList) == 0 {
				panic("invalid return result null")
			}
			for _, str := range tokenList {
				g.StartPattern[str] = append(g.StartPattern[str], *lhs)
			}
		}
	}
}


// setup the lhs start non terminal lookup(s)
func (g GrammaryLibrary) getRulesByFirstLetter(firstLetter string) []GrammarLhs {
	return g.StartPattern[firstLetter]
}


// resolve references to other rules where possible
func (g GrammaryLibrary) resolveReferences() {
	// resolve any reference
	for _, lhs := range g.GrammarMap {
		rhsList := lhs.RhsList
		if rhsList == nil {
			panic("invalid Grammar rule, no rhs '" + lhs.Name + "'")
		}
		for index, rhs := range rhsList {
			if len(rhs.Text) > 0 {
				// text, but not abc or number marker
				if ! ( rhs.Text == "abc" || rhs.Text == "number" || rhs.Text == "space") && isABC(rhs.Text[0:1]) {
					if rhs.Text == lhs.Name {
						panic("rule '" + lhs.Name + "' cyclic reference")
					}
					// can this be resolved?
					if reference, ok := g.GrammarMap[rhs.Text]; ok {
						rhsList[index].Text = ""
						rhsList[index].Reference = reference
					}
				}
			}
		}
	}
}


// extend a list of tokens with another list of tokens
func list_append(list1 []model.Token, list2 []model.Token) []model.Token {
	for _, item := range list2 {
		list1 = append(list1, item)
	}
	return list1
}


// match the longest possible chain of rules from a list of rules
func match_rhs( tokenList []model.Token, index int, rhs *GrammarRhs ) *Match {
	if rhs != nil && len(tokenList) > 0 && index < len(tokenList) {
		r_match := Match{}
		if rhs.Reference != nil { // other rule reference
			m := match_lhs(tokenList, index + r_match.index, rhs.Reference)
			if m != nil { // update recursive state
				r_match.index = r_match.index + m.index
				r_match.resultList = list_append(r_match.resultList, m.resultList)

				// repeat of reference?
				if rhs.IsRepeat {
					m = match_lhs(tokenList, index + r_match.index, rhs.Reference)
					for m != nil {
						r_match.index = r_match.index + m.index
						r_match.resultList = list_append(r_match.resultList, m.resultList)
						m = match_lhs(tokenList, index + r_match.index, rhs.Reference)
					}
				}

			} else { // fail
				r_match.index = -1
			}

		} else if len(rhs.Text) > 0 { // literal
			i_token := tokenList[index]
			if rhs.Text == "abc" && isABC(i_token.Text) { // text is valid
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else if rhs.Text == "number" && isNumeric(i_token.Text) { // text is valid
				// range check?
				if rhs.NumberRangeStart != rhs.NumberRangeEnd || rhs.NumberRangeStart != 0 {
					// can only really check numbers in the 64 bit range
					if len(i_token.Text) <= 12 {
						value, _ := strconv.Atoi(i_token.Text)
						if rhs.NumberRangeStart <= value && value <= rhs.NumberRangeEnd { // within range?
							r_match.resultList = append(r_match.resultList, i_token)
							r_match.index++
						} else {
							r_match.index = -1 // outside range, fail
						}
					} else {
						r_match.index = -1 // too big - fail
					}
				} else {
					r_match.resultList = append(r_match.resultList, i_token)
					r_match.index++
				}
			} else if rhs.Text == "space" && i_token.Text == " " {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else if rhs.Text == i_token.Text {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else { // fail
				r_match.index = -1
			}

		} else if len(rhs.PatternSet) > 0 { // literal
			i_token := tokenList[index]
			if _, ok := rhs.PatternSet["abc"]; ok && isABC(i_token.Text) {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else if _, ok := rhs.PatternSet["number"]; ok && isNumeric(i_token.Text) {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else if _, ok := rhs.PatternSet["space"]; ok && i_token.Text == " " {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else if _, ok := rhs.PatternSet[i_token.Text]; ok {
				r_match.resultList = append(r_match.resultList, i_token)
				r_match.index++
			} else { // fail
				r_match.index = -1
			}
		}

		// did we get a valid match
		if ( r_match.index > 0 ) {
			if rhs.IsRepeat {  // try again / recurse?
				temp := match_rhs( tokenList, index + r_match.index, rhs )
				if temp != nil {
					r_match.index += temp.index
					r_match.resultList = list_append(r_match.resultList, temp.resultList)
				}
			}
			return &r_match
		}
	}
	return nil
}


/**
 * apply the specified modification to the token list, for now the only modifcation supported
 * is space@index, insert a space @ index
 * @param modification the requested modification
 * @param list the list to modify
 * @return the modified list
 */
func (g *GrammaryLibrary) modifySet(modification string, list []model.Token) []model.Token {
	if len(modification) >0 && !strings.HasPrefix(modification, "space@") {
		panic("bad modification string:" + modification)
	}
	if len(list) > 0 {
		index, _:= strconv.Atoi(strings.Split(modification, "@")[1])
		list[index] = model.Token{Text: " "}
	}
	return list
}

// match the longest possible chain of rules from a list of rules
func match_lhs_list( tokenList []model.Token, index int, ruleSet []GrammarLhs) *Match {
	if len(tokenList) > 0 && index < len(tokenList) && len(ruleSet) > 0 {
		r_match := Match{index: -1}
		for _, lhs := range ruleSet {
			temp := match_lhs( tokenList, index, &lhs )
			if temp != nil {
				// careful - this doesn't allow for two matching rules
				// if two rules match exactly - the first one will be chosen
				if r_match.index == -1 || r_match.index < temp.index {
					r_match.index = temp.index
					r_match.ruleName = temp.ruleName
					r_match.resultList = temp.resultList
				}
			}
		}
		if r_match.index > 0 {
			r_match.index = r_match.index + index // offset into token set
			return &r_match
		}
	}
	return nil
}

// match the longest possible chain of rules from a list of rules
func match_lhs( tokenList []model.Token, index int, rule *GrammarLhs ) *Match {
	if len(tokenList) > 0 && index < len(tokenList) && rule != nil && len(rule.RhsList) > 0 {
		match := Match{index: 0}
		for _, rhs := range rule.RhsList {
			if index + match.index >= len(tokenList) { // failed - no more tokens - but rule not finished
				match.index = -1
				break
			}
			temp := match_rhs( tokenList, index + match.index, &rhs )
			if temp != nil {
				match.index += temp.index
				match.resultList = list_append(match.resultList, temp.resultList)
			} else {
				match.index = -1
				break
			}
		}
		if ( match.index > 0 ) {
			match.ruleName = rule.Name
			return &match
		}
	}
	return nil
}


// find any Grammar rules that match and apply them - return
// new tokens based on the Grammar rules that applied and that didn't
func (g GrammaryLibrary) ParseSentenceList(sentenceList []model.Sentence) []model.Sentence {
	if g.initialised {
		newSentenceList := make([]model.Sentence, 0)
		for _, sentence := range sentenceList {
			newTokenList := make([]model.Token,0)
			for _, t_token := range sentence.TokenList {
				newTokenList = append(newTokenList, t_token)
				newTokenList = append(newTokenList, model.Token{Tag: " ", Text: " "})
			}
			newSentenceList = append(newSentenceList, model.Sentence{TokenList: g.Parse(newTokenList)})
		}
		return newSentenceList
	} else {
		return sentenceList
	}
}

// find any Grammar rules that match and apply them - return
// new tokens based on the Grammar rules that applied and that didn't
func (g GrammaryLibrary) Parse(tokenList []model.Token) []model.Token {
	if g.initialised {
		newTokenList := make([]model.Token, 0)

		// split spacy tokens back into basic components where possible
		correctedTokenList := make([]model.Token, 0)
		for _, t_token := range tokenList {
			manyTokenList := tokenizer.Retokenize(t_token)
			correctedTokenList = list_append(correctedTokenList, manyTokenList)
		}
		tokenList = correctedTokenList

		for i := 0; i < len(tokenList); {
			t_token := tokenList[i]

			// literal first - more specific
			result := &Match{index: -1}
			ruleSet := g.getRulesByFirstLetter(t_token.Text)
			if len(ruleSet) > 0 {
				result = match_lhs_list(tokenList, i, ruleSet)
			}
			if result == nil || result.index == -1 {
				if tokenizer.IsABC(t_token.Text) {
					ruleSet = g.getRulesByFirstLetter("abc")
					result = match_lhs_list(tokenList, i, ruleSet)
				} else if tokenizer.IsNumeric(t_token.Text) {
					ruleSet = g.getRulesByFirstLetter("number")
					result = match_lhs_list(tokenList, i, ruleSet)
				} else if t_token.Text == " " {
					ruleSet = g.getRulesByFirstLetter(" ")
					result = match_lhs_list(tokenList, i, ruleSet)
				}
			}
			if result != nil && len(result.ruleName) > 0 {
				if rule, ok := g.GrammarModificationMap[result.ruleName]; ok {
					result.resultList = g.modifySet(rule, result.resultList)
				}
				resultStr := ""
				for _, t := range result.resultList {
					resultStr += t.Text
				}
				newToken := model.Token{
					Text: resultStr,
					Dep: t_token.Dep,
					Index: t_token.Index,
					AncestorList: t_token.AncestorList,
					Semantic: result.ruleName,
					Tag: "CD",
					SynId: t_token.SynId,
				}
				i = result.index
				newTokenList = append(newTokenList, newToken)
			} else {
				newTokenList = append(newTokenList, t_token)
				i += 1
			}
		}
		return tokenizer.FilterOutSpaces(newTokenList)
	} else {
		return tokenList
	}
}

// load the complete library from file
func (g *GrammaryLibrary) initFromFile() error {
	filename := util.GetDataPath() + "/grammar/grammar-rules.txt"
	logger.Log.Info("NLU: loading %s", filename)
	text, err := util.LoadTextFile(filename)
	if err != nil { return err }
	g.InitFromString(text)
	logger.Log.Info("NLU: grammar loading done")
	return nil
}

// load the complete library from a string
func (g *GrammaryLibrary) InitFromString(grammarRules string) {
	if !g.initialised {
		g.Init()
		grammarPatternList := strings.Split(grammarRules, "\n")
		for _, pattern := range grammarPatternList {
			line := strings.TrimSpace(pattern)
			if len(line) > 0 && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") {
				lhs := g.processPattern(line)
				if ( lhs == nil ) {
					panic("invalid line in Grammar @ \"" + line + "\"")
				}
				if len(lhs.ConversionPattern) == 0 && len(lhs.Modifier) == 0 {
					if _, ok := g.GrammarMap[lhs.Name]; ok {
						panic("duplicate rule '" + lhs.Name + "'")
					}
				}
				// store in the maps
				if len(lhs.ConversionPattern) > 0 {
					g.GrammarConversionMap[lhs.Name] = lhs.ConversionPattern
				} else if len(lhs.Modifier) > 0 {
					g.GrammarModificationMap[lhs.Name] = lhs.Modifier
				} else {
					g.GrammarMap[lhs.Name] = lhs
				}
			}
		}
		// resolve reference to patterns internally
		g.resolveReferences()

		// setup first letter lookup
		g.setupFirstLetterLookup()
	}
}

// singleton access to the Grammar library system
var Grammar = GrammaryLibrary{}

// setup
func init() {
	// setup the grammar system
	err := Grammar.initFromFile()
	if err != nil {
		panic(err)
	}
}

