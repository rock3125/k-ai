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

package db

import (
	"github.com/gocql/gocql"
	"k-ai/util"
	"k-ai/logger"
	"fmt"
	"time"
	"strings"
	"strconv"
	"errors"
)

// the cassandra system
type CCassandra struct {
	keyspace string
	host string
	replicationFactor int
	Initialised bool

	cluster *gocql.ClusterConfig
	Session *gocql.Session
}

// Freebase word to id or vice versa
type WordId struct {
	Word string
	Id   int
	Is_predicate bool
}

// setup the cassandra
var Cassandra = CCassandra{}


// create a keyspace if it does not exist
func createKeyspace(host string, keyspace string, replicationFactor int) error {

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "system"
	cluster.Timeout = 20 * time.Second

	sess, err := cluster.CreateSession()
	for err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			logger.Log.Warning("Cassandra not responding, waiting %d seconds.  ", 5)
			time.Sleep(5 * time.Second)
			sess, err = cluster.CreateSession()
		} else {
			return err
		}
	}
	str := fmt.Sprintf("CREATE KEYSPACE %s WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor': %d}",
								keyspace, replicationFactor)
	err = sess.Query(str).Exec()
	if err != nil && fmt.Sprintf("%T",err) != "*gocql.RequestErrAlreadyExists" {
		return err
	} else if err == nil {
		logger.Log.Info(fmt.Sprintf("Cassandra: create keyspace %q, rf:%d", keyspace, replicationFactor))
	}
	defer sess.Close()
	return nil
}

// create a keyspace if it does not exist
func DropKeyspace(host string, keyspace string) error {

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "system"
	cluster.Timeout = 20 * time.Second

	sess, err := cluster.CreateSession()
	if err != nil { return err }
	str := fmt.Sprintf("DROP KEYSPACE %s", keyspace)
	err = sess.Query(str).Exec()
	if err == nil {
		logger.Log.Info(fmt.Sprintf("Cassandra: DROP keyspace %q\n", keyspace))
	}
	defer sess.Close()
	return nil
}

// setup db tables etc.
func (c *CCassandra) setupDb() error {
	cql, err := util.LoadTextFile(util.GetDataPath() + "/cql/database.cql")
	if err != nil {
		return err
	}
	cmd_str := ""
	for _, line := range strings.Split(cql, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 && !strings.HasPrefix(line, "//") {
			cmd_str += strings.TrimSpace(" " + line)
			if strings.HasSuffix(cmd_str, ");") {
				cmd_str = strings.Replace(cmd_str, "<ks>", c.keyspace, 1)
				err := c.Session.Query(cmd_str).Exec()
				if err != nil {
					return err
				}
				cmd_str = ""   // reset
			}
		}
	}
	return nil
}

// setup a cassandra connection
func (c *CCassandra) InitCassandraConnection(host string, keyspace string, replicationFactor int) error {

	logger.Log.Info(fmt.Sprintf("Cassandra: connecting to %q, keyspace %q, rf:%d", host, keyspace, replicationFactor))
	createKeyspace(host, keyspace, replicationFactor)

	c.cluster = gocql.NewCluster(host)
	c.cluster.Keyspace = keyspace
	c.cluster.Consistency = gocql.Quorum
	c.cluster.Timeout = 120 * time.Second
	c.host = host
	c.replicationFactor = replicationFactor
	c.keyspace = keyspace

	var err error
	c.Session, err = c.cluster.CreateSession()
	if err != nil { return err }

	// setup tables
	err = c.setupDb()
	if err != nil { return err }
	c.Initialised = true  // db has initialised

	logger.Log.Info("Cassandra setup: done")
	return nil
}

// convert a type to a CQL string part
func typeToString(value interface{}) string {
	str := ""
	switch v := value.(type) {
	case int:
		str += strconv.Itoa(int(v))
	case string:
		// replace single quotes with double single quotes for Cassandra
		str += "'" + strings.Replace(string(v), "'", "''", -1) + "'"
	case gocql.UUID:
		str += gocql.UUID(v).String()
	case *gocql.UUID:
		if v == nil {
			return "<nil>"
		} else {
			str += v.String()
		}
	case float64:
		str += fmt.Sprintf("%f", float64(v))
	case float32:
		str += fmt.Sprintf("%f", float32(v))
	case int64:
		str += fmt.Sprintf("%d", int64(v))
	default:
		logger.Log.Error(fmt.Sprintf("unknown type %T", v))
		panic(fmt.Sprintf("unknown type %T", v))
	}
	return str
}

// convert a type to a CQL string part
func TypeToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(int(v))
	case string:
		// replace single quotes with double single quotes for Cassandra
		return "'" + strings.Replace(string(v), "'", "''", -1) + "'"
	case gocql.UUID:
		return gocql.UUID(v).String()
	case *gocql.UUID:
		if v == nil {
			return "null"
		} else {
			return v.String()
		}
	case float64:
		return fmt.Sprintf("%f", float64(v))
	case float32:
		return fmt.Sprintf("%f", float32(v))
	case int64:
		return fmt.Sprintf("%d", int64(v))
	default:
		return ""
	}
}

