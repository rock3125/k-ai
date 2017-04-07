import unittest

#
# Copyright (c) 2017 by Peter de Vocht
#
# All rights reserved. No part of this publication may be reproduced, distributed, or
# transmitted in any form or by any means, including photocopying, recording, or other
# electronic or mechanical methods, without the prior written permission of the publisher,
# except in the case of brief quotations embodied in critical reviews and certain other
# noncommercial uses permitted by copyright law.
#

from kai.parser.parser import Parser


# test the parser
class ParserTest(unittest.TestCase):
    # initialise the class
    def __init__(self, methodName: str):
        unittest.TestCase.__init__(self, methodName)
        self.parser = Parser()

    # test we can parse and words are compounded and spaces are removed
    def test_parser_1(self):
        sentence_list = self.parser.parse_document("Peter de Vocht was here.  He then moved to Wellington.")
        self.assertTrue(len(sentence_list) == 2)

        sentence_1 = sentence_list[0]
        self.assertTrue(len(sentence_1) == 6)

        sentence_2 = sentence_list[1]
        self.assertTrue(len(sentence_2) == 6)
