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
	"k-ai/util"
	"strings"
	"k-ai/logger"
	"k-ai/nlu/model"
	"k-ai/nlu/tokenizer"
	"sync"
)


type SLexicon struct {
	initialised  bool                       // init status

	plural       map[string]string        // lwr(plural) -> lwr(singular)
	verb         map[string]string        // lwr(non-vb-verb) -> lwr(vb_verb)
	Semantic     map[string]string        // lwr(noun) -> semantic_for_noun
	LWord        map[string][]model.Token // longest word multiple nouns
	Undesirables map[string]bool          // list of undesirable words

	stemSet      map[string]map[string]bool // stemmed word -> list of related words
	synonymSet   map[string]map[string]bool // word -> list of synonyms

	seen		map[string] bool			// temp map for speeding up loading

	sync.Mutex
}

// the lexicon global placeholder
var Lexi = SLexicon{}

// set a related word if dne
func (l *SLexicon) add_stem_word(baseWord string, relatedWord string) {
	if baseWord != relatedWord {
		if map1, ok := l.stemSet[baseWord]; ok {
			if _, found := map1[relatedWord]; !found {
				map1[relatedWord] = true
			}
		} else {
			l.stemSet[baseWord] = make(map[string]bool, 0)
			l.stemSet[baseWord][relatedWord] = true
		}
	}
}

// set a synonym if dne
func (l *SLexicon) add_synonym(word1 string, word2 string) {
	if word1 != word2 {
		if map1, ok := l.synonymSet[word1]; ok {
			if _, found := map1[word2]; !found {
				map1[word2] = true
			}
		} else {
			l.synonymSet[word1] = make(map[string]bool, 0)
			l.synonymSet[word1][word2] = true
		}
	}
}

// load the plural to singular map
func (l *SLexicon) loadPlurals(dataDir string) error {
	l.plural = make(map[string]string,0)
	file_contents, err := util.LoadTextFile(dataDir + "/lexicon/plurals.txt")
	if err != nil { return err }

	for _, line := range strings.Split(file_contents, "\n") {
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			singular := strings.ToLower(parts[0])
			plural := strings.ToLower(parts[1])
			l.plural[plural] = singular                 // plural -> singular
			l.add_stem_word(singular, plural)           // record relationship
			l.testAndAddCompoundWordWithCache(singular) // setup longest word if applicable
			l.testAndAddCompoundWordWithCache(plural)
		}
	}
	return nil
}

// load the synonyms
func (l *SLexicon) loadSynonyms(dataDir string) error {
	l.synonymSet = make(map[string]map[string]bool,0) // setup synonym lookup
	file_contents, err := util.LoadTextFile(dataDir + "/lexicon/synonyms.txt")
	if err != nil { return err }

	for _, line := range strings.Split(file_contents, "\n") {
		parts := strings.Split(line, ",")
		if len(parts) > 1 {
			for i, word1 := range parts {
				for j, word2 := range parts {
					if i != j {
						l.add_synonym(word1, word2)
						l.testAndAddCompoundWordWithCache(word1) // setup longest word if applicable
						l.testAndAddCompoundWordWithCache(word2)
					}
				}
			}
		}
	}
	return nil
}

// load the plural to singular map
func (l *SLexicon) loadVerbs(dataDir string) error {
	l.verb = make(map[string]string,0)
	file_contents, err := util.LoadTextFile(dataDir + "/lexicon/verbs.txt")
	if err != nil { return err }

	for _, line := range strings.Split(file_contents, "\n") {
		parts := strings.Split(line, "|")
		if len(parts) >= 6 {
			lwrStr := strings.ToLower(parts[0])
			l.testAndAddCompoundWordWithCache(lwrStr)
			for i := 1; i < len(parts); i++ {
				conjugateStr := strings.ToLower(parts[i])
				if conjugateStr != lwrStr {
					l.verb[conjugateStr] = lwrStr
					l.add_stem_word(lwrStr, conjugateStr)           // record relationship
					l.testAndAddCompoundWordWithCache(conjugateStr) // setup longest word if applicable
				}
			}
		}
	}
	return nil
}

// see if a word is a compound word and add it tot he system for longest word fixing
func (l *SLexicon) testAndAddCompoundWordWithCache(word_str string) {
	if _, ok := l.seen[word_str]; !ok {
		l.seen[word_str] = true
		isCompound := false
		for _, ch := range word_str {
			if ch == ' ' || ch == '-' {
				isCompound = true
				break
			}
		}
		if isCompound {
			parts := tokenizer.FilterOutSpaces(tokenizer.Tokenize(word_str))
			if len(parts) > 1 {
				l.LWord[strings.ToLower(word_str)] = parts
			}
		}
	}
}

