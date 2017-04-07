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
	"os"
	"fmt"
	"sync"
)

const log_file = "/var/log/kai/system.log"

type appLogger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger

	sync.Mutex
}

// create a new file logger
func FileLogger() Logger {
	logger := &appLogger{}

	f_handle, err := os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Error opening log-file(%s): %v", log_file, err))
	}
	logger.debug = log.New(f_handle, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.info = log.New(f_handle, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.warning = log.New(f_handle, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.error = log.New(f_handle, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)

	return logger
}

func (logger *appLogger) String() string {
	return log_file
}

func (logger *appLogger) Debug(format string, args ...interface{}) {
	if logger != nil && logger.debug != nil {
		logger.Lock()
		defer logger.Unlock()
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		logger.debug.Printf(format, args...)
	}
}

func (logger *appLogger) Info(format string, args ...interface{}) {
	if logger != nil && logger.info != nil {
		logger.Lock()
		defer logger.Unlock()
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		logger.info.Printf(format, args...)
	}
}

func (logger *appLogger) Warning(format string, args ...interface{}) {
	if logger != nil && logger.warning != nil {
		logger.Lock()
		defer logger.Unlock()
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		logger.warning.Printf(format, args...)
	}
}

func (logger *appLogger) Error(format string, args ...interface{}) {
	if logger != nil && logger.error != nil {
		logger.Lock()
		defer logger.Unlock()
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		logger.error.Printf(format, args...)
	}
}


// setup the system main logger
func init() {
	Log = FileLogger()
}

