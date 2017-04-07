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
	"k-ai/util"
	"strings"
	"fmt"
	"path/filepath"
	"io/ioutil"
)

// check a file is valid in the right directory
func isValidPath(path string) bool {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	// check we're inside the system's own data folder
	data_dir_parts := strings.Split(util.GetDataPath() + "/web", "/")
	check_dir_parts := strings.Split(abs_path, "/")
	if len(check_dir_parts) < len(data_dir_parts) { // must be at least the same size
		return false
	}
	for i, part := range  check_dir_parts {
		if i < len(data_dir_parts) {
			if part != data_dir_parts[i] {
				return false
			}
		}
	}
	return true
}

// check from a list of suffixes that the url is right
func isValidSuffix(path string, suffix_list...string) bool {
	for _, suffix := range suffix_list {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}

func StaticIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	main_url := r.URL.String()

	basePath := util.GetDataPath() + "/web/"

	if len(url) == 0 {  // return index

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		str, err := util.LoadTextFile(basePath + "index.html")
		if err != nil {
			fmt.Fprint(w, "File error:" + err.Error())
			return
		}
		str = strings.Replace(str, "$version", versionStr, -1)
		w.Write([]byte(str))

	} else {
		// serve whatever file they wanted
		if strings.Contains(main_url, "/css/") {
			// css
			file_path := basePath + "css/" + url
			if isValidPath(file_path) && util.Exists(file_path) && isValidSuffix(file_path, ".css") {
				w.Header().Set("Content-Type", "text/css")
				str, err := util.LoadTextFile(file_path)
				if err != nil {
					fmt.Fprint(w, "File error:" + err.Error())
					return
				}
				w.Write([]byte(str))
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "File not found 404")
			}
		} else if strings.Contains(main_url, "/images/") {
			// image
			file_path := basePath + "images/" + url
			if isValidPath(file_path) && util.Exists(file_path) && isValidSuffix(file_path, ".png", ".jpg", ".jpeg", ".gif") {
				if isValidSuffix(file_path, ".png") {
					w.Header().Set("Content-Type", "image/png")
				} else if isValidSuffix(file_path, ".jpg", ".jpeg") {
					w.Header().Set("Content-Type", "image/jpeg")
				} else {
					w.Header().Set("Content-Type", "image/gif")
				}
				data, err := ioutil.ReadFile(file_path)
				if err != nil {
					fmt.Fprint(w, "File error:" + err.Error())
					return
				}
				w.Write(data)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "File not found 404")
			}
		} else if strings.Contains(main_url, "/fonts/") {
			// image
			file_path := basePath + "fonts/" + url
			if isValidPath(file_path) && util.Exists(file_path) && isValidSuffix(file_path, ".woff", ".woff2", ".ttf", ".svg", ".eot") {
				w.Header().Set("Content-Type", "application/octet-stream")
				data, err := ioutil.ReadFile(file_path)
				if err != nil {
					fmt.Fprint(w, "File error:" + err.Error())
					return
				}
				w.Write(data)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "File not found 404")
			}
		} else if strings.Contains(main_url, "/js/") {
			// image
			file_path := basePath + "js/" + url
			if isValidPath(file_path) && util.Exists(file_path) && isValidSuffix(file_path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
				data, err := ioutil.ReadFile(file_path)
				if err != nil {
					fmt.Fprint(w, "File error:" + err.Error())
					return
				}
				w.Write(data)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "File not found 404")
			}
		} else {
			// root dir
			file_path := basePath + "" + url
			if isValidPath(file_path) && util.Exists(file_path) && isValidSuffix(file_path, ".html", ".ico") {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				str, err := util.LoadTextFile(file_path)
				if err != nil {
					fmt.Fprint(w, "File error:" + err.Error())
					return
				}
				str = strings.Replace(str, "$version", versionStr, -1)
				w.Write([]byte(str))
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "File not found 404")
			}
		}

	}
}

