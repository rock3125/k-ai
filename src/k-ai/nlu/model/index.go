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

import "github.com/gocql/gocql"

// index item after a match result
type IndexMatch struct {
	// NB. this is a copy of model.Index!!!  keep in sync
	Sentence_id gocql.UUID	// the sentence owner id

	Word        string 		// the word, main index
	Tag         string 		// the Penn tag of the word
	Shard       int    		// shard spreading across systems

	Offset      int    		// unique offset for repeating words
	Topic		string 		// what spawned it, what is the sentence owner?
	Score		float64		// the value of this index relative to others

	// separate additional
	KeywordIndex int    // the keyword's index / match
}

// convert an index to an index match
func Convert(sentence_id gocql.UUID, word string, tag string, shard int, offset int,
			 topic string, score float64, keywordIndex int) *IndexMatch {

	return &IndexMatch{
			KeywordIndex: keywordIndex,

			Sentence_id: sentence_id,

			Word:        word,
			Tag:         tag,
			Shard:       shard,

			Offset:      offset,
			Topic:       topic,
			Score:		 score,
		}
}

