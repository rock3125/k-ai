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
	"strings"
	"k-ai/util"
)


// load the semantics into its map
func (l *SLexicon) loadSemantics(dataDir string) error {
	l.Semantic = make(map[string]string,0)
	files, err := util.GetFilesInDirectory(dataDir + "/lexicon/semantics/*.txt")
	if err != nil { return err }
	for _, filename := range files {
		file_contents, err := util.LoadTextFile(filename)
		if err != nil { return err }
		for _, line := range strings.Split(file_contents, "\n") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				l.AddSemantic(parts[0], parts[1])
				l.testAndAddCompoundWordWithCache(parts[0]) // add this item to the compound word set
			}
		}
	}
	return nil
}


// add a single new semantic to the lexicon - seldom used
func (l *SLexicon) AddSemantic(word string, semantic string) {
	l.Lock()            // one at a time
	defer l.Unlock()

	l.Semantic[word] = strings.ToLower(semantic)
	l.testAndAddCompoundWord(word) // add this item to the compound word set
}


// return the semantic for a noun if it exists, otherwise empty string
func (l *SLexicon) GetSemantic(word string) string {
	l.Lock()            // one at a time
	defer l.Unlock()

	if val, ok := l.Semantic[word]; ok { // non case sensitive first
		return val
	}
	lwrStr := strings.ToLower(word)
	stemmedWord := l.GetStem(lwrStr)
	if val, ok := l.Semantic[stemmedWord]; ok {
		return val
	}
	return ""
}

