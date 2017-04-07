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
	"k-ai/nlu/model"
	"errors"
	"fmt"
	"k-ai/db/db_model"
	"k-ai/nlu/tokenizer"
	"strings"
)

/**
 * Created by peter on 25/04/16.
 *
 * parse super search statements and convert them to actionable search items
 *
 * grammar:
 *
 * where ->  ['exact'] 'any' '(' text ')'
 *           ['exact'] 'person' '(' text ')' |
 *           ['exact'] 'location' '(' text ')' |
 *           'url'  '('  text  ')'  |
 *           'date' 'between' time 'and' time |
 *           'date' 'before' time |
 *           'date' 'after' time |
 *           'date' 'exact' time  |
 *           where 'and' where |
 *           where 'and' 'not' where |
 *           where 'or' where |
 *           '(' where ')'
 *
 * text ->   ...  |
 *           ... , tag
 *
 * time ->  yyyy-mm-dd hh ':' mm |
 *          yyyy-mm-dd |
 *          yyyy-mm |
 *          yyyy
 *
 * tag -> verb | noun | adjective | proper noun | penn-tag
 *
 */

/**
 * convert a super search text to a search item set
 *
 * @param index the index into the token list
 * @param tokenList a list of tokens to be converted / parsed
 * @return the parsed search system
 */
func parse_string(query string) (*SSTree, error) {
	return parse(0, tokenizer.Tokenize(query))
}


/**
 * convert a super search text to a search item set
 *
 * @param index the index into the token list
 * @param tokenList a list of tokens to be converted / parsed
 * @return the parsed search system
 */
func parse(index int, tokenList []model.Token) (*SSTree, error) {
	if len(tokenList) > 0 && index < len(tokenList) {

		var item *SSTree
		var err error
		tokenWithIndex := getNextSkippingSpace(index, tokenList)
		if tokenWithIndex == nil {
			return nil, errors.New(fmt.Sprintf("expected statement start @ %d", index))
		}
		tokenStr := tokenWithIndex.Text
		switch tokenStr {
			case "exact": {
				item, err = parseExactWord(tokenWithIndex.Index, tokenList)
				if err != nil { return nil, err }
			}
			case "location", "person", "any": {
				item, err = parseWord(tokenWithIndex.Index, tokenList, tokenStr, false)
				if err != nil { return nil, err }
			}
			case "(": {
				item, err = parseBrackets(index, tokenList)
				if err != nil { return nil, err }
			}
		}

		if item == nil {
			return nil, errors.New(fmt.Sprintf("unknown token @ %d", index))
		}

		index = item.Offset  // update index to next token

		// or/and ?
		if index < len(tokenList) {
			tokenWithIndex = getNextSkippingSpace(index, tokenList)
			if tokenWithIndex != nil && (tokenWithIndex.Text == "or" || tokenWithIndex.Text == "and" ) {
				return parseAndOr( item, index, tokenList )
			}
		}

		// just return the item itself
		return item, nil
	}
	return nil, nil
}


/**
 * get the next token skipping any spaces and return its index
 * @param index the index to start @
 * @param tokenList the token-list to scan
 * @return the token with its index
 */
func getNextSkippingSpace( index int, tokenList []model.Token ) *SSTokenWithIndex {
	// skip any white-spaces automatically
	if index < len(tokenList) {
		next := &tokenList[index]
		for next != nil && next.Text == " " && index < len(tokenList) {
			index++
			if index < len(tokenList) {
				next = &tokenList[index]
			} else {
				next = nil
			}
		}
		if next != nil {
			return &SSTokenWithIndex{*next, index + 1}
		}
	}
	return nil
}



/**
 * 'exact' 'word' '(' text ')'
 * @param index the offset into the array
 * @param tokenList the array of items
 * @return a parsed item if successful or null
 */
