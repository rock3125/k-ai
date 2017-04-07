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

// aiml structure node
type Aiml struct {
	Text         string           // the text to match on
	Origin       string           // whence it came or its origin
	TemplateList []string         // possible answers
	NodeSet      map[string]*Aiml // other nodes of this node
}

// binding
type AimlBinding struct {
	Text      string // the text
	Origin    string // whence it came
	Offset    int
	TokenList []Token
}

