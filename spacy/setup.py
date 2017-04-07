#!/usr/bin/env python3

#
# Copyright (c) 2017 by Peter de Vocht
#
# All rights reserved. No part of this publication may be reproduced, distributed, or
# transmitted in any form or by any means, including photocopying, recording, or other
# electronic or mechanical methods, without the prior written permission of the publisher,
# except in the case of brief quotations embodied in critical reviews and certain other
# noncommercial uses permitted by copyright law.
#

from setuptools import setup

setup(
    name='kai-parser',
    version='1.0',
    description='KAI AI natural language parser',
    author='Peter de Vocht',
    author_email='peter@peter.co.nz',
    long_description='The kai-parser is a micro-service around the spaCy parser for KAI.',
    url='https://github.com/peter3125/kai',
    keywords='KAI natural language parser',
    license='http://www.apache.org/licenses/LICENSE-2.0',
    classifiers=['License :: OSI Approved :: Apache Software License','Natural Language :: English',
                 'Operating System :: OS Independent','Programming Language :: Python :: 3 :: Only'],
    requires=['Flask','nltk','numpy','spacy','sputnik'],
    # package_dir={'': 'src'},
    packages=['kai', 'kai/parser', 'kai/test'],
    package_data={'kai': ['kai-parser.sh']},
    test_suite='kai',
 )
