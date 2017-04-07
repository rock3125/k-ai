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
	"k-ai/nlu/model"
	"encoding/json"
	"k-ai/util"
)

// write a json error with header to the output
func JsonError(writer http.ResponseWriter, error_message string) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("{\"error\": \"" + error_message + "\"}"))
}

// write a success message with header to the output
func JsonMessage(writer http.ResponseWriter, info_code int, info_message string) {
	writer.WriteHeader(info_code)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("{\"message\": \"" + info_message + "\"}"))
}

// write a teach/ask json error with header to the output
func ATJsonError(writer http.ResponseWriter, error_message string) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Header().Set("Content-Type", "application/json")
	obj := model.ATResultList{Error: error_message}
	json_bytes, _ := json.Marshal(obj)
	writer.Write(json_bytes)
}

// write a teach/ask success message with header to the output
func ATJsonMessage(writer http.ResponseWriter, info_code int, info_message string) {
	writer.WriteHeader(info_code)
	writer.Header().Set("Content-Type", "application/json")
	obj := model.ATResultList{ ResultList: make([]model.ATResult,0) }
	obj.ResultList = append(obj.ResultList, model.ATResult{Text: info_message,
		Timestamp: util.GetTimeNowSting(), Topic: "K/AI"})
	json_bytes, _ := json.Marshal(obj)
	writer.Write(json_bytes)
}

