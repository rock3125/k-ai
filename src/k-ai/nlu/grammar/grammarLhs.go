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

import "strings"

type GrammarLhs struct {
	Name string
    IsPublic bool
    ConversionPattern string	// pattern for converting this to a system entity (e.g. date/time)
    Modifier string             // correct badly parsed entities
    RhsList []GrammarRhs
}

// setup the items that could be nil
func (g *GrammarLhs) Init() {
	g.RhsList = make([]GrammarRhs, 0)
}

func (g GrammarLhs) rhsToString() string {
	str := ""
	for _, rhs := range g.RhsList {
		str += rhs.ToString() + " "
	}
	return strings.TrimSpace(str)
}

func (g GrammarLhs) ToString() string {
	str := ""
	if len(g.ConversionPattern) > 0 {
		str += "pattern " + g.Name + " = " + g.ConversionPattern
	} else {
		if g.IsPublic {
			str += "public "
		} else {
			str += "private "
		}
		str += g.Name + " = " + g.rhsToString()
	}
	return str
}

// get the tokens that can start this Grammar rule
func (g GrammarLhs) GetStartTokens() []string {
	if len(g.RhsList) > 0 {
		rhs := g.RhsList[0]
		// a reference to another rule?
		if rhs.Reference != nil {
			return rhs.Reference.GetStartTokens()
		}
		if len(rhs.Text) > 0 {
			resultSet := make([]string,0)
			resultSet = append(resultSet, rhs.Text)
			return resultSet
		}

		if len(rhs.PatternSet) > 0 {
			resultSet := make([]string,0)
			for key, _ := range rhs.PatternSet {
				resultSet = append(resultSet, key)
			}
			return resultSet;
		}
	}
	panic("invalid pattern - can't get start tokens() for '" + g.Name + "'")
}

