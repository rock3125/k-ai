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

package logger

type (
	Logger interface {
		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warning(format string, args ...interface{})
		Error(format string, args ...interface{})
		String() string  // various meanings for different loggers
	}
)

var Log Logger
