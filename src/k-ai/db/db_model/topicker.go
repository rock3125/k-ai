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
	"k-ai/nlu/model"
	"k-ai/db"
	"k-ai/nlu/lexicon"
	"sort"
	"errors"
	"github.com/gocql/gocql"
	"k-ai/logger"
	"fmt"
	"k-ai/util"
)

// a topic
type Topic struct {
	Topic string	`json:"name"`
	Body string		`json:"text"`
}

type Unindex struct {
	Sentence_id gocql.UUID
	Word string
	Tag string
}

type Topics []Topic

//////////////////////////////////////////////
// a topic name and its score
type TopicScore struct {
	Topic string
	Score float64
}

type TopicScores []TopicScore

// sort interface
func (slice TopicScores) Len() int {
	return len(slice)
}

// highest score first
func (slice TopicScores) Less(i, j int) bool {
	return slice[i].Score > slice[j].Score
}

func (slice TopicScores) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//////////////////////////////////////////////
// check the topic is valid
func isValidForTopic(t_token *model.Token) (bool, string, string) {
	tag := t_token.Tag
	if strings.HasPrefix(tag, "NN") || strings.HasPrefix(tag, "VB") {
		tag_str := "NN"
		if strings.HasPrefix(tag, "VB") {
			tag_str = "VB"
		}
		stemmed := lexicon.Lexi.GetStem(t_token.Text)
		return true, stemmed, tag_str
	}
	return false, "", ""
}

// return a list of topics from the db
func GetTopicList(prev string, page_size int) (Topics,error) {
	if page_size <= 0 {
		return nil, errors.New("GetTopicList() invalid parameters")
	}

	cols := []string{"topic", "body"}
	where_map := make(map[string]interface{})

	// setup first page
	if prev == "null" {
		prev = ""
	}
	cql_str := db.Cassandra.SelectPaginated("topic", cols, where_map, "topic", prev, page_size)
	iter := db.Cassandra.Session.Query(cql_str).Iter()

	list := make(Topics,0)

	var topic, body string

	for iter.Scan(&topic, &body) {
		list = append(list, Topic{Topic: topic, Body: body})
	}
	return list, iter.Close()
}

// save a topic to the db and index
func SaveTopic(topic string, body string, sentence_list []model.Sentence) error {
	if len(topic) == 0 || len(body) == 0 || len(sentence_list ) == 0 {
		return errors.New("SaveTopic() invalid parameters")
	}

	value_map := make(map[string]interface{})
	value_map["topic"] = topic
	value_map["body"] = body

	err := SaveText(sentence_list, topic)  // save the sentences themselves
	if err != nil { return err }

	err = db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("topic", value_map))
	if err != nil { return err }

	err = indexTopic(topic, sentence_list)
	if err != nil { return err }

	err = IndexText(topic, 0, sentence_list, 0.98) // index into factoid system too
	if err != nil { return err }

	err = IndexText("global", 0, sentence_list, 0.98) // index into factoid system too
	if err != nil { return err }

	return nil
}

// use the topic_unindexes to get sentence ids for a topic
func GetUnindexesForTopic(topic string) ([]Unindex, error) {
	unindex_list := make([]Unindex, 0)

	columns := []string{"sentence_id", "word", "tag"}
	whereMap := make(map[string]interface{},0)
	whereMap["topic"] = topic

	select_str := db.Cassandra.SelectPaginated("topic_unindex", columns, whereMap, "", nil, 0)
	iter := db.Cassandra.Session.Query(select_str).Iter()

	var word, tag string
	var sentence_id gocql.UUID
	for iter.Scan(&sentence_id,&word,&tag) {
		if !util.IsEmpty(&sentence_id) {
			uidx := Unindex{Word: word, Tag: tag}
			util.CopyUUID(&uidx.Sentence_id, &sentence_id)
			unindex_list = append(unindex_list, uidx)
		}
	}
	return unindex_list, iter.Close()
}

