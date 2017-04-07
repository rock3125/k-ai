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

package db_model

import (
	"strings"
	"k-ai/db"
	"k-ai/nlu/lexicon"
	"k-ai/nlu/model"
	"github.com/gocql/gocql"
	"k-ai/util"
	"errors"
)


// index item - NB. keep in sync with IndexMatch !!!!
type Index struct {
	Sentence_id gocql.UUID	// the sentence owner id

	Word        string 		// the word, main index
	Tag         string 		// the Penn tag of the word
	Shard       int    		// shard spreading across systems

	Offset      int    		// unique offset for repeating words
	Topic		string 		// what spawned it, what is the sentence owner?
	Score		float64		// the value of this index relative to others
}

// unindex item
type UnIndex struct {
	Sentence_id gocql.UUID	// the sentence owner id
	Word        string 		// the word, main index
	Shard       int    		// shard spreading across systems
}


// add an index into the system, and its unindex equivalent for later removal
func addIndex(sentence_id *gocql.UUID, word string, tag string, shard int,
				topic string, offset int, score float64) error {
	// add the index
	indexValueSet := make(map[string]interface{}, 0)
	indexValueSet["sentence_id"] = sentence_id
	indexValueSet["word"] = word
	indexValueSet["tag"] = tag
	indexValueSet["shard"] = shard

	indexValueSet["offset"] = offset
	indexValueSet["topic"] = topic
	indexValueSet["score"] = score

	err := db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("word_index", indexValueSet))
	if err != nil { return err }

	// add the unindex
	// url text, origin text, shard int, word text, kb text,
	// primary key((url,origin,kb), word, shard)
	unindexValueSet := make(map[string]interface{}, 0)
	unindexValueSet["sentence_id"] = sentence_id
	unindexValueSet["word"] = word
	unindexValueSet["shard"] = shard
	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("word_unindex", unindexValueSet))
}

// index a text string into the system
func IndexText(topic string, shard int, sentence_list []model.Sentence, score_dropoff float64) error {

	if len(topic) > 0 && len(sentence_list) > 0 {
		offset := 0
		score := 1.0

		for _, sentence := range sentence_list {

			if util.IsEmpty(&sentence.Id) { return errors.New("invalid guid for sentence") }

			for _, t_token :=  range sentence.TokenList {

				stemmed := lexicon.Lexi.GetStem(t_token.Text)
				// never index the auxiliary verbs
				if len(stemmed) > 0 && !lexicon.Lexi.IsUndesirable(stemmed) && t_token.Dep != "aux" { // only index valid words

					err := addIndex(&sentence.Id, stemmed, t_token.Tag, shard, topic, offset, score)
					if err != nil { return err }

					////////////////////////////////////////////////////////////////////////
					// also index sub parts of compound words like "New York" -> "New" and "York"

					// spaces in the token (multi word index)
					parts := make([]string, 0)
					if strings.Contains(t_token.Text, " ") {
						for _, str := range strings.Split(t_token.Text, " ") {
							parts = append(parts, str)
						}
					}
					if strings.Contains(t_token.Text, "-") {
						for _, str := range strings.Split(t_token.Text, "-") {
							parts = append(parts, str)
						}
					}
					if len(parts) > 1 {
						for _, part := range parts {
							part_lcase := strings.ToLower(strings.TrimSpace(part))
							if len(part_lcase) > 0 && !lexicon.Lexi.IsUndesirable(part_lcase) {
								// add an index for parts of the words
								err := addIndex(&sentence.Id, part_lcase, t_token.Tag, shard,topic, offset, score * 0.5)
								if err != nil { return err }
							}
						}
					}

					////////////////////////////////////////////////////////////////////////
					// also index semantics of words as a reference, e.g. "New York":city -> index city reference

					if len(t_token.Semantic) > 0 {
						token_semantic := strings.ToLower(strings.TrimSpace(t_token.Semantic))
						if len(token_semantic) > 0 && !lexicon.Lexi.IsUndesirable(token_semantic) {
							// add an index for parts of the words
							err := addIndex(&sentence.Id, token_semantic, t_token.Tag, shard,topic, offset, score * 0.5)
							if err != nil { return err }
						}
					}

				} // if valid word for index

				offset += 1
			} // for each token

			score *= score_dropoff

		} // for each sentence

	}
	return nil
}

