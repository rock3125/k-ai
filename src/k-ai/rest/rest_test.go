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

package rest

import (
	"testing"
	"k-ai/util"
	"k-ai/util_ut"
)

// test a few of the little utilities in here

func TestValidPath(t *testing.T) {
	// check GetDataPath() doesn't have any .. in it since it uses Abs()
	path := util.GetDataPath() + "/../../"
	util_ut.IsTrue(t, !isValidPath(path))

	path_2 := util.GetDataPath() + "/web/"
	util_ut.IsTrue(t, isValidPath(path_2))
}

// test suffix check routine
func TestValidSuffix(t *testing.T) {
	// check GetDataPath() doesn't have any .. in it since it uses Abs()
	path := util.GetDataPath() + "/test/test.jpg"
	util_ut.IsTrue(t, !isValidSuffix(path,".png", ".jpeg"))
	util_ut.IsTrue(t, isValidSuffix(path,".png", ".jpeg", ".jpg"))
	util_ut.IsTrue(t, isValidSuffix(path,".jpg"))
}


