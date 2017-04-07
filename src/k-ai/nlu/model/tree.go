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
)

type Tree struct {
	Tokens  Sentence        // the tokens of this part of the tree
    Srl     string          // srl label
    Offset  int             // start offset of item in the SRL tree
    Left    *Tree
	Right   *Tree
}

type TreeList []Tree


// convert a parsed sentence structure to a Tree, and return the root of the tree
func SentenceToTuple(sentence Sentence) (*Tree) {
	root := &Tree{}
	if sentence.Len() > 0 {
		// create a lookup for Index(int) -> Tree item
		lookup := make(map[int]*Tree)
		for _, token := range sentence.TokenList {
			s := Sentence{}
			s.TokenList = append(s.TokenList, token)
			lookup[token.Index] = &Tree{Srl: token.Dep, Tokens: s, Offset: token.Index}
		}
		// hook the items into the tree using the lookup map for parentage using the AncestorList
		for _, token := range sentence.TokenList {
			if len(token.AncestorList) == 0 {  // root has no parents / ancestors
				root = lookup[token.Index]
			} else {
				current := lookup[token.Index]  // get a tree item
				if current != nil {
					anc := token.AncestorList[0]  // find the correct root - can't be itself
					for i := 0; anc == token.Index && i < len(token.AncestorList); i++ {
						anc = token.AncestorList[i]
					}
					if anc != token.Index {  // can never have itself as a parent
						parent := lookup[anc]
						if parent != nil {
							parent.AddNode(current)  // setup parent child relationship
						}
					}
				}
			}
		}
	}
	return root
}

// simple part string like
// [Craig{person}:nsubj] [has{person}:ROOT] [a:det] [boat{vehicle}:dobj] [in:prep] [the:det] [harbour{location}:pobj] [.:punct]
func (t *Tree) ToString() (string) {
	str := ""
	if t.Left != nil {
		str += t.Left.ToString()
	}
	str = str + " [" + t.Tokens.ToString() + ":" + t.Srl + "]"
	if t.Right != nil {
		str += t.Right.ToString()
	}
	return str
}

// more complex string with verbs as main items
// _has{VBZ}_ (Craig{person} | a boat{vehicle} in the harbour{location} .)
func (t *Tree) ToStringIndent() (string) {
	if t.Tokens.Len() == 1 && strings.HasPrefix(t.Tokens.TokenList[0].Tag, "VB") {
		str := ""
		str += "_" + t.Tokens.TokenList[0].Text + "{" + t.Tokens.TokenList[0].Tag + "}_ ("
		if t.Left != nil {
			str += strings.TrimSpace(t.Left.ToStringIndent()) + " | "
		}
		if t.Right != nil {
			str += strings.TrimSpace(t.Right.ToStringIndent())
		}
		str += ")"
		return str
	} else {
		str := ""
		if t.Left != nil {
			str += t.Left.ToStringIndent()
		}
		str += t.Tokens.ToString() + " "
		if t.Right != nil {
			str += t.Right.ToStringIndent()
		}
		return str
	}
}

// recursively add an item into the tuple tree
func (tt *Tree) add(srl string, offset int, tokens Sentence) {
    if len(srl) > 0 && tokens.Len() > 0 {
        if (offset < tt.Offset) {
            if tt.Left == nil {
                tt.Left = &Tree{Srl: srl, Offset: offset, Tokens: tokens}
            } else {
                tt.Left.add(srl, offset, tokens);
            }
        } else if (offset > tt.Offset) {
            if tt.Right == nil {
                tt.Right = &Tree{Srl: srl, Offset: offset, Tokens: tokens}
            } else {
                tt.Right.add(srl, offset, tokens);
            }
        }
    }
}

// recursively add an item into the tuple tree
func (tt *Tree) AddNode(child *Tree) {
    if child != nil {
        if child.Offset < tt.Offset {
            if tt.Left == nil {
                tt.Left = child
            } else {
                tt.Left.AddNode(child)
            }
        } else if child.Offset > tt.Offset {
            if tt.Right == nil {
                tt.Right = child
            } else {
                tt.Right.AddNode(child)
            }
        }
    }
}

