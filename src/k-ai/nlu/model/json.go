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

// {"num_sentences": 1, "packetList": [
// {"metadata": "body", "spacyTokenList": {"num_sentences": 1, "processing_time": 0, "num_tokens": 2,
// "sentence_list": [[{"index": 0, "list": ["1"], "tag": "NN", "text": "test", "dep": "compound", "synid": -1},
//                    {"index": 1, "list": [], "tag": "NN", "text": "sentence", "dep": "ROOT", "synid": -1}]],
// ]}

type SpacyList struct {
	NumSentences int        `json:"num_sentences"`
	Processing_time int64   `json:"processing_time"`
	NumTokens int           `json:"num_tokens"`
	SentenceList [][]Token  `json:"sentence_list"`
}
