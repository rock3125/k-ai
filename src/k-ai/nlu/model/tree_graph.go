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
	"time"
	"os"
	"strconv"
	"io/ioutil"
	"os/exec"
	"strings"
	"errors"
)

// helper for ToGraphVizDot below - creates the nodes and the strings for the GraphViz tree
func (t *Tree) toGraphVizDotHelper() (string) {
	indent := "    "
	if t.Tokens.Len() == 1 && strings.HasPrefix(t.Tokens.TokenList[0].Tag, "VB") {
		str := ""
		tt := t.Tokens.TokenList[0]
		label := t.Tokens.ToString() + ", " + tt.Tag
		node := "node" + strconv.Itoa(tt.Index)
		str += indent + node + " [shape=box,label=\"" + label + "\"];\n"
		if t.Left != nil {
			ttl := t.Left.Tokens.TokenList[0]
			nodel := "node" + strconv.Itoa(ttl.Index)
			str += indent + node + " -> " + nodel + ";\n"
			str += t.Left.toGraphVizDotHelper()
		}
		if t.Right != nil {
			ttr := t.Right.Tokens.TokenList[0]
			noder := "node" + strconv.Itoa(ttr.Index)
			str += indent + node + " -> " + noder + ";\n"
			str += t.Right.toGraphVizDotHelper()
		}
		return str
	} else {
		str := ""
		if len(t.Tokens.TokenList) > 0 {
			tt := t.Tokens.TokenList[0]
			label := t.Tokens.ToString()

			// add some extras to the label
			if tt.Text != tt.Tag {
				label += ", " + tt.Tag
			}
			if len(tt.Anaphora) > 0 {
				label += ", ref:" + tt.Anaphora
			}
			if len(tt.Semantic) > 0 {
				label += ", sem:" + tt.Semantic
			}

			node := "node" + strconv.Itoa(tt.Index)
			str += indent + node + " [label=\"" + label + "\"];\n"

			if t.Left != nil {
				ttl := t.Left.Tokens.TokenList[0]
				nodel := "node" + strconv.Itoa(ttl.Index)
				str += indent + node + " -> " + nodel + ";\n"
				str += t.Left.toGraphVizDotHelper()
			}

			if t.Right != nil {
				ttr := t.Right.Tokens.TokenList[0]
				noder := "node" + strconv.Itoa(ttr.Index)
				str += indent + node + " -> " + noder + ";\n"
				str += t.Right.toGraphVizDotHelper()
			}
		}
		return str
	}
}

// convert the tree to a drawable graphviz Dot file
// see /usr/bin/dot
func (t *Tree) ToGraphVizDot() (string) {
	str := "digraph G {\n"
	str += t.toGraphVizDotHelper()
	str += "}\n"
	return str
}

// convert the tree to a drawable graphviz Dot file
// see /usr/bin/dot
func (t TreeList) ToGraphVizDot() (string) {
	str := "digraph G {\n"
	for _, ttree := range t {
		str += ttree.toGraphVizDotHelper()
	}
	str += "}\n"
	return str
}

// get unique number fo this process
func uid() int {
	value := int(time.Now().UnixNano() + int64(os.Getpid()))
	if value < 0 { value = -value }
	return value
}

// convert a dot text description to a PNG file and return the PNG bytes
// this requires to DOT command be installed on the system, part of graphviz
func ToPng(dotContent string) ([]byte, error) {
	if len(dotContent) > 0 {
		// generate the dot and write it to a temp file
		dot_temp := "/tmp/bt-dot-" + strconv.Itoa(uid()) + ".dot"
		ioutil.WriteFile(dot_temp, []byte(dotContent), 0644)

		png_temp := "/tmp/bt-png-" + strconv.Itoa(uid()) + ".png"
		bc := make(chan string, 1)
		errs := make(chan error, 1)
		go func() {
			body, err := exec.Command("dot", "-Tpng", dot_temp, "-o", png_temp).Output()
			if err != nil {
				errs <- err
			} else {
				bc <- string(body)
			}
			close(bc)
			close(errs)
		}()

		// check for errors
		err, _ := <-errs
		if err != nil {
			return nil, err
		}

		<-bc // just read the stdout from the command
		b, err := ioutil.ReadFile(png_temp)

		// remove temp files
		os.Remove(dot_temp)
		os.Remove(png_temp)

		return b, err
	}
	return nil, errors.New("invalid parameter")
}
