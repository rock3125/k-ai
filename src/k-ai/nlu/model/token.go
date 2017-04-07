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

package model

import "fmt"

type Token struct {
	Index int               `json:"index"`
	AncestorList []int      `json:"list"`
	Tag   string            `json:"tag"`
	Text  string            `json:"text"`
	Dep   string            `json:"dep"`
	SynId int               `json:"synid"`
	Semantic string         `json:"semantic"`
	Anaphora string			`json:"-"`
}

func (t Token) ToString() (string) {
	return fmt.Sprintf("%#v", t)
}

