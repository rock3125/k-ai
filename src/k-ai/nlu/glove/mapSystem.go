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

package glove

import (
	"math"
	"sort"
)

// a system for referencing all lookup information
// for a glove vector and its distances
type MapSystem struct {
	g_map map[string][]float64          // the glove string -> []vector
	size float64                        // size of the indexes
	location_array [][]MatchVector      // [index][list of items @ index]
}

// calculate distance squared between two vectors
// and return the distance value and the index
func (m MapSystem)to_index(v2 []float64) (float64, int) {
	value := 0.0
	for i := 0; i < len(v2); i++ {
		diff := v2[i]
		value += diff * diff
	}
	score := math.Sqrt(value) / float64(len(v2))
	return score, int(score * m.size)
}

// return the angle between two vectors - their relatedness
func (m MapSystem) GetVectorAngleUsingNames(v1_name string, v2_name string) (float64) {
	return GetVectorAngle(m.g_map[v1_name], m.g_map[v2_name])
}

// test a vector and get its closest relatives
// word: the word to find relationships for
// dist: the distance around the word 0..x
func (m MapSystem) GetRelatedWords(word string, dist int) (MatchVectors) {

	return_value := make(MatchVectors,0)
	int_size := int(m.size)
	divisor := float64(dist)
	if divisor == 0.0 { divisor = 1.0 }

	if val, ok := m.g_map[word]; ok {  // does it exist in the map?
		_, index := m.to_index(val)  // get the index of the closest neighbour
		if index < int_size {
			start := index - dist
			if start < 0 { start = 0 }
			end := index + dist + 1
			if end > int_size { end = int_size }
			for index_offset := start; index_offset < end; index_offset++ {
				value := m.location_array[index_offset]
				if value != nil {
					for _, mv := range value {
						return_value = append(return_value,
							MatchVector{Text: mv.Text,
								Score: m.GetVectorAngleUsingNames(word, mv.Text),
								Vector: mv.Vector})
					}
				}
			}
		}
	}
	if len(return_value) > 0 {
		sort.Sort(return_value)
	}
	return return_value
}

// setup a the MapSystem with related items and distances
// g_map: the existing map with all glove vectors
// size: the scaling offset 100,000 to 1,000,000 recommended
func (m *MapSystem) setup(g_map map[string][]float64, size float64) {
	m.g_map = g_map
	// pick an arbitrary vector
	m.size = size

	// setup a new slice for the size of the matches
	m.location_array = make([][]MatchVector, int(size), int(size))

	// setup this array
	for key, value := range g_map {
		_, index := m.to_index(value)
		if index < int(size) {
			if m.location_array[index] == nil {
				m.location_array[index] = make([]MatchVector, 0)
			}
			mv := MatchVector {Text: key, Vector: value, Score: 0.0}
			m.location_array[index] = append(m.location_array[index], mv)
		}
	}
}

