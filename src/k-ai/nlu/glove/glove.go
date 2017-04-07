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
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"k-ai/logger"
	"math"
	"errors"
)

// example use

//g := glove.CreateGloveMap(util.GetDataPath() + "/glove.6B.50d.txt", 1000000)
//
//// test a vector and get its closest relatives
//relations := g.GetRelatedWords("peter", 1)
//for _, value := range relations {
//	Printf("score=%f, %s\n", value.Score, value.Text)
//}


// load a set of glove vectors from the file specified
// the glove file is a simple one vector per line with format:
// word_str f64 f64 ... f64
// returned is a map of word_str -> []float64 of size whatever the file is
// the file must be consistent, i.e. all vectors are the same size in the file
func LoadGlove(filename string) (map[string][]float64, error) {

	g_map := make(map[string][]float64)

	// open the file and check it
	file, err := os.Open(filename)
	if err != nil { return nil, err }
	defer file.Close()

	vector_size := 0  // initial size
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())  // split
		this_vector_size := len(words) - 1
		if vector_size == 0 {
			vector_size = this_vector_size
		} else {
			if this_vector_size != vector_size {
				return nil, errors.New(fmt.Sprintf("inconsistent vector size, expected vectors of size %d but got %d", vector_size, this_vector_size))
			}
		}
		word := words[0]
		vector := make([]float64, vector_size, vector_size)
		for i := 0; i < vector_size; i++ {
			vector[i], _ = strconv.ParseFloat(words[i+1], 64)
		}
		g_map[word] = vector
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return g_map, nil
}


// return the angle between two vectors - their relatedness
func GetVectorAngle(v1 []float64, v2 []float64) (float64) {
	value := 0.0
	l_v1 := 0.0
	l_v2 := 0.0
	for i := 0; i < len(v2); i++ {
		l_v1 += v1[i] * v1[i]   // length of v1 and v2
		l_v2 += v2[i] * v2[i]
		value += v1[i] * v2[i]  // dot product
	}
	dotp := value / (math.Sqrt(l_v1) * math.Sqrt(l_v2))
	return math.Cosh(dotp)
}

// create a glove lookup map system
func CreateGloveMap(filename string, index_size int) (*MapSystem, error) {
	logger.Log.Info("loading glove vectors ")
	g_map, err := LoadGlove(filename)
	if err != nil { return nil, err }
	logger.Log.Info(fmt.Sprintf("%d vectors loaded, setup", len(g_map)))
	result := new(MapSystem)
	result.setup(g_map, float64(index_size))
	logger.Log.Info(" done")
	return result, nil
}

