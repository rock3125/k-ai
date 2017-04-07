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
	"k-ai/db"
)

// for unit tests only
func Delete_and_create_keyspace_for_unit_test(keyspace string) {
	// init cassandra
	db.DropKeyspace("localhost", keyspace)
	db.Cassandra.InitCassandraConnection("localhost", keyspace, 1)
}

// for unit tests only
func Delete_keyspace_after_unit_test(keyspace string) {
	// init cassandra
	db.DropKeyspace("localhost", keyspace)
}

