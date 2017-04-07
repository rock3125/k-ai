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

package super_search

import (
	"k-ai/db/db_model"
	"k-ai/nlu/model"
)

// a parse tree structure for a super search
type SSTree struct {

	db_model.Index                              // what we're looking for

	TType string                                // the tree type of this item, one of {or,and,and not,word}
	Offset int                                  // the offset

	Semantic string                             // the semantic to look for
	Exact bool                                  // exact match required?

	Left*   SSTree                              // left and right children
	Right*  SSTree

	IndexMap map[string][]model.IndexMatch      // the result items (only applicable for ttype "word")
}



type SSTokenWithIndex struct {

	model.Token                     // a token
	Index int                       // the index of the token

}

