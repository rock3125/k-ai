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

package anaphora

import (
	"strings"
	"math"
	"sort"
	"k-ai/nlu/model"
)

//
// Anaphora Resolution
//
// see: https://www.cl.cam.ac.uk/teaching/1011/L104/lec12-2x2.pdf
// Two different operations are performed:
//
// Maintaining and updating a discourse model consisting of a set of co-reference classes:
// * Each co-reference class corresponds to one entity that has been evoked in the discourse
// * Each co-reference class has an updated salience value
//
// Resolving each Pronoun from left to right
// * Collect potential referents from up to 4 sentences back
// * Filter out coreference classes that don’t satisfy  agreement/syntax constraints
// * Select remaining co-reference class with the highest salience value; add pronoun to class.
//
// The salience of a referent is calculated on the basis of recency and grammatical function.
// Salience Factor        Example                          Weight
// Current sentence                                          100
// Subject emphasis       John opened the door                80
// Existential emphasis   There was a dog standing outside    70
// Accusative emphasis    John liked the dog                  50
// Indirect object        John gave a biscuit to the dog      40
// Adverbial emphasis     Inside the house, the cat looked on 50
// Head Noun emphasis     The cat in the house looked on      80

// The salience of a referent is the sum of all applicable weights
// The salience of a referent is halved each time a sentence  boundary is crossed
// This, along with the weight for being in the current sentence,  makes more recent referents more salient
// Weights are calculated for each member of the salience class
// Previous mentions can boost the salience of a coreference class
// This accounts for the repetition effect
// Lappin and Leass report 86% accuracy for their algorithm on a corpus of Computer manuals


type LappinLeass struct {
	pronoun_set map[string]LLPronoun	// list of pronouns to look for in the text (ones we can resolve)
	n_back int							// go back up to n-sentences in the list for resolution
}

type LLReferent struct {
	anaphora string			// the anaphora text resolved to
	salience float64		// the salient score of the referent
}

type LLReferentList []*LLReferent

// sort interface for LLReferent
func (llr LLReferentList) Len() int {
	return len(llr)
}
// sort by highest salience
func (llr LLReferentList) Less(i, j int) bool {
	return llr[i].salience > llr[j].salience
}
func (llr LLReferentList) Swap(i, j int) {
	llr[i], llr[j] = llr[j], llr[i]
}

type LLPronoun struct {
	text string					// the pronoun's text (e.g. "he")
	semantics []string			// specific semantics, one of {"man', "woman", "person", or "other"}, where other will match anything
	number string				// the pronoun's number, one of {"s", "p"}  (singular/plural)
}


// return true if this pronoun contains a semantic matching the parameter
func (pr LLPronoun) containsSemantic(semantic string) bool {
	for _, sem := range pr.semantics {
		if sem == semantic {
			return true
		}
	}
	return false
}


func (ll *LappinLeass) Init() {
	ll.n_back = 4 // number of sentences to scan back for pronouns, default: 4

	// setup a lookup map for valid 3rd person pronouns
	ll.pronoun_set = make(map[string]LLPronoun,0)

	ll.pronoun_set["he"] = LLPronoun{text: "he", number: "s", semantics: []string{"male","person"}}
	ll.pronoun_set["she"] = LLPronoun{text: "she", number: "s", semantics: []string{"female","person"}}
	ll.pronoun_set["it"] = LLPronoun{text: "it", number: "s", semantics: []string{"other"}}

	ll.pronoun_set["him"] = LLPronoun{text: "him", number: "s", semantics: []string{"male","person"}}
	ll.pronoun_set["her"] = LLPronoun{text: "her", number: "s", semantics: []string{"female","person"}}

	ll.pronoun_set["himself"] = LLPronoun{text: "himself", number: "s", semantics: []string{"male","person"}}
	ll.pronoun_set["herself"] = LLPronoun{text: "herself", number: "s", semantics: []string{"female","person"}}
	ll.pronoun_set["itself"] = LLPronoun{text: "itself", number: "s", semantics: []string{"other"}}

	ll.pronoun_set["his"] = LLPronoun{text: "his", number: "s", semantics: []string{"male","person"}}
	ll.pronoun_set["her"] = LLPronoun{text: "her", number: "s", semantics: []string{"female","person"}}
	ll.pronoun_set["hers"] = LLPronoun{text: "hers", number: "s", semantics: []string{"female","person"}}
	ll.pronoun_set["its"] = LLPronoun{text: "its", number: "s", semantics: []string{"other"}}

	ll.pronoun_set["they"] = LLPronoun{text: "they", number: "p", semantics: []string{"male", "female", "person", "other"}}
	ll.pronoun_set["them"] = LLPronoun{text: "them", number: "p", semantics: []string{"male", "female", "person", "other"}}
	ll.pronoun_set["themselves"] = LLPronoun{text: "themselves", number: "p", semantics: []string{"male", "female", "person"}}
}

// is the current "nsubj" Noun succeeded by another np? (a DET or IN, no verb)
func (ll LappinLeass) hasHeadNounEmphasis(index int, sentence model.Sentence) bool {
	// before?
	for i := index + 1; i < len(sentence.TokenList); i++ {
		t_token := sentence.TokenList[i]
		if t_token.Tag == "IN" || t_token.Tag == "DET" {
			return true
		}
		if strings.HasPrefix(t_token.Tag, "VB") {
			break
		}
	}
	return false
}

// is the current "nsubj" Noun preceeded by another np? (a DET or IN, no verb)
func (ll LappinLeass) hasAdverbialEmphasis(index int, sentence model.Sentence) bool {
	// before?
	found_det_or_in := false
	for i := index - 1; i >= 0; i-- {
		t_token := sentence.TokenList[i]
		if t_token.Tag == "IN" || t_token.Tag == "DET" { // mark the start of a new np
			found_det_or_in = true
		}
		if strings.HasPrefix(t_token.Tag, "NN") && found_det_or_in { // found another np?
			return true
		}
		if strings.HasPrefix(t_token.Tag, "VB") { // verbs are bad
			break
		}
	}
	return false
}

