package freebase

import (
	"testing"
	"k-ai/nlu/model"
	"strings"
)

// create a fake token list from a space separated string
func helper_create_tokenList(text string) []model.Token {
	tokenList := make([]model.Token,0)
	parts := strings.Split(text, " ")
	counter := 0
	for _, part := range parts {
		tokenList = append(tokenList, model.Token{Text: part, Tag: "NN", Index: counter})
		counter += 1
	}
	return tokenList
}

// test simple patterns first
func TestFBPatterns_1(t *testing.T) {
	// add a simple pattern
	nodes := make(map[string]*FBPatternTree, 0)
	err := AddPatternToTree("(who|whose) are %1", "%1 is_a %2", nodes)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	// match this pattern against all possible
	vars := make(map[string][]model.Token, 0)
	exec := MatchTokenList(helper_create_tokenList("Whose are these things?"), nodes, vars)
	if len(exec) == 0 || len(vars) == 0 {
		t.Error("match failed")
		t.FailNow()
	}

}

