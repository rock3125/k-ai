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

package logger

import (
	"testing"
	"strings"
	"k-ai/util_ut"
)

// test the logger works by logging into a string buffer
func TestLogging(t *testing.T) {

	logger := MemoryLogger()

	logger.Debug("debug to stdout")
	logger.Info("info to stdout")
	logger.Warning("warning to stdout")
	logger.Error("error to stdout")

	log_str := logger.String()
	util_ut.IsTrue(t, strings.Contains(log_str, "debug to stdout"))
	util_ut.IsTrue(t, strings.Contains(log_str, "info to stdout"))
	util_ut.IsTrue(t, strings.Contains(log_str, "warning to stdout"))
	util_ut.IsTrue(t, strings.Contains(log_str, "error to stdout"))
}