// calculate the salience value using grammatical constructs for a noun
// is_last_sentence:  true if this sentence is the one with the pronoun
// seen_existential:  true if this sentence thusfar has seen an EX tag
// token: the noun token under investigation
// index: its index into sentence
// sentence: the sentence of this token
func (ll LappinLeass) calculateSalience(seen_existential bool,
										token *model.Token, index int, sentence model.Sentence) float64 {
	// calculate salience
	salience := 100.0  // basic score
	// have we seen an existential marker?
	if seen_existential {
		salience += 70.0
	}
	// subject?
	if token.Dep == "nsubj" {
		salience += 80
	}
	// accusative emphasis
	if token.Dep == "dobj" {
		salience += 50
	}
	// indirect object
	if token.Dep == "pobj" {
		salience += 40
	}
	// Adverbial emphasis (head verb (nsubj) is preceeded by another np)
	if token.Dep == "nsubj" && ll.hasAdverbialEmphasis(index, sentence) {
		salience += 50
	}
	// head noun emphases (head verb (nsubj is followed by another np)
	if token.Dep == "nsubj" && ll.hasHeadNounEmphasis(index, sentence) {
		salience += 80
	}
	return salience
}


// is this token semantically compatible with the pronoun?
func (ll LappinLeass) isSemanticMatch(token *model.Token, pronoun *LLPronoun) bool {
	if pronoun.containsSemantic(token.Semantic) { // otherwise - it must be one of its semantics
		return true
	}
	if token.Semantic == "male" || token.Semantic == "female" || token.Semantic == "person" {
		return false
	}
	if pronoun.containsSemantic("other") {
		return true
	}
	return false
}


// is this token number compatible with the pronoun
func (ll LappinLeass) matchesNumber(token *model.Token, pronoun *LLPronoun) bool {
	if pronoun.number == "p" {  // plural
		return token.Tag == "NNS" || token.Tag == "NNPS"
	}
	return true
}


// find suitable references to resolve pronoun
// pronoun: the pronoun that needs resolving
// s_index: the sentence index for the sentence to process
// t_index: the pronoun's offset into sentence_list[s_index]
// sentence_list: a window of sentences to use for pronoun resolution if possible
// returns a sorted list of LLReferent with the most likely referent at position [0] (highest score)
func (ll LappinLeass) find_pronouns(pronoun *LLPronoun, s_index int, t_index int, sentence_list model.SentenceList) LLReferentList {
	// collect noun phrases that might match
	referent_array := make(LLReferentList, 0)
	num_sentences_back := s_index - ll.n_back
	if num_sentences_back < 0 { num_sentences_back = 0 }

	for sentence_id := num_sentences_back; sentence_id <= s_index; sentence_id++ {
		salient_dropoff := math.Pow(2.0, float64(sentence_id - s_index))
		sentence := sentence_list[sentence_id]
		is_last := sentence_id == s_index

		seen_existential := false
		for i, t_token := range sentence.TokenList {
			if (is_last && i < t_index) || !is_last { // restrict to anything before the part in the last sentence
				if t_token.Tag == "EX" {
					seen_existential = true
				}
				if strings.HasPrefix(t_token.Tag, "NN") { // noun
					// is this of the right semantic for the pronoun?
					if ll.isSemanticMatch(&t_token, pronoun) && ll.matchesNumber(&t_token, pronoun) {
						// calculate its salience
						salience := ll.calculateSalience(seen_existential, &t_token, i, sentence)
						// multiply with drop-off for farther away sentences
						salience *= salient_dropoff
						// add new referent
						referent_array = append(referent_array, &LLReferent{anaphora: t_token.Text, salience: salience})

					} // if is semantic match

				} // if is noun

			} // if right part of the sentence

		} // for each token in the sentence

	} // for each back sentence

	// sort salience array
	sort.Sort(referent_array)
	return referent_array
}


// resolve pronouns in a sentence list - it is assumed that the last sentence
//   sentence_list: the last sentence is assumed to have a pronoun reference
// return the number of pronouns that did get resolved in this sentence
func (ll LappinLeass) ResolvePronouns(sentence_list model.SentenceList) int {
	num_pronouns_resolved := 0

	for s_index, _ := range sentence_list {
		// find the pronoun(s) to be resolved, from left to right
		for index, t_token := range sentence_list[s_index].TokenList {
			if t_token.Tag == "PRP" || t_token.Tag == "PRP$" {
				if prp, ok := ll.pronoun_set[strings.ToLower(t_token.Text)]; ok {
					referent_list := ll.find_pronouns(&prp, s_index, index, sentence_list)
					if len(referent_list) > 0 {
						num_pronouns_resolved += 1
						sentence_list[s_index].TokenList[index].Anaphora = referent_list[0].anaphora
					} else {
						sentence_list[s_index].TokenList[index].Anaphora = "?" // not found marker
					}
				}
			}
		}
	}
	return num_pronouns_resolved
}


// does a sentence have a pronoun in it we can try and resolve?
func (ll LappinLeass) HasPronoun(sentence model.Sentence) bool {
	for _, t_token := range sentence.TokenList {
		if t_token.Tag == "PRP" || t_token.Tag == "PRP$" {
			if _, ok := ll.pronoun_set[strings.ToLower(t_token.Text)]; ok {
				return true
			}
		}
	}
	return false
}

// share access to LL algorithm
var LL LappinLeass


// initializer
func init() {
	LL.Init()
}

