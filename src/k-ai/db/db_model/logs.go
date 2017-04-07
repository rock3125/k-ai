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
	"time"
	"k-ai/db"
	"errors"
)

// add a new action log entry to the system recording an action
func AddLogEntry(who string, what string) error {
	if len(who) == 0 || len(what) == 0  {
		return errors.New("LogEntry.Save() invalid parameters")
	}

	value_map := make(map[string]interface{})
	value_map["when"] = time.Now().UnixNano()
	value_map["who"] = who
	value_map["what"] = what

	return db.Cassandra.ExecuteWithRetry(db.Cassandra.Insert("logs", value_map))
}

