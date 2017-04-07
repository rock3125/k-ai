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
	"testing"
)

func TestGloveAngles(t *testing.T) {
	v1 := []float64{1.0, 0.0}
	v2 := []float64{0.0, 1.0}
	angle := GetVectorAngle(v1,v2)
	if angle != 1.0 {
		t.Errorf("angle != 1.0 but %f", angle)
	}
}

