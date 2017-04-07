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
	"k-ai/nlu/tokenizer"
	"strings"
	"k-ai/nlu/model"
)


/**
 * add a pattern for prcessing
 * @param patternList the list of patterns to use
 * @param origin the origin of this aiml pattern (whence it came)
 * @param aimlTemplateList the list of templates to associate with this pattern
 */
func (mgr *AimlManager) AddPattern(patternList []string, origin string, templateList []string) {
	for _, pattern1 := range patternList {
		expandedPattern := expandBrackets(pattern1)
		for _, pattern := range expandedPattern {
			tokenList := tokenizer.FilterOutPunctuation(tokenizer.FilterOutSpaces(tokenizer.Tokenize(pattern)))
			if len(tokenList) > 1 { // must at least have two items in a pattern
				t_token := tokenList[0];
				key := strings.ToLower(t_token.Text)
				if key == "*" {
					panic("error: pattern cannot start with *");
				}
				root, ok := mgr.NodeSet[key];
				if !ok {
					root = &model.Aiml{Text: key, Origin: origin, NodeSet: make(map[string]*model.Aiml,0)}
					mgr.NodeSet[key] = root;
				}
				addPatternHelper(root, origin, 1, tokenList, templateList);
			}
		}
	}
}

/**
 * Expand brackets for (item1|item2|)  (last one's empty)
 * @param str the string to examine and expand
 * @return a list of expansions (or string itself if not the case)
 */
func expandBrackets( str string ) []string {
	resultList := make([]string,0)
	if strings.Contains(str,"(") {
		tokenList := tokenizer.FilterOutSpaces(tokenizer.Tokenize(str))
		builder := make([]string,0)
		sb := ""
		for i := 0; i < len(tokenList); {
			text := tokenList[i].Text
			if text == "(" {

				// finalise the previous results
				builder, sb = finish(builder, sb);

				itemList := make([]string, 0)
				j := i + 1;
				item := ""
				for j < len(tokenList) {
					t2 := tokenList[j].Text
					if t2 == ")" {
						itemList = append(itemList, item)
						j++;
						break;
					} else if t2 == "|" {
						itemList = append(itemList, item)
						item = ""
					} else {
						if len(item) > 0 {
							item += " "
						}
						item += t2;
					}
					j++;
				}

				// generate new list
				newBuilder := make([]string,0)
				for _, str1 := range builder {
					for _, str2 := range itemList {
						str3 := str1 + " " + str2;
						newBuilder = append(newBuilder, strings.TrimSpace(str3))
					}
				}
				builder = newBuilder;
				i = j; // advance
			} else {
				if len(sb) > 0 {
					sb += " ";
				}
				sb += text;
				i++;
			}
		}
		// finalise the results
		builder,_ = finish(builder, sb);
		return builder
	} else {
		// no ( | )
		resultList = append(resultList, str)
	}
	return resultList;
}



/**
 * finalise dealing with the builder string given a string builder that has
 * been collecting information
 * @param builder the builder to add sb to
 * @param sb the string builder
 * @return the modified builder with sb appended
 */
func finish(builder []string, sb string) ([]string, string) {
	if len(sb) > 0 {
		// add the current sb content to all previous builder items
		if len(builder) == 0 {
			builder = append(builder, sb)
			sb = ""
		} else {
			newBuilder := make([]string,0)
			for _, str1 := range builder {
				str3 := str1 + " " + sb;
				newBuilder = append(newBuilder, strings.TrimSpace(str3))
			}
			return newBuilder, ""
		}
	} else if len(builder) == 0 {
		// make sure the builder has an initial value to proceed with
		builder = append(builder, "")
	}
	return builder, ""
}


/**
 * process the patterns and create a tree of patterns that can be matched to user input strings
 * @param nodeSet the parent set of nodes - recursively updated
 * @param index the index into tokenList
 * @param tokenList the list of token making up the pattern
 * @param templateList the list of template to be added to the last node
 */
func addPatternHelper( nodeSet *model.Aiml, origin string, index int, tokenList []model.Token, templateList []string ) {
	if index + 1 == len(tokenList) { // last item insert
		// last node
		t_token := tokenList[index]
		key := strings.ToLower(t_token.Text)
		template, ok := nodeSet.NodeSet[key];
		if !ok { // new item
			template = &model.Aiml{Text: key, Origin: origin, NodeSet: make(map[string]*model.Aiml,0)}
			template.TemplateList = make([]string,0)
			for _, item := range templateList {
				template.TemplateList = append(template.TemplateList, item)
			}
			if len(template.TemplateList) != len(templateList) {
				panic("assert: len == len after insert")
			}
			nodeSet.NodeSet[key] = template
		} else { // existing template - all these sets as alternatives
			for _, item := range templateList {
				template.TemplateList = append(template.TemplateList, item)
			}
		}
	} else if index < len(tokenList) {
		// in between node
		t_token := tokenList[index]
		key := strings.ToLower(t_token.Text)
		template, ok := nodeSet.NodeSet[key]
		if !ok {
			template = &model.Aiml{Text: key, Origin: origin, NodeSet: make(map[string]*model.Aiml,0)}
			nodeSet.NodeSet[key] = template
		}
		addPatternHelper(template, origin, index + 1, tokenList, templateList);
	}
}

