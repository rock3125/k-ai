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

import "strings"

var pennStrings = [...]string {
	"cc", "cd", "dt", "ex", "fw", "in", "jj", "jjr", "jjs", "ls", "md", "nn", "nns", "nnp", "nnps", "pdt", "pos", "prp", "prp$",
	"prps", "rb", "rbr", "rbs", "rp", "sym", "to", "intj", "uh", "vb", "vbd", "vbg", "vbn", "vbp", "vbz", "wdt", "wp", "wp$", "wps",
	"wrb", "rsb", "-rsb-", "rcb", "-rcb-", "rrb", "-rrb-", "-lsb-", "lsb", "-lcb-", "lcb", "lrb", "-lrb-", "np-tmp",
	"pun", "hyph", ".", "sqt", "eqt", "x", "xx", "sp", " ",
	"adjp", "advp", "conjp", "np", "vp", "pp", "qp", "s", "sq", "sbarq", "sbar", "sinv", "ucp", "whadjp", "whadvp", "whnp", "whpp", "root", "prn",
	"frag", "prt", "rrc", "nx", "nac", "lst", "add", "afx", "gw", "bes", "hvs", "nfp",
}

var pennStringMap = make(map[string]bool,0)

// return true if str (case insensitive) is one of the known Penn-types
func IsPennType(str string) bool {
	if len(pennStringMap) == 0 {
		for _, item := range pennStrings {
			pennStringMap[item] = true
		}
	}
	_, ok := pennStringMap[strings.ToLower(str)]
	return ok
}