// read indexes using word and meta-data fields
// word text, shard int, tag text, url text, kb text, offset int, meta_c_type int,
func readIndexes(word string, topic string, shard int) ([]Index,error) {
	return_list := make([]Index,0)

	columns := []string{"offset", "sentence_id", "tag", "score"}
	whereMap := make(map[string]interface{},0)
	whereMap["word"] = strings.ToLower(word)
	whereMap["shard"] = shard
	whereMap["topic"] = topic

	select_str := db.Cassandra.SelectPaginated("word_index", columns, whereMap, "", nil, 0)
	iter := db.Cassandra.Session.Query(select_str).Iter()

	var offset int
	var score float64
	var tag string
	var sentence_id gocql.UUID

	for iter.Scan(&offset, &sentence_id, &tag, &score) {
		return_list = append(return_list, Index{Topic: topic, Offset: offset, Score: score,
			Sentence_id: sentence_id, Shard: shard, Tag: tag, Word: word})
	}
	return return_list, iter.Close()
}

// return true if a verb is a verb and a noun is a noun - otherwise we don't care and assume its ok
func compatible_tag(tag1 string, tag2 string) bool {
	if len(tag1) >= 2 && len(tag2) >= 2 {
		if strings.HasPrefix(tag1, "NN") || strings.HasPrefix(tag2, "NN") {
			// either matching nouns, or one of them is an adjective (e.g. blue can be a noun or JJ depending on use)
			return tag1[0:2] == tag2[0:2] || tag1[0:2] == "JJ" || tag2[0:2] == "JJ"
		}
		if strings.HasPrefix(tag1, "VB") || strings.HasPrefix(tag2, "VB") {
			return tag1[0:2] == tag2[0:2]
		}
	}
	return true  // all other cases are ok, including missing tags
}

// return how many valid tokens there are in the tokenList
func GetNumSearchTokens(token_list []model.Token) int {
	count := 0
	for _, t_token := range token_list { // for each token
		stemmed := lexicon.Lexi.GetStem(t_token.Text) // unstem it
		if len(stemmed) > 0 && !lexicon.Lexi.IsUndesirable(stemmed) { // must be index-able
			count += 1
		}
	}
	return count
}

/**
 * read a set of indexes using words as a filter for a specific set of meta-data
 * @param organisation_id the id of the organisation to read from
 * @param tokenList a parsed + filtered set of tokens to search through the indexes
 * @param shard the shard of the index
 * @return a list of URLs that matched
 */
func ReadIndexesWithFilterForTokens(token_list []model.Token, topic string, shard int) (map[gocql.UUID][]model.IndexMatch, error) {
	combined_indexes := make(map[gocql.UUID][]model.IndexMatch, 0)
	i := 0
	for _, t_token := range token_list { // for each token
		stemmed := lexicon.Lexi.GetStem(t_token.Text) // unstem it
		if len(stemmed) > 0 && !lexicon.Lexi.IsUndesirable(stemmed) { // must be index-able
			indexes, err := readIndexes(stemmed, topic, shard) // read the indexes
 			if err != nil {
				return nil, err
			}
			if len(indexes) == 0 {
				return make(map[gocql.UUID][]model.IndexMatch, 0), nil // fail if nothing returned at any one point
			} else if i == 0 { // first index all items are just added
				for _, index := range indexes {
					if compatible_tag(t_token.Tag, index.Tag) {
						if list, ok := combined_indexes[index.Sentence_id]; ok {
							list = append(list, *model.Convert(index.Sentence_id, index.Word, index.Tag, index.Shard, index.Offset,
																index.Topic, index.Score, i))
						} else {
							combined_indexes[index.Sentence_id] = make([]model.IndexMatch, 0)
							combined_indexes[index.Sentence_id] = append(combined_indexes[index.Sentence_id],
								*model.Convert(index.Sentence_id, index.Word, index.Tag, index.Shard, index.Offset, index.Topic, index.Score, i))
						}
					}
				}

			} else if i > 0 {
				// all other rounds are intersection rounds
				new_combined_indexes := make(map[gocql.UUID][]model.IndexMatch, 0)
				for _, index := range indexes {
					if compatible_tag(t_token.Tag, index.Tag) {
						if list, ok := combined_indexes[index.Sentence_id]; ok {
							new_combined_indexes[index.Sentence_id] = make([]model.IndexMatch, 0)
							for _, item := range list {
								new_combined_indexes[index.Sentence_id] = append(new_combined_indexes[index.Sentence_id], item)
							}
							new_combined_indexes[index.Sentence_id] = append(new_combined_indexes[index.Sentence_id],
								*model.Convert(index.Sentence_id, index.Word, index.Tag, index.Shard, index.Offset, index.Topic, index.Score, i))
						}
					}
				}
				combined_indexes = new_combined_indexes
				if len(combined_indexes) == 0 {
					return make(map[gocql.UUID][]model.IndexMatch, 0), nil // failed after combining indexes - nothing left
				}
			}

			i += 1

		} // if valid
	}
	return combined_indexes, nil
}


