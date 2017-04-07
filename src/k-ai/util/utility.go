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
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/gocql/gocql"
	"regexp"
	"time"
	"bytes"
	"strconv"
)

// exists returns whether the given file or directory exists or not
func Exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

var cached_data_path string = ""

// get path to data directory
func GetDataPath() (string) {
	if len(cached_data_path) == 0 {
		// get cwd
		myPath, _ := os.Getwd()
		dataDir := myPath + "/data"
		for i := 0; i < 5 && !Exists(dataDir); i++ { // pick altenative director if not found
			myPath += "/.."
			dataDir = myPath + "/data"
		}
		if !Exists(dataDir) {
			panic("data directory not found @ " + dataDir)
		}
		var err error
		cached_data_path, err = filepath.Abs(dataDir)
		if err != nil {
			panic(err)
		}
	}
	return cached_data_path
}


// read a text file
func LoadTextFile(filename string) (string, error) {
	// open the file and check it
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// write a text file
func SaveTextFile(filename string, str string) error {
	return ioutil.WriteFile(filename, []byte(str), 0644)
}

// return all files in a directory by glob
func GetFilesInDirectory(directoryGlob string) ([]string, error) {
	file_list, err := filepath.Glob(directoryGlob)
	return file_list, err
}

// UUID empty?
func IsEmpty(uuid *gocql.UUID) bool {
	if uuid == nil {
		return true
	}
	for _, value := range uuid {
		if value != 0 {
			return false
		}
	}
	return true
}

// UUID empty?  uuid1 = uuid2 content
func CopyUUID(uuid1 *gocql.UUID, uuid2 *gocql.UUID) {
	if uuid1 != nil && uuid2 != nil {
		for i, value := range uuid2 {
			uuid1[i] = value
		}
	}
}

// return true if the two are equal
func EqualUUID(uuid1 *gocql.UUID, uuid2 *gocql.UUID) bool {
	if uuid1 != nil && uuid2 != nil {
		for i, value := range uuid2 {
			if uuid1[i] != value {
				return false
			}
		}
	}
	return true
}

// validate an email address, return true if the email address is valid
func ValidateEmail(email string) bool {
	var isEmail = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	return len(email) < 50 && isEmail.MatchString(email)
}

// get the formatted time "now"
func GetTimeNowSting() string {
	return time.Now().Format(time.RFC850)
}

// convert a number to a thousand string
func NumberToString(n int) string {
	s := strconv.Itoa(n)

	startOffset := 0
	var buff bytes.Buffer

	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}

	l := len(s)
	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}

	for i := startOffset; i < l; i++ {
		if (commaIndex == 3) {
			buff.WriteRune(',')
			commaIndex = 0
		}
		commaIndex++
		buff.WriteByte(s[i])
	}
	return buff.String()
}

