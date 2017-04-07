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

package service_layer

import (
	"net/http"
	"github.com/gorilla/mux"
	"strings"
	"github.com/gocql/gocql"
	"k-ai/db/db_model"
)

// remove a previously indexed factoid
//
func UnTeach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// check session is a valid
	session := strings.ToLower(strings.TrimSpace(vars["session"]))
	session_obj, err := db_model.ValidateSession(session)
	if err != nil {
		ATJsonError(w, err.Error())
		return
	}
	username := session_obj.GetUserName()

	// check id
	factoid_str := vars["id"]
	id, err := gocql.ParseUUID(factoid_str)
	if err != nil {
		ATJsonError(w, "Invalid id:"+factoid_str)
		return
	}

	// log the event
	db_model.AddLogEntry(username, "unteach:" + username + "/" + factoid_str)

	err = db_model.DeleteText(&id, username)
	if err != nil {
		ATJsonError(w, "DeleteText("+err.Error() + ")")
		return
	}

	db_model.RemoveIndexes(id, username)
	if err != nil {
		ATJsonError(w, "RemoveIndexes(" + username + "," + err.Error() + ")")
		return
	}

	db_model.RemoveIndexes(id, "global")
	if err != nil {
		ATJsonError(w, "DeleteText(global,"+err.Error() + ")")
		return
	}
	ATJsonMessage(w, http.StatusOK, "removed factoid " + id.String())
}

