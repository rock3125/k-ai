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

package environment

import (
	"os"
	"fmt"
	"github.com/BurntSushi/toml"
	"k-ai/util"
)

// K/AI Configuration structure to match properties.ini
type Config struct {

	Version string  // version of K/AI

	Logger string // logger, stdout or path to file

	// web server host port and cert details
	KaiServerPort int
	KaiUseHTTPS bool
	KaiCertLocation string
	KaiKeyLocation string

	// cassandra
	CassandraServer string
	Keyspace string
	ReplicationFactor int

	// spacy
	SpacyEndpoint string
}

// Reads info from config file
func ReadConfig() Config {
	var config_file = util.GetDataPath() + "/properties.ini"
	_, err := os.Stat(config_file)
	if err != nil {
		panic(fmt.Sprintf("Config file is missing: %s", config_file))
	}
	var config Config
	if _, err := toml.DecodeFile(config_file, &config); err != nil {
		panic(err)
	}
	return config
}