// read the list of un-indexes for a url / origin / kb
func readUnindexes(sentence_id gocql.UUID) ([]UnIndex,error) {
	return_list := make([]UnIndex,0)

	columns := []string{"word", "shard"}
	whereMap := make(map[string]interface{},0)
	whereMap["sentence_id"] = sentence_id

	select_str := db.Cassandra.SelectPaginated("word_unindex", columns, whereMap, "", nil, 0)
	iter := db.Cassandra.Session.Query(select_str).Iter()

	var word string
	var shard int

	for iter.Scan(&word, &shard) {
		return_list = append(return_list, UnIndex{Sentence_id: sentence_id, Shard: shard, Word: word })
	}
	return return_list, iter.Close()
}


// delete an index item: word,origin,kb,shard
func deleteIndex(topic string, unindex *UnIndex) error {
	where_map := make(map[string]interface{})
	where_map["topic"] = topic
	where_map["word"] = unindex.Word
	where_map["shard"] = unindex.Shard
	where_map["sentence_id"] = unindex.Sentence_id
	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("word_index", where_map))
}

// delete an unindex item: url,origin,kb
func deleteUnIndex(sentence_id gocql.UUID) error {
	where_map := make(map[string]interface{})
	where_map["sentence_id"] = sentence_id
	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("word_unindex", where_map))
}

// remove indexes for a given sentence
func RemoveIndexes(sentence_id gocql.UUID, topic string) error {
	// read the unindexes
	unindex_list, err := readUnindexes(sentence_id)
	if err != nil { return err }

	for _, unindex := range unindex_list {
		// delete each word index
		err = deleteIndex(topic, &unindex)
		if err != nil { return err }
	}
	return deleteUnIndex(sentence_id)
}

// remove indexes for a given list of sentences
func RemoveIndexesForSentenceList(sentence_list []model.Sentence, topic string) error {
	for _, sentence := range sentence_list {
		if !util.IsEmpty(&sentence.Id) {
			err := RemoveIndexes(sentence.Id, topic)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Amend any pronoun references with replacement as an NNP token (look it's semantic up in the lexicon)
func resolvePronounReferences(token_list []model.Token, replacement string, pronouns...string) []model.Token {
	// does the desired pronoun occur?
	new_token_list := make([]model.Token,0)
	var pronoun_token *model.Token
	for _, t_token := range token_list {
		str := strings.ToLower(t_token.Text)
		found := false
		for _, pron := range pronouns {
			if pron == str && (t_token.Tag == "PRP" || t_token.Tag == "PRP$") {
				found = true
				break
			}
		}
		if found { // process the pronoun
			if pronoun_token == nil { // get the replacement token
				pronoun_token = &model.Token{Text: replacement, Semantic: lexicon.Lexi.GetSemantic(replacement),
					Tag:                           "NNP", AncestorList: t_token.AncestorList, Dep: t_token.Dep, Index: t_token.Index, SynId: t_token.SynId}
			}
			new_token_list = append(new_token_list, t_token)
			new_token_list = append(new_token_list, model.Token{Text: "[", Index: t_token.Index, Tag: "-LRB-"})
			new_token_list = append(new_token_list, *pronoun_token)
			new_token_list = append(new_token_list, model.Token{Text: "]", Index: t_token.Index, Tag: "-RRB-"})
		} else {
			new_token_list = append(new_token_list, t_token)
		}
	}
	return new_token_list
}

// resolve first and second person pronoun references
func ResolveFirstAndSecondPerson(first_person string, second_person string, sentence *model.Sentence) {
	// replace you, your, yourself pronoun references with KAI
	sentence.TokenList = resolvePronounReferences(sentence.TokenList, second_person,
		"you", "your", "yourself", "yours")

	sentence.TokenList = resolvePronounReferences(sentence.TokenList, first_person,
		"i", "my", "myself", "me", "mine")
}