func parseExactWord(index int, tokenList []model.Token) (*SSTree, error) {
	if index < len(tokenList) {

		tokenWithIndex := getNextSkippingSpace(index, tokenList)
		if tokenWithIndex == nil {
			return nil, errors.New(fmt.Sprintf("expression 'exact' must be followed by other tokens @ %d", index))
		}
		return parseWord( tokenWithIndex.Index, tokenList, tokenWithIndex.Text, true)
	} else {
		return nil, errors.New(fmt.Sprintf("expression 'exact' must be followed by other tokens @ %d", index))
	}
}


/**
 * make sure the token @ index is word
 * @param index the index of the system
 * @param tokenList the  list of tokens / stream
 * @param words the word(s) to check
 */
func getNextCompulsary( index int, tokenList []model.Token, words...string ) (int, error) {

	wordList := ""
	for _, word := range words {
		if len(wordList) > 0 {
			wordList +=  ", "
		}
		wordList += word
	}

	if index < len(tokenList) {
		// skip any white-spaces automatically
		next := getNextSkippingSpace( index, tokenList )
		if next == nil {
			return 0, errors.New(fmt.Sprintf("expected token(s) %s @ %d'", wordList, index))
		}

		// check it is the word(s)
		found := false
		for _, word := range words {
			if word == next.Text {
				found = true
				break
			}
		}

		if !found {
			return 0, errors.New(fmt.Sprintf("expected token(s) %s @ %d'", wordList, index))
		}
		return next.Index, nil

	} else {
		return 0, errors.New(fmt.Sprintf("expected token(s) %s @ %d'", wordList, index))
	}
}


/**
 * ' text ' |
 * ' text ' , tag
 * @param index the offset into the array
 * @param tokenList the array of items
 * @return a parsed item if successful or null
 */
func parseText(index int, tokenList []model.Token) (*SSTree, error) {
	if index < len(tokenList) {
		text := ""
		ttoken := &tokenList[index]
		for ttoken != nil && ttoken.Text != ")" && ttoken.Text != "," && index < len(tokenList) {
			text += ttoken.Text
			index += 1
			if index < len(tokenList) {
				ttoken = &tokenList[index]
			} else {
				ttoken = nil
			}
		}
		if ttoken == nil || (ttoken.Text != ")" && ttoken.Text != ",") {
			return nil, errors.New(fmt.Sprintf("unterminated text @ %d", index - 1))
		}
		index += 1

		// is it followed by an optional tag?
		if ttoken.Text == "," {
			next := getNextSkippingSpace(index, tokenList)
			if next == nil { // || next.getItem().getType() != TokenizerConstants.Type.Text )
				return nil, errors.New("invalid token following text , penn-type")
			}

			index = next.Index
			tag := next.Text
			isPennTag := model.IsPennType(tag)
			if tag != "noun" && tag != "proper noun" && tag != "adjective" && tag != "verb" && !isPennTag {
				return nil, errors.New(fmt.Sprintf("invalid token following text , penn-type: %s", tag))
			}
			if isPennTag {
				tag = strings.ToUpper(tag)
			}
			return &SSTree{TType: "word", Offset: index, Index: db_model.Index{ Word: text, Tag: tag}}, nil

		} else {
			return &SSTree{TType: "word", Offset: index -1, Index: db_model.Index{ Word: text, }}, nil
		}
	}
	return nil, nil
}


/**
 * semantic '(' text ')'
 * @param index the offset into the array
 * @param tokenList the array of items
 * @return a parsed item if successful or null
 */
