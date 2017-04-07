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

package grammar

import (
	"fmt"
)

type GrammarRhs struct {
    Text string					// the text / literal
    Reference *GrammarLhs		// a reference to another rule
    IsRepeat bool				// + at the end of strings length > 1
    PatternSet map[string]bool	// bag of words / or equivalent
    NumberRangeStart int		// number.range(1,31) type of filtering
    NumberRangeEnd int
}

// setup the items that could be nil
func (g *GrammarRhs) Init() {
	g.PatternSet = make(map[string]bool, 0)
}

func (g GrammarRhs) ToString() string {
	if len(g.Text) > 0 {
		if g.IsRepeat {
			return g.Text + "+";
		}
		if g.NumberRangeStart != 0 || g.NumberRangeEnd != 0 {
			return fmt.Sprintf("%s.range(%d,%d)", g.Text, g.NumberRangeStart, g.NumberRangeEnd)
		}
		return g.Text
	}
	if g.Reference != nil {
		return "<" + g.Reference.rhsToString() + ">"
	}
	if len(g.PatternSet) > 0 {
		str := "[ ";
		for key, _ := range g.PatternSet {
			str += key + " "
		}
		str += "]"
		return str
	}
	return "<null>"
}


