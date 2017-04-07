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

import "github.com/gocql/gocql"

// the result of an ask or teach request
type ATResult struct {
	Text string						`json:"text"`
	Timestamp string				`json:"timestamp"`
	Topic string					`json:"topic"`
	Sentence_id gocql.UUID          `json:"sentence_id"`    // id for the item if applicable
	KB_id gocql.UUID          		`json:"kb_id"`
}

// a list of ask teach results with error / message fields
type ATResultList struct {
	Message string				`json:"message"`
	Error string				`json:"error"`
	ResultList []ATResult		`json:"result_list"`
}