// remove a topic from the db and indexes
func DeleteTopic(topic string) error {
	if len(topic) == 0 {
		return errors.New("DeleteTopic() invalid parameters")
	}

	// remove the topic itself from the topic table
	where_map := make(map[string]interface{})
	where_map["topic"] = topic

	err := db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("topic", where_map))
	if err != nil { return err }

	// get the unindexes for further sentence removal
	unindex_list, err := GetUnindexesForTopic(topic)
	if err != nil { return err }

	seen := make(map[gocql.UUID]bool, 0)
	for _, unindex := range unindex_list {
		// remove each sentence exactly once in the text and indexes
		if _, ok := seen[unindex.Sentence_id]; !ok {

			seen[unindex.Sentence_id] = true
			err = DeleteText(&unindex.Sentence_id, topic)
			if err != nil {
				logger.Log.Warning(fmt.Sprintf("DeleteTopic: DeleteText: %s", err.Error()))
			}

			err = RemoveIndexes(unindex.Sentence_id, topic) // remove potential previous indexes
			if err != nil {
				logger.Log.Warning(fmt.Sprintf("DeleteTopic: RemoveIndexes(%s): %s", topic, err.Error()))
			}

			err = RemoveIndexes(unindex.Sentence_id, "global") // remove potential previous indexes
			if err != nil {
				logger.Log.Warning(fmt.Sprintf("DeleteTopic: RemoveIndexes(global): %s", err.Error()))
			}
		}
		// remove each index for this topic too
		deleteTopicIndex(topic, unindex.Word, unindex.Tag)
	}

	return nil
}

// index all sentences of a topic if the topic isn't an email
func indexTopic(topic_name string, sentence_list []model.Sentence) error {

	score_map := make(map[string]float64, 0)
	detection_map := make(map[string]bool,0)  // make sure we unindex items only once for optimization
	sentence_score := 1.0
	score_dropoff := 0.98
	for _, sentence := range sentence_list {
		for _, t_token := range sentence.TokenList {
			valid, stemmed, tag_str := isValidForTopic(&t_token)
			if valid && len(stemmed) > 0 && len(tag_str) > 0 {
				if _, ok := score_map[stemmed]; !ok {
					score_map[stemmed+":"+tag_str] = sentence_score
				} else {
					score_map[stemmed+":"+tag_str] += sentence_score
				}

				// add an unindex for this sentence word combination if not done so already
				key := stemmed + ":" + tag_str + ":" + sentence.Id.String()
				if _, ok := detection_map[key]; !ok {
					detection_map[key] = true // seen
					// add unindex for this sentence
					topicUnindexSet := make(map[string]interface{}, 0)
					topicUnindexSet["word"] = stemmed
					topicUnindexSet["tag"] = tag_str
					topicUnindexSet["topic"] = topic_name
					topicUnindexSet["sentence_id"] = sentence.Id
					err := db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("topic_unindex", topicUnindexSet))
					if err != nil {
						return err
					}
				}

			}
		} // for each token
		sentence_score *= score_dropoff  // score lessens as we move down the document
	}

	// add all items for scoring on this topic
	for key, value := range score_map {
		parts := strings.Split(key, ":")
		if len(parts[0]) > 0 && len(parts[1]) > 0 {
			topicSet := make(map[string]interface{}, 0)
			topicSet["word"] = parts[0]
			topicSet["tag"] = parts[1]
			topicSet["topic"] = topic_name
			topicSet["score"] = float32(value)
			err := db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("topic_index", topicSet))
			if err != nil { return err }
		}
	}
	return nil
}

// read a set of topic indexes (if available) for a given word
func readTopicIndexes(topic_set map[string]float64, word string, tag string) error {

	columns := []string{"topic", "score"}
	whereMap := make(map[string]interface{},0)
	whereMap["word"] = word
	whereMap["tag"] = tag

	select_str := db.Cassandra.SelectPaginated("topic_index", columns, whereMap, "", nil, 0)
	iter := db.Cassandra.Session.Query(select_str).Iter()

	var topic string
	var score float64
	for iter.Scan(&topic, &score) {
		if value, ok := topic_set[topic]; ok {
			topic_set[topic] = value + score
		} else {
			topic_set[topic] = score
		}
	}
	return iter.Close()
}

// read a set of topic indexes (if available) for a given word
func deleteTopicIndex(topic string, word string, tag string) error {

	whereMap := make(map[string]interface{},0)
	whereMap["word"] = word
	whereMap["tag"] = tag
	whereMap["topic"] = topic

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Delete("topic_index", whereMap))
}

// convert a topic set map to a list of ordered topics (highest first)
func toTopicList(topic_set map[string]float64) TopicScores {
	topics := make(TopicScores,0)
	for key, value := range topic_set {
		topics = append(topics, TopicScore{Topic: key, Score: value})
	}
	sort.Sort(topics)
	return topics
}

// given a set of tokens, get the best possible topic set (with scores)
func GetTopTopics(token_list []model.Token) (TopicScores, error) {
	topic_set := make(map[string]float64,0)
	for _, t_token := range token_list {
		valid, stemmed, tag_str := isValidForTopic(&t_token)
		if valid {
			err := readTopicIndexes(topic_set, stemmed, tag_str)
			if err != nil { return make(TopicScores,0), err }
		}
	}
	return toTopicList(topic_set), nil
}

