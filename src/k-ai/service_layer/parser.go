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
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "k-ai/nlu/parser"
    "k-ai/nlu/model"
)

// parse a piece of raw text from the query string
//
func Parse(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    text_to_parse := vars["text"]
    sentenceList, err := parser.ParseText(text_to_parse)
    if err != nil { panic(err) }
    tupleList := make([]string,0)
    for _, sentence := range sentenceList {
        tupleList = append(tupleList, model.SentenceToTuple(sentence).ToStringIndent())
    }
    if err := json.NewEncoder(w).Encode(tupleList); err != nil {
        panic(err)
    }
}

// parse the text and return a PNG
//
func ParseToPng(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    text_to_parse := vars["text"]

    ttList, err := parser.ParseTextToTupleTree(text_to_parse)
    if err != nil {
        w.Write([]byte("error parsing " + err.Error()))
    } else if len(ttList) > 0 {
        png, err := model.ToPng(ttList.ToGraphVizDot())
        if err != nil {
            w.Write([]byte("error parsing: " + err.Error()))
        } else {
            w.Header().Set("Content-Type", "image/png")
            w.Write(png)
        }
    } else {
        w.Write([]byte("error parsing"))
    }
}

