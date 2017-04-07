package freebase

import (
	"k-ai/util"
	"strings"
	"errors"
	"k-ai/nlu/model"
	"k-ai/nlu/lexicon"
)

// a single freebase word / tag match
type FreebaseComponent struct {
	base        string      // the base of the word
	tag         string      // the required word tag
}

// the list defined
type FreebaseComponentList []FreebaseComponent

// a predicate and its matches
type FreebaseMatch struct {
	predicate   string                  // the freebase "predicate" to use
	item_list   FreebaseComponentList   // list of items required for the match (without 5wh)
}

// the list defined
type FreebaseMatchList []FreebaseMatch

// the system
type FreebaseMatchSystem struct {
	freebase_map        map[string]FreebaseMatchList
}

// a freebase search set
type FreebaseSearch struct {
	Predicate       string
	TokenList       []model.Token
}

// load the access verbs for the freebase system
func loadFreebaseAccessVerbs() (map[string]FreebaseMatchList, error) {
	freebase_access_map := make(map[string]FreebaseMatchList, 0)

	lines, err := util.LoadTextFile(util.GetDataPath() + "/freebase/freebase-access.txt")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(lines, "\n") {
		parts := strings.Split(line, ",")
		if len(parts) > 1 {
			predicate := parts[0]
			for _, item := range parts[1:] {
				var w5h string
				list := make(FreebaseComponentList,0)
				words := strings.Split(item, " ")
				w5h = words[0]
				for _, word := range words[1:] {
					word_tag := strings.Split(word, ":")
					if len(word_tag) != 2 {
						return nil, errors.New("invalid line:" + line)
					}
					fc := FreebaseComponent{base: word_tag[0], tag: word_tag[1]}
					list = append(list, fc)
				}
				if len(list) > 0 {
					new_item := FreebaseMatch{predicate: predicate, item_list: list}
					if _, ok := freebase_access_map[w5h]; !ok {
						freebase_access_map[w5h] = make(FreebaseMatchList, 0)
					}
					freebase_access_map[w5h] = append(freebase_access_map[w5h], new_item)
				}
			}
		}
	}
	return freebase_access_map, nil
}

// the system
var MatchSystem = FreebaseMatchSystem{}

// setup/load the data
func (fb* FreebaseMatchSystem)Setup() error {
	var err error
	fb.freebase_map, err = loadFreebaseAccessVerbs()
	if err != nil {
		return err
	}
	return nil
}

// is any of the tokens in the token-list a match for this template item?
func (match FreebaseComponent)IsMatch(token_list []model.Token) bool {
	for _, t_token := range token_list {
		stem := strings.ToLower(lexicon.Lexi.GetStem(t_token.Text))
		if match.base == stem {
			if match.tag == "n" {
				return lexicon.Lexi.IsNoun(t_token.Tag)
			} else if match.tag == "v" {
				return lexicon.Lexi.IsVerb(t_token.Tag)
			}
		}
	}
	return false
}

// is any of the tokens in the token-list a match for this item?
func (match FreebaseComponent)IsMatchToken(t_token model.Token) bool {
	stem := strings.ToLower(lexicon.Lexi.GetStem(t_token.Text))
	if match.base == stem {
		if match.tag == "n" {
			return lexicon.Lexi.IsNoun(t_token.Tag)
		} else if match.tag == "v" {
			return lexicon.Lexi.IsVerb(t_token.Tag)
		}
	}
	return false
}

// create a reduced search token list
func (fb FreebaseMatchSystem)createSearchTokenList(match FreebaseMatch, token_list []model.Token) []model.Token {
	new_list := make([]model.Token,0)
	for _, t_token := range token_list {
		// check the token is "search worthy"
		if lexicon.Lexi.IsNoun(t_token.Tag) || lexicon.Lexi.IsVerb(t_token.Tag) ||
			lexicon.Lexi.IsAdj(t_token.Tag) || lexicon.Lexi.IsAdv(t_token.Tag) ||
			lexicon.Lexi.IsNumber(t_token.Tag) {

			// then check it ISN'T part of the match item itself
			if !match.IsMatchToken(t_token) {
				new_list = append(new_list, t_token)
			}

		}
	}
	return new_list
}

// do we have match for the token set
func (match FreebaseMatch)IsMatch(token_list []model.Token) bool {
	for _, item := range match.item_list {
		// does the item exist in the token_list?
		if item.IsMatch(token_list) {
			return true
		} else {
			return false
		}
	}
	return false
}

// do we have match for the token set
func (match FreebaseMatch)IsMatchToken(t_token model.Token) bool {
	for _, item := range match.item_list {
		// does the item exist in the token_list?
		if item.IsMatchToken(t_token) {
			return true
		}
	}
	return false
}

// make a match and return the valid words
func (fb FreebaseMatchSystem)Match(token_list []model.Token) (*FreebaseSearch, error) {
	if len(token_list) > 0 {
		// word0 is the w5h accessor, and must exist for now all queries are w5h
		word0 := strings.ToLower(token_list[0].Text)
		if item_list, exists := fb.freebase_map[word0]; exists {
			for _, item := range item_list { // these are all items for this w5h
				if item.IsMatch(token_list[1:]) {  // do the other tokens match?
					token_list := fb.createSearchTokenList(item, token_list[1:])
					if len(token_list) > 0 {
						return &FreebaseSearch{Predicate: item.predicate, TokenList: token_list}, nil
					}
				}
			}
		}
	}
	return nil, errors.New("not found")
}