func parseWord(index int, tokenList []model.Token, semantic string, exact bool) (*SSTree, error) {
	if index < len(tokenList) {
		index, err := getNextCompulsary(index, tokenList, "(" )
		if err != nil { return nil, err }
		item, err := parseText(index, tokenList)
		if err != nil { return nil, err }
		if item == nil || item.TType != "word" {
			return nil, errors.New(fmt.Sprintf("expression word( must be followed by text @ %d", index))
		}

		var pennTag string
		switch semantic {
			case "any": {
				semantic = ""  // effectively an empty semantic
				pennTag = ""
			}
			case "person": {
				if len(item.Tag) > 0 {
					return nil, errors.New("'person' field cannot be followed by a tag specifier, is assumed NNP")
				}
				pennTag = "NNP"
			}
			case "location": {
				if len(item.Tag) > 0 {
					return nil, errors.New("'location' field cannot be followed by a tag specifier, is assumed NNP")
				}
				pennTag = "NNP"
			}
			default: {
				return nil, errors.New(fmt.Sprintf("unknown meta-data tag %s", semantic))
			}
		}
		item.Exact = exact
		item.Semantic = semantic
		if len(pennTag) > 0 { // don't overwrite
			item.Tag = pennTag
		}

		index = item.Offset
		index, err = getNextCompulsary(index, tokenList, ")" )
		if err != nil { return nil, err }
		item.Offset = index
		return item, nil
	}
	return nil, nil
}


/**
 * '(' ssearch ')'
 * @param index the offset into the array
 * @param tokenList the array of items
 * @return a parsed item if successful or null
 */
func parseBrackets(index int, tokenList []model.Token) (*SSTree, error) {
	if index < len(tokenList) {
		index, err := getNextCompulsary( index, tokenList, "(" )
		if err != nil { return nil, err }
		item, err := parse( index, tokenList )
		if err != nil { return nil, err }
		if item == nil {
			return nil, errors.New(fmt.Sprintf("expected 'expression )' after '(' @ %d", index ))
		}
		index = item.Offset

		// next token can be either 'and', 'or' or ')'
		if index < len(tokenList) {
			tokenWithIndex := getNextSkippingSpace(index, tokenList)
			if tokenWithIndex != nil {
				next := tokenWithIndex
				if len(next.Text) == 0 || !(next.Text == "or" || next.Text == "and" || next.Text == ")") {
					return nil, errors.New(fmt.Sprintf("expected token 'and'/'or'/')' @ %d", index ))
				}

				// update and / or parsing
				if next.Text == "or" || next.Text == "and" {
					item, err = parseAndOr(item, index, tokenList)
					if err != nil { return nil, err }
				}
				index, err = getNextCompulsary(index, tokenList, ")")
				if err != nil { return nil, err }
				item.Offset = index
				return item, nil

			} else {
				return nil, errors.New(fmt.Sprintf("expected ')' @ %d", index ))
			}
		} else {
			return nil, errors.New(fmt.Sprintf("expected ')' @ %d", index ))
		}
	}
	return nil, nil
}


/**
 * and / or parser helper
 * @param item1 the first item
 * @param index the current index at which 'or' or 'and' was found
 * @param tokenList the list of tokens
 * @return the and/or parsed entity with updated index
 */
func parseAndOr(item1 *SSTree, index int, tokenList []model.Token) (*SSTree, error) {
	// or/and ?
	if index < len(tokenList) {
		tokenWithIndex := getNextSkippingSpace(index, tokenList)
		if tokenWithIndex != nil {
			tokenStr := tokenWithIndex.Text
			switch tokenStr {
				case "or": {
					item2, err := parse(tokenWithIndex.Index, tokenList)
					if err != nil { return nil, err }
					if item2 == nil {
						return nil, errors.New("'or' missing rhs expression")
					}
					return &SSTree{Offset: item2.Offset, TType: "or", Left: item1, Right: item2}, nil
				}
				case "and": {
					// and not?
					next := getNextSkippingSpace(tokenWithIndex.Index, tokenList)
					if next != nil && next.Text == "not" {
						item2, err := parse(next.Index, tokenList)
						if err != nil { return nil, err }
						if item2 == nil {
							return nil, errors.New("'and' missing rhs expression")
						}
						return &SSTree{Offset: item2.Offset, TType: "and not", Left: item1, Right: item2}, nil
					} else {
						item2, err := parse(tokenWithIndex.Index, tokenList)
						if err != nil { return nil, err }
						if item2 == nil {
							return nil, errors.New("'and' missing rhs expression")
						}
						return &SSTree{Offset: item2.Offset, TType: "and", Left: item1, Right: item2}, nil
					}
				}
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("invalid 'and' / 'or' block @ %d", index))
}

