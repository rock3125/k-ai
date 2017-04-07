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

import "encoding/xml"

// xml structure
type Categories struct {
	XMLName xml.Name	`xml:"aiml"`
	Version string   	`xml:"version,attr"`
	Cats []Category		`xml:"category"`
}

type Category struct {
	XMLName xml.Name 		`xml:"category"`
	Template string			`xml:"template"`
 	PatternList []string 	`xml:"pattern"`
}


