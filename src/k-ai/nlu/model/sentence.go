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

package model

import (
	"strings"
	"github.com/gocql/gocql"
	"k-ai/util"
)

type Sentence struct {
	Id gocql.UUID           `json:"id"`
	Topic string			`json:"topic"`
	TokenList []Token       `json:"tokenList"`
}

type SentenceList []Sentence

// len for sort interface
func (s Sentence) Len() int {
    return len(s.TokenList)
}

// generate a random id for this sentence
func (s *Sentence) RandomId() error {
	uuid, err := gocql.RandomUUID()
	if err != nil { return err }
	util.CopyUUID(&s.Id, &uuid)
	return nil
}

// less for sorting, sort by token index
func (s Sentence) Less(i, j int) bool {
	return s.TokenList[i].Index < s.TokenList[j].Index;
}

// sort interface
func (s Sentence) Swap(i, j int) {
    s.TokenList[i], s.TokenList[j] = s.TokenList[j], s.TokenList[i]
}

func (s Sentence) ToString() (string) {
	str := ""
	for _, token := range s.TokenList {
		str += token.Text
		if len(token.Semantic) > 0 {
			s_str := token.Semantic
			index := strings.Index(s_str, ".")  // remove . substring of grammar item if present
			if index > 0 {
				s_str = s_str[:index]
			}
			str += "{" + s_str + "} "
		} else if strings.Contains(token.Dep, "sub") || strings.Contains(token.Dep, "obj") {
			str += "{" + token.Dep + "} " // display grammar rule
		} else {
			str += " "
		}
	}
	return strings.TrimSpace(str)
}

// is this sentence a question sentence?
// for now - check last token is a question mark
// or the first token is who|what|where|why|when|how
// or the first token is a aux verb: be | have | do
func (s Sentence) IsQuestion() bool {
	if len(s.TokenList) > 1 {  // min sentence size
		if s.TokenList[len(s.TokenList)-1].Text == "?" {
			return true
		}
		token0 := s.TokenList[0]
		if strings.HasPrefix(token0.Tag, "VB") {
			lwr_verb := strings.ToLower(token0.Text)
			if lwr_verb == "do" || lwr_verb == "did" || lwr_verb == "does" ||
				lwr_verb == "are" || lwr_verb == "is" || lwr_verb == "was" ||
				lwr_verb == "have" || lwr_verb == "had" {
				return true
			}
		}
		// any of the tokens is a WDT tag?
		for _, token := range s.TokenList {
			if token.Tag == "WDT" || token.Tag == "WP" {
				return true
			}
		}
	}
	return false
}

// is this sentence an imperative statement (command)
// no question mark at the end, first word is a verb other than an AUX
func (s Sentence) IsImperative() bool {
	if len(s.TokenList) > 1 { // min sentence size
		if s.TokenList[len(s.TokenList)-1].Text == "?" {
			return false
		}
		token0 := s.TokenList[0]
		if strings.HasPrefix(token0.Tag, "VB") {
			if token0.Tag == "VB" || token0.Tag == "VBP" {
				return true
			}
		}
	}
	return false
}


// does the sentence have a verb?
func (s Sentence) HasVerb() bool {
	if len(s.TokenList) > 1 { // min sentence size
		num_verbs := 0
		for _, token := range s.TokenList {
			if strings.HasPrefix(token.Tag, "VB") {
				num_verbs += 1
			}
		}
		return num_verbs > 0
	}
	return false
}

