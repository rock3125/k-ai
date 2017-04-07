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
	"k-ai/logger"
	"net/http"
	"strconv"
	"os"
	"errors"
	"fmt"
)

// start the server
func StartServer(port int, ssl bool, cert_file string, key_file string, version string) error {
	if ssl {
		if _, err := os.Stat(cert_file); os.IsNotExist(err) {
			return errors.New("cert file not found: " + cert_file)
		}
		if _, err := os.Stat(key_file); os.IsNotExist(err) {
			return errors.New("key file not found: " + key_file)
		}
	}
	router := NewRouter(port, version) // setup the router

	if ssl {
		logger.Log.Info(fmt.Sprintf("Starting REST server at https://localhost:%d\n", port))
		return http.ListenAndServeTLS(":"+strconv.Itoa(port), cert_file, key_file, router)
	} else {
		logger.Log.Info(fmt.Sprintf("Starting REST server at http://localhost:%d\n", port))
		return http.ListenAndServe(":"+strconv.Itoa(port), router)
	}
}

