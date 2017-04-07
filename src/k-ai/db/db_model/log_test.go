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

package db_model

import (
	"testing"
	"k-ai/db"
)

// perform further index tests multiple keyword
func TestLog1(t *testing.T) {

	// init cassandra
	db.DropKeyspace("localhost", "kai_ai_log")
	db.Cassandra.InitCassandraConnection("localhost", "kai_ai_log", 1)

	AddLogEntry("peter", "This is a add to log test")

	db.DropKeyspace("localhost", "kai_ai_log")
}
