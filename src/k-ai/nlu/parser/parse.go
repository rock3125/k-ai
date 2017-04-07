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

package parser

import (
    "bytes"
    "net/http"
	"strconv"
	"encoding/json"
	"k-ai/nlu/lexicon"
	"k-ai/nlu/model"
	"k-ai/nlu/tokenizer"
	"k-ai/nlu/anaphora"
	"strings"
	"errors"
)

// spacy's default endpoint configuration
var SpacyEndpoint = "http://localhost:9000/parse"

// convert [][]model.Token to []model.Sentence
func convertTokensToSentenceList(tokenList [][]model.Token) ([]model.Sentence) {
	result_list := make([]model.Sentence,0)
	for _, ttoken := range tokenList {
		sent := model.Sentence{TokenList: ttoken}
		result_list = append(result_list, sent)
	}
	return result_list
}

// convert a json string to a parser result map
func JsonToParseResult(jsonStr string) (model.SentenceList) {
	res := model.SpacyList{} // json string to parser response
	json.Unmarshal([]byte(jsonStr), &res)
	sentence_list := make(model.SentenceList, 0)
	for _, tokenList := range res.SentenceList {
		sentence_list = append(sentence_list, model.Sentence{TokenList: tokenList})
	}
	return sentence_list
}

// post a request to the parser server (parsey)
func PostRequest(url string, text string) ([]model.Sentence, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(text))
	req.Header.Set("Content-Type", "text/plain")
	cl := strconv.Itoa(len(text))
	req.Header.Set("Content-Length", cl)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if strings.HasPrefix(buf.String(), "<!DOCTYPE HTML PUBLIC") {
		return nil, errors.New("parser in distress")
	}
	return JsonToParseResult(buf.String()), nil
}


// setup the semantics after a parse
func setupSemantics(sentenceList []model.Sentence) {
	for _, sentence := range sentenceList {
		for i, token := range sentence.TokenList {
			if len(sentence.TokenList[i].Semantic) == 0 { // only assign if not yet set
				sentence.TokenList[i].Semantic = lexicon.Lexi.GetSemantic(token.Text)
			}
		}
	}
}


// parse a piece of text and return its []model.Sentence
func ParseText(text string) ([]model.Sentence, error) {
	// parse the text
	sentence_list, err := PostRequest(SpacyEndpoint, text)
	if err != nil { return nil, err }

	// setup the semantics and sequences and grammar items
	// put spaces back into the sentence after spacy to avoid mistakes
	sentence_list = lexicon.Lexi.GetLongestWordSequenceForList(sentence_list) // 1. apply longest sentence

	setupSemantics(sentence_list)  // 2. setup semantics for items

	anaphora.LL.ResolvePronouns(sentence_list)  // 3. resolve third person pronouns

	return tokenizer.FilterOutSpacesForSentences(sentence_list), nil
}

// parse some text and return its JSON
func ParseTextToJson(text string) (string, error) {
	sentence_list, err := ParseText(text)
	if err != nil { return "", err }
	b, err := json.Marshal(sentence_list)
	if err != nil { return "", err }
	return string(b), nil
}

// parse a piece of text and return its list of tuple trees
func ParseTextToTupleTree(text string) (model.TreeList, error) {
	// convert the lot to tuple trees
	sentenceList, err := ParseText(text)
	if err != nil { return nil, err }
	resultList := make(model.TreeList,0)
	for _, sentence := range sentenceList {
		resultList = append(resultList, *model.SentenceToTuple(sentence))
	}
	return resultList, nil
}

// convert an existing list of sentences to a list of trees
func SentenceListToTupleTrees(sentence_list []model.Sentence) model.TreeList {
	resultList := make(model.TreeList,0)
	for _, sentence := range sentence_list {
		resultList = append(resultList, *model.SentenceToTuple(sentence))
	}
	return resultList
}

