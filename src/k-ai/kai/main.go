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

package main

import (
	"fmt"
	"k-ai/db"
	"k-ai/rest"
	"k-ai/logger"
	"k-ai/nlu/parser"
	"k-ai/nlu/aiml"
	"k-ai/environment"
	"k-ai/db/freebase"
)


// K/AI main start
func main() {

	// get configuration
	env := environment.ReadConfig()

	logger.Log.Info("******  start  ******")

	logger.Log.Info(fmt.Sprintf("K/AI System, version %s", env.Version))

	// init cassandra
	logger.Log.Info(fmt.Sprintf("connecting to Cassandra %s @ %s", env.Keyspace, env.CassandraServer))
	db.Cassandra.InitCassandraConnection(env.CassandraServer, env.Keyspace, env.ReplicationFactor)

	logger.Log.Info("Setting up Freebase Match System")
	err := freebase.MatchSystem.Setup()
	if err != nil {
		logger.Log.Error("Error connecting to Spacy %s", err.Error())
		return
	}

	// test spacy is up
	logger.Log.Info(fmt.Sprintf("connecting to Spacy @ %s", env.SpacyEndpoint))
	parser.SpacyEndpoint = env.SpacyEndpoint
	sl, err := parser.ParseText("Test text.")
	if err != nil {
		logger.Log.Error("Error connecting to Spacy %s", err.Error())
	} else if len(sl) != 1 || len(sl[0].TokenList) != 3 {
		logger.Log.Error("Error Spacy parser interface not working")
	} else {
		// setup db schema for aiml
		aiml.Aiml.SetupDbSchema()

		//// test freebase
		//sentence_list, err := parser.ParseText("what recordings did bastard souls make?")
		//if err == nil {
		//	answer, err := freebase.FreebaseQueryBySearch(sentence_list[0].TokenList)
		//	if err == nil {
		//		print(answer)
		//	}
		//}

		// start server on port
		err = rest.StartServer(env.KaiServerPort, env.KaiUseHTTPS, env.KaiCertLocation, env.KaiKeyLocation, env.Version)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

}


