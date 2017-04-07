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
	"log"
	"bytes"
	"sync"
)

type memoryLogger struct {
	Buffer  *bytes.Buffer

	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger

	sync.Mutex
}


// create a memory logger
func MemoryLogger() Logger {
	logger := &memoryLogger{}

	logger.Buffer = &bytes.Buffer{}

	logger.debug = log.New(logger.Buffer, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.info = log.New(logger.Buffer, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.warning = log.New(logger.Buffer, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.error = log.New(logger.Buffer, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)

	return logger
}

func (logger *memoryLogger) String() string {
	return logger.Buffer.String()
}

func (logger *memoryLogger) Debug(format string, args ...interface{}) {
	if logger != nil && logger.debug != nil {
		logger.Lock()
		defer logger.Unlock()
		logger.debug.Printf(format, args...)
	}
}

func (logger *memoryLogger) Info(format string, args ...interface{}) {
	if logger != nil && logger.info != nil {
		logger.Lock()
		defer logger.Unlock()
		logger.info.Printf(format, args...)
	}
}

func (logger *memoryLogger) Warning(format string, args ...interface{}) {
	if logger != nil && logger.warning != nil {
		logger.Lock()
		defer logger.Unlock()
		logger.warning.Printf(format, args...)
	}
}

func (logger *memoryLogger) Error(format string, args ...interface{}) {
	if logger != nil && logger.error != nil {
		logger.Lock()
		defer logger.Unlock()
		logger.error.Printf(format, args...)
	}
}

