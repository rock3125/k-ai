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

package aiml

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"errors"
	"k-ai/nlu/model"
	"k-ai/util"
	"k-ai/db/db_model"
	"k-ai/logger"
	"sync"
	"path"
	"fmt"
)


// the manager system
type AimlManager struct {
	NodeSet map[string]*model.Aiml

	sync.Mutex
}


// singleton access to the Grammar library system
var Aiml AimlManager

// read the AIML xml files from the system
func (mgr *AimlManager) initFromFile() error {
	if mgr.NodeSet == nil || len(mgr.NodeSet) == 0 {
		mgr.NodeSet = make(map[string]*model.Aiml, 0)

		// get the list of AI/ml files to load
		// careful: won't work with   <?xml version="1.0" encoding="ISO-8859-1"?>
		glob_str := util.GetDataPath() + "/aiml/*.aiml"
		file_list, err := util.GetFilesInDirectory(glob_str)
		if err != nil { return err }
		logger.Log.Info(fmt.Sprintf("NLU: loading %s", glob_str))
		for _, fileName := range file_list {
			// read the xml file
			base_filename := path.Base(fileName)
			if strings.Contains(base_filename, ".") {
				base_filename = strings.Split(base_filename, ".")[0]
			}
			if len(base_filename) == 0 || base_filename == ".." {
				base_filename = "K/AI"
			} else {
				base_filename += " module"
			}
			xmlFile, err := ioutil.ReadFile(fileName)
			if err != nil { return err }

			// xml to internal structure
			var category Categories
			xml.Unmarshal(xmlFile, &category)

			// go through each category in the file and expand into the node system
			for _, cat := range category.Cats {
				templateList := strings.Split(cat.Template, "|")
				finalTemplateList := make([]string, 0)
				for _, template := range templateList {
					template := strings.TrimSpace(template)
					if len(template) > 0 {
						finalTemplateList = append(finalTemplateList, template)
					}
				}
				mgr.AddPattern(cat.PatternList, base_filename, finalTemplateList)
			}
		}
		logger.Log.Info("NLU: aiml loading done")
	}
	return nil
}

// read the AIML xml files from the db system, BUT ONLY AFTER the load from file above
func (mgr *AimlManager) SetupDbSchema() error {
	if mgr.NodeSet == nil || len(mgr.NodeSet) == 0 {
		return errors.New("aiml not initialized")
	}
	// get a list of all schema items
	schema_list, err := db_model.GetSchemaList()
	logger.Log.Info("NLU: loading schemas from db")
	if err != nil {
		return err
	} else {
		for _, schema_item := range schema_list {
			for _, field := range schema_item.Field_list {
				if len(field.Aiml) > 0 {
					aiml_list := make([]string,0)
					for _, aiml_str := range strings.Split(field.Aiml, "\n") {
						aiml_str = strings.TrimSpace(aiml_str)
						if len(aiml_str) > 0 {
							aiml_list = append(aiml_list, aiml_str)
						}
					}
					if len(aiml_list) > 0 {
						template_list := make([]string, 0)
						kbStr := "db_search:"+schema_item.Name+":"+field.Name
						template_list = append(template_list, kbStr)
						mgr.AddPattern(aiml_list, schema_item.Origin, template_list)
					}
				}
			}
		}
	}
	logger.Log.Info("NLU: aiml db loading done")
	return nil
}

// todo: cleanup and fix
// reload the entire AIML system - for now - incredibly badly written
// lock is held open way too long - however - it'll do in the proto-type
func (mgr *AimlManager) Reload() error {
	mgr.Lock()
	defer mgr.Unlock()

	// re-create the map
	mgr.NodeSet = make(map[string]*model.Aiml, 0)

	err := mgr.initFromFile()
	if err != nil {
		return err
	}

	err = mgr.SetupDbSchema()
	if err != nil {
		return err
	}

	return nil
}


// setup
func init() {
	// init aiml system
	Aiml = AimlManager{}
	// and load
	err := Aiml.initFromFile()
	if err != nil {
		panic(err)
	}
}

