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

package glove

// a single entry - the text, the vector, and its score relative to some query
type MatchVector struct {
	Text string             // glove text
	Score float64           // scoring for searching
	Vector []float64        // glove vector
}

// for sorting the vector - a data-type defn.
type MatchVectors []MatchVector

// len for sort interface
func (slice MatchVectors) Len() int {
    return len(slice)
}

// less for sorting, sort by highest score, then ABC text
func (slice MatchVectors) Less(i, j int) bool {
	if slice[i].Score != slice[j].Score {
		return slice[i].Score > slice[j].Score;
	} else {
		return slice[i].Text < slice[j].Text
	}
}

// sort interface
func (slice MatchVectors) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

