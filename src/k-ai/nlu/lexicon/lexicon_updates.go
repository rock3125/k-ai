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

// apply all updates from the lexicon update file to the semantics of the lexicon system
func (l *SLexicon) applySemanticUpdates() error {

	file_contents, err := util.LoadTextFile(util.GetDataPath() + "/lexicon/lexicon_updates.txt")
	if err != nil { return err }

	for _, line := range strings.Split(file_contents, "\n") {
		parts := strings.Split(line, "|")
		if len(parts) == 3 {
			word_sem := strings.Split(parts[2], ":")
			if len(word_sem) == 2 {
				switch (parts[0]) {
				case "del":
					delete(l.Semantic, word_sem[0])
				case "save":
					l.Semantic[word_sem[0]] = word_sem[1]
				default:
					panic("unknown instruction in line " + line)
				}
			}
		}
	}
	return nil
}