/**
 * setup a simple select for a column family
 * @param cf the column family to select from
 * @param columns a set of columns or null
 * @return the select statement
 */
func (c *CCassandra) SelectPaginated(cf string, columns []string,
					whereSet map[string]interface{}, paginationField string,
					paginationValue interface{}, pageSize int) string {
	str := "SELECT "
	if len(columns) > 0 {
		for i, col := range columns {
			if i > 0 {
				str += ","
			}
			str += col
		}
	} else {
		str += "*"
	}
	str += " FROM " + c.keyspace + "." + cf
	counter := 0
	for name, value := range whereSet {
		if counter == 0 {
			str += " WHERE "
		}
		if counter > 0 {
			str += " AND "
		}
		str += name + "=" + typeToString(value)
		counter += 1
	}
	if len(paginationField) > 0 {
		pv := typeToString(paginationValue)
		if pv != "<nil>" && pv != "''"{
			if counter > 0 {
				str += " AND " + paginationField + ">" + pv
			} else {
				str += " WHERE " + paginationField + ">" + pv
			}
		}
	}
	if pageSize > 0 {
		str += " LIMIT " + strconv.Itoa(pageSize)
	}
	str += ";"
	return str
}

/**
 * setup a simple delete for a column family
 * @return the delete statement
 */
func (c *CCassandra) Delete(cf string, whereSet map[string]interface{}) string {
	str := "DELETE FROM " + c.keyspace + "." + cf
	counter := 0
	for name, value := range whereSet {
		if counter == 0 {
			str += " WHERE "
		}
		if counter > 0 {
			str += " AND "
		}
		str += name + "="
		str += typeToString(value)
		counter += 1
	}
	str += ";"
	return str
}

/**
 * setup a simple select for a column family
 * @param cf the column family to select from
 * @param columns a set of columns or null
 * @return the select statement
 */
func (c *CCassandra) Insert(cf string, valueSet map[string]interface{}) string {
	str := "INSERT INTO " + c.keyspace + "." + cf
	counter := 0
	names := ""
	values := ""
	for name, value := range valueSet {
		if counter == 0 {
			names += " (" + name
			values += " VALUES ("
		}
		if counter > 0 {
			names += "," + name
			values += ","
		}
		values += typeToString(value)
		counter += 1
	}
	str += names + ")" + values + ");"
	return str
}

/**
 * execute the statement with the allowance for re-tries
 * @param statement the statement to execute
 * @return the result set or throws an exception eventually after timeoutRetryCount reaches 0
 */
func (c *CCassandra) ExecuteWithRetry(str string) error {
	retryCount := 10
	for retryCount > 0 {
		var err error
		err = c.Session.Query(str).Exec()
		if err != nil {
			retryCount = retryCount - 1;
			if retryCount > 0 {
				logger.Log.Error("execute timed-out (%s), retrying %d more times (waiting 5 seconds: %s)\n", err.Error(), retryCount, str)
				time.Sleep(5 * time.Second)
			} else {
				return errors.New(fmt.Sprint("db exception %s", err.Error()))
			}
		} else {
			return nil // done
		}
	}
	return nil
}

// return a set of words from the freebase lookup system converted from string to id
func (c *CCassandra) FreebaseStringsToIdList(word_list []string) ([]WordId, error) {
	cql_str := "select word, id, is_predicate from " + c.keyspace + ".freebase_word "
	cql_str += "where word in ("
	for i, word := range word_list {
		cql_str += "'" + strings.ToLower(word) + "'"
		if i + 1 < len(word_list) {
			cql_str += ","
		}
	}
	cql_str += ");"

	iter := c.Session.Query(cql_str).Iter()

	var id int
	var word string
	var is_predicate bool

	return_list := make([]WordId,0)
	for iter.Scan(&word, &id, &is_predicate) {
		return_list = append(return_list, WordId{Word: word, Id: id, Is_predicate: is_predicate})
	}
	return return_list, iter.Close()
}


// return a set of ids from the freebase lookup system converted from id to string
func (c *CCassandra) FreebaseIdsToStringList(id_list []int) ([]WordId, error) {
	cql_str := "select word, id, is_predicate from " + c.keyspace + ".freebase_vocab "
	cql_str += "where id in ("
	for i, id := range id_list {
		cql_str += strconv.Itoa(id)
		if i + 1 < len(id_list) {
			cql_str += ","
		}
	}
	cql_str += ");"

	iter := c.Session.Query(cql_str).Iter()

	var id int
	var word string
	var is_predicate bool

	return_list := make([]WordId,0)
	for iter.Scan(&word, &id, &is_predicate) {
		return_list = append(return_list, WordId{Word: word, Id: id, Is_predicate: is_predicate})
	}
	return return_list, iter.Close()
}