// see if a word is a compound word and add it tot he system for longest word fixing
func (l *SLexicon) testAndAddCompoundWord(word_str string) {
	isCompound := false
	for _, ch := range word_str {
		if ch == ' ' || ch == '-' {
			isCompound = true
			break
		}
	}
	if isCompound {
		parts := tokenizer.FilterOutSpaces(tokenizer.Tokenize(word_str))
		if len(parts) > 1 {
			l.LWord[strings.ToLower(word_str)] = parts
		}
	}
}

// load the longest space combinations of words
func (l *SLexicon) loadLongestWords(dataDir string) error {
	file_contents, err := util.LoadTextFile(dataDir + "/lexicon/compound_nouns.txt")
	if err != nil { return err }
	for _, line := range strings.Split(file_contents, "\n") {
		l.testAndAddCompoundWordWithCache(strings.TrimSpace(line)) // setup longest word if applicable
	}
	return nil
}

// setup the lexicon's undesirables
func (l *SLexicon) setupUndesirables() {
	l.Undesirables = make(map[string]bool,0)
	for _, str := range undesirableList {
		l.Undesirables[str] = true
	}
}

// setup the lexicon
func (l *SLexicon) initFromFile() error {
	if !l.initialised {
		l.initialised = true
		dataDir := util.GetDataPath()
		logger.Log.Info("NLU: loading %s", dataDir)

		l.stemSet = make(map[string]map[string]bool,0) // setup stem word lookup
		l.LWord = make(map[string][]model.Token,0) // setup longest word

		l.seen = make(map[string]bool,0) // temp for speeding up loading

		err := l.loadPlurals(dataDir)
		if err != nil { return err }

		err = l.loadVerbs(dataDir)
		if err != nil { return err }

		err = l.loadSemantics(dataDir)
		if err != nil { return err }

		err = l.loadLongestWords(dataDir)
		if err != nil { return err }

		l.setupUndesirables()
		err = l.loadSynonyms(dataDir)
		if err != nil { return err }

		l.seen = nil // release map

		logger.Log.Info("NLU: lexicon loaded %d stem items, %d compound words, %d semantics, %d undesirables, and %d synonyms", len(l.plural)+len(l.verb),
			len(l.LWord), len(l.Semantic), len(l.Undesirables), len(l.synonymSet))

		// apply any updates from the UI to the semantics system
		l.applySemanticUpdates()

		return err
	}
	return nil
}

// return the stem of a word if it exists, otherwise the word to lower case is returned
func (l *SLexicon) GetStem(word string) string {
	lwrStr := strings.ToLower(word)
	if l.initialised {
		if val, ok := l.plural[lwrStr]; ok {
			return val
		}
		if val, ok := l.verb[lwrStr]; ok {
			return val
		}
	}
	return lwrStr
}

// return true if this word is in the undesirables list
func (l *SLexicon) IsUndesirable(word string) bool {
	_, ok := l.Undesirables[strings.ToLower(word)]
	return ok
}

// get a list of related words for a base (stemmed) word if it exists
func (l *SLexicon) GetStemList(stemmed_word string) []string {
	word_list := make([]string,0)
	if map1, ok := l.stemSet[stemmed_word]; ok {
		for key, _ := range map1 {
			word_list = append(word_list, key)
		}
	}
	return word_list
}

// is the tag a noun tag?
func (l SLexicon)IsNoun(tag string) bool {
	return tag == "NN" || tag == "NNS" || tag == "NNP" || tag == "NNPS"
}

// is the tag a verb tag?
func (l SLexicon)IsVerb(tag string) bool {
	return strings.Contains(tag, "VB")
}

// is the tag an adjective
func (l SLexicon)IsAdj(tag string) bool {
	return tag == "JJ" || tag == "JJR" || tag == "JJS"
}

// is the tag an adverb
func (l SLexicon)IsAdv(tag string) bool {
	return tag == "RB" || tag == "RBR" || tag == "RBS"
}

// is the tag a number
func (l SLexicon)IsNumber(tag string) bool {
	return tag == "CD"
}

// get a list of synonyms
func (l *SLexicon) GetSynonymList(stemmed_word string) []string {
	word_list := make([]string,0)
	if map1, ok := l.synonymSet[stemmed_word]; ok {
		for key, _ := range map1 {
			word_list = append(word_list, key)
		}
	}
	return word_list
}

// init the lexicon from file
func init() {
	err := Lexi.initFromFile()
	if err != nil {
		panic(err)
	}
}

