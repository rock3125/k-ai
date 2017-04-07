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

package util

import (
	"testing"
	"github.com/gocql/gocql"
	"strings"
	"k-ai/util_ut"
)

// test a simple rule with name substitution (from ai.aiml, line 82)
func TestUtilUUID(t *testing.T) {
	test_uuid := gocql.UUID{}
	util_ut.IsTrue(t, IsEmpty(&test_uuid)) // test empty one

	// setup
	for i, _ := range test_uuid {
		test_uuid[i] = byte(i)
	}
	util_ut.IsTrue(t, !IsEmpty(&test_uuid)) // test empty one

	test_uuid2 := gocql.UUID{}
	util_ut.IsTrue(t, IsEmpty(&test_uuid2)) // test empty one

	CopyUUID(&test_uuid2, &test_uuid)
	util_ut.IsTrue(t, !IsEmpty(&test_uuid2)) // test empty one

	for i, _ := range test_uuid {
		util_ut.IsTrue(t, test_uuid[i] == test_uuid2[i])
	}
}

func TestPathFns(t *testing.T) {
	util_ut.IsTrue(t, len(GetDataPath()) > 0)
	a_psycho, err := LoadTextFile(GetDataPath() + "/test_data/american-psycho.txt")
	util_ut.Check(t, err)
	util_ut.IsTrue(t, len(a_psycho) > 0)
	// check GetDataPath() doesn't have any .. in it since it uses Abs()
	util_ut.IsTrue(t, !strings.Contains(GetDataPath(),".."))
}

func TestValidEmail(t *testing.T) {
	util_ut.IsTrue(t, ValidateEmail("peter@peter.co.nz"))
	util_ut.IsTrue(t, !ValidateEmail("peterpeter.co.nz"))
	util_ut.IsTrue(t, !ValidateEmail("peter@peter"))
	util_ut.IsTrue(t, !ValidateEmail("peter@peter.co.nzasdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"))
}

