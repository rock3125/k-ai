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
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "k-ai/service_layer"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    Example     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route


// server main entry point
func ServiceLayerIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<h2>K/AI server</h2>\n\n<h3>service layer access points</h3>\n\n%s<br/><br/>" +
        "<small>Copyright Peter de Vocht, 2017.  K/AI version %s</small>", infoStr, versionStr)
}



var routes = Routes {
    Route{
        "This message",
        "GET",
        "/sl/",
        "/sl/",
        ServiceLayerIndex,
    },
    Route{
        "Construct a tuple for a parseable piece of text",
        "GET",
        "/sl/parse/{text}",
        "/sl/parse/Mark lives in his car.",
        service_layer.Parse,
    },
    Route{
        "Construct a parser png (image) for a parseable piece of text",
        "GET",
        "/sl/parse-to-png/{text}",
        "/sl/parse-to-png/Peter and Sherry went to the beach at 12:45 to view the boats comming in.",
        service_layer.ParseToPng,
    },

    /////////////////////////////////////////////////////////////////
    // KB ui

    Route{
        "Get a list of entities",
        "GET",
        "/kb-entity/get_list/{session}/{topic}/{prev}/{page_size}/{json_field}/{query_str}",
        "",
        service_layer.ListEntities,
    },
    Route{
        "Delete an existing entity",
        "DELETE",
        "/kb-entity/delete/{session}/{topic}/{id}",
        "",
        service_layer.DeleteEntity,
    },
    Route{
        "Save an existing entity",
        "POST",
        "/kb-entity/save/{session}",
        "",
        service_layer.SaveEntity,
    },
    Route{
        "Upload a set of instances",
        "POST",
        "/kb-entity/upload/{session}/{id}",
        "",
        service_layer.UploadInstances,
    },

    /////////////////////////////////////////////////////////////////
    // query ui

    Route{
        "ASK a Question",
        "POST",
        "/ask/{session}",
        "",
        service_layer.Ask,
    },
    Route{
        "TEACH K/AI a factoid",
        "POST",
        "/teach/{session}",
        "",
        service_layer.Teach,
    },
    Route{
        "REMOVE an existing factoid by id",
        "DELETE",
        "/remove/factoid/{session}/{id}",
        "",
        service_layer.UnTeach,
    },

    /////////////////////////////////////////////////////////////////
    // semantic entities

    Route{
        "Semantic Entities: find",
        "GET",
        "/entities/find/{session}/{name}",
        "",
        service_layer.FindSemanticEntities,
    },
    Route{
        "Semantic Entities: delete",
        "DELETE",
        "/entities/delete/{session}/{name}",
        "",
        service_layer.DeleteSemanticEntity,
    },
    Route{
        "Semantic Entities: save",
        "GET",
        "/entities/save/{session}/{name}/{semantic}",
        "",
        service_layer.SaveSemanticEntity,
    },

    /////////////////////////////////////////////////////////////////
    // topic entities (unstructured topic data)

    Route{
        "Return a list of topics",
        "GET",
        "/topic/get_list/{session}/{page}/{page_size}/{filter_text}",
        "",
        service_layer.GetTopicList,
    },
    Route{
        "Remove a topic from the system",
        "DELETE",
        "/topic/delete/{session}/{topic_name}",
        "",
        service_layer.DeleteTopic,
    },
    Route{
        "Save a topic (insert or update)",
        "POST",
        "/topic/save/{session}/{topic_name}",
        "",
        service_layer.SaveTopic,
    },

    /////////////////////////////////////////////////////////////////
    // topic entities (unstructured topic data)

    Route{
        "Create a new user account",
        "POST",
        "/user/create",
        "",
        service_layer.CreateUser,
    },

    Route{
        "Signin an existing user",
        "POST",
        "/user/signin",
        "",
        service_layer.Signin,
    },

    Route{
        "Signout an existing user",
        "GET",
        "/user/signout/{session}",
        "",
        service_layer.Signout,
    },

}

// this becomes the info / feedback string for the / part of the interface
var infoStr = ""
var versionStr = ""

// setup the RESTful router
func NewRouter(port int, version string) *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    versionStr = version  // set version
    str := "<table style=\"width:100%\">\n<tr><th>method</th><th>url</th><th>description</th></tr>\n"
    for _, route := range routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
        if len(route.Example) > 0 {
            fmtStr := "<tr><td>%s</td><td><a href=\"%s\" target='_blank'>http://%s:%d%s</a></td><td>%s</td></tr>\n"
            str += fmt.Sprintf(fmtStr, route.Method, route.Example, "localhost", port, route.Pattern, route.Name)
        } else {
            fmtStr := "<tr><td>%s</td><td>http://%s:%d%s</td><td>%s</td></tr>\n"
            str += fmt.Sprintf(fmtStr, route.Method, "localhost", port, route.Pattern, route.Name)
        }
    }
    infoStr = str + "\n</table>\n"

    //////////////////////////////////////////////////////////////////////
    // add a static router for the web and all other pages
    var handlerFunc http.HandlerFunc = StaticIndex
    router.
        Methods("GET").
        Path("/{url}").
        Name("static web root specific").
        Handler(handlerFunc)

    router.
    Methods("GET").
        Path("/").
        Name("static web root /").
        Handler(handlerFunc)

    router.
    Methods("GET").
        Path("/css/{url}").
        Name("static web css").
        Handler(handlerFunc)

    router.
    Methods("GET").
        Path("/images/{url}").
        Name("static web images").
        Handler(handlerFunc)

    router.
    Methods("GET").
        Path("/js/{url}").
        Name("static web javascript").
        Handler(handlerFunc)

    router.
    Methods("GET").
        Path("/fonts/{url}").
        Name("static web fonts").
        Handler(handlerFunc)

    return router
}

