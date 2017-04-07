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
	"github.com/gocql/gocql"
)

/**
 * perfom the serach using the objects itself
 * @param query_str a super search query string to be parsed
 * @return the matching indexes
 */
func SuperSearch(query_str string, origin string) (map[gocql.UUID][]model.IndexMatch, error) {
	query, err := parse_string(query_str)
	if err != nil { return nil, err }
	return doSearch(query, origin)
}


/**
 * perfom the serach using the objects itself
 * @param searchItem a series of search object
 * @return the matching indexes
 */
func doSearch(searchItem *SSTree, origin string) (map[gocql.UUID][]model.IndexMatch, error) {

	if searchItem != nil {

		switch searchItem.TType {
			case "and": {
				if searchItem.Left == nil || searchItem.Right == nil {
					return nil, errors.New("and: left or right nil")
				}
				set1, err := doSearch(searchItem.Left, origin)
				if err != nil { return nil, err }
				set2, err := doSearch(searchItem.Right, origin)
				if err != nil { return nil, err }
				return intersection(set1, set2)
			}
			case "and not": {
				if searchItem.Left == nil || searchItem.Right == nil {
					return nil, errors.New("and not: left or right nil")
				}
				set1, err := doSearch(searchItem.Left, origin)
				if err != nil { return nil, err }
				set2, err := doSearch(searchItem.Right, origin)
				if err != nil { return nil, err }
				return intersectionNot(set1, set2)
			}
			case "or": {
				if searchItem.Left == nil || searchItem.Right == nil {
					return nil, errors.New("or: left or right nil")
				}
				set1, err := doSearch(searchItem.Left, origin)
				if err != nil { return nil, err }
				set2, err := doSearch(searchItem.Right, origin)
				if err != nil { return nil, err }
				return union(set1, set2)
			}
			case "word": {
				return readIndexesForTerm(searchItem, origin)
			}
			default: {
				return nil, errors.New(fmt.Sprintf("unknown/unhandled super search type (%s)", searchItem.TType))
			}
		}
	}
	return nil, nil
}

// combine two lists of indexes into one
func combine(v1 []model.IndexMatch, v2 []model.IndexMatch) []model.IndexMatch {
	list := make([]model.IndexMatch,0)
	for _, item := range v1 {
		list = append(list, item)
	}
	for _, item := range v2 {
		list = append(list, item)
	}
	return list
}

/**
 * intersect two sets of indexes together into a single one
 * @param set1 first set
 * @param set2 second set
 * @return the intersection of the two sets (at url level)
 */
func intersection(set1 map[gocql.UUID][]model.IndexMatch, set2 map[gocql.UUID][]model.IndexMatch) (map[gocql.UUID][]model.IndexMatch, error) {
	result_set := make(map[gocql.UUID][]model.IndexMatch, 0)

	// either empty?
	if len(set1) == 0 || len(set2) == 0 {
		return result_set, nil
	}

	// first set intersected with second set
	for url1, value1_list := range set1 {
		if value2_list, ok := set2[url1]; ok {
			result_set[url1] = combine(value1_list, value2_list)
		}
	}
	return result_set, nil
}


/**
 * intersect two sets of indexes together for an AND NOT
 * @param set1 first set
 * @param set2 second set
 * @return the AND NOT intersection of the two sets (at url level)
 */
func intersectionNot(set1 map[gocql.UUID][]model.IndexMatch, set2 map[gocql.UUID][]model.IndexMatch) (map[gocql.UUID][]model.IndexMatch, error) {
	result_set := make(map[gocql.UUID][]model.IndexMatch, 0)

	// either empty - its set1
	if len(set2) == 0 {
		return set1, nil
	}

	// first set intersected with second set
	for url1, value1_list := range set1 {
		if _, ok := set2[url1]; !ok {
			result_set[url1] = value1_list
		}
	}
	return result_set, nil
}

/**
 * put two sets of indexes together into a single one
 * @param set1 first set
 * @param set2 second set
 * @return the union of the two sets
 */
func union(set1 map[gocql.UUID][]model.IndexMatch, set2 map[gocql.UUID][]model.IndexMatch) (map[gocql.UUID][]model.IndexMatch, error) {
	result_set := make(map[gocql.UUID][]model.IndexMatch, 0)

	if len(set1) == 0  {
		return set2, nil
	}

	if len(set2) == 0  {
		return set1, nil
	}

	// first set union-ed with second set
	for url1, value1_list := range set1 {
		if value2_list, ok := set2[url1]; ok {
			result_set[url1] = combine(value1_list, value2_list)
		} else {
			result_set[url1] = value1_list
		}
	}

	// take care of any items not yet in the set
	for url2, value2_list := range set2 {
		if _, ok := set1[url2]; !ok {
			result_set[url2] = value2_list
		}
	}

	return result_set, nil
}

// read indexes for a "word"
func readIndexesForTerm(item *SSTree, topic string) (map[gocql.UUID][]model.IndexMatch, error) {
	if item != nil {
		token_list := tokenizer.Tokenize(item.Word)
		return db_model.ReadIndexesWithFilterForTokens(token_list, topic,  0)
	}
	return nil, nil
}

