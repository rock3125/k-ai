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

package util_ut

import (
	"testing"
	"runtime/debug"
	"log"
)

// test fn, unit testing aid
func IsTrue(t *testing.T, cond bool) {
	if !cond {
		debug.PrintStack()
		t.FailNow()
	}
}


// check there is an error - and print it and stop if there is
func Check(t *testing.T, e error) {
	if e != nil {
		log.Print(e)
		debug.PrintStack()
		t.FailNow()
	}
}

