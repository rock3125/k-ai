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

package environment

import "testing"

func TestProperties(t *testing.T) {
	config := ReadConfig()
	if config.KaiServerPort == 0 || len(config.CassandraServer) == 0 || len(config.SpacyEndpoint) == 0 {
		t.Error("config file not loaded")
	}
}

