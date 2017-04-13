from typing import List
import spacy
import json
import os

#
# Copyright (c) 2017 by Peter de Vocht
#
# All rights reserved. No part of this publication may be reproduced, distributed, or
# transmitted in any form or by any means, including photocopying, recording, or other
# electronic or mechanical methods, without the prior written permission of the publisher,
# except in the case of brief quotations embodied in critical reviews and certain other
# noncommercial uses permitted by copyright law.
#


logging.info("loading spacy...")
en_nlp = spacy.load('en_core_web_sm')
logging.info("loading spacy done!")


# sentence holder, this is what is returned
class Token:
    def __init__(self, text, index, tag, dep, ancestor_list):
        self.text = text                        # text of the token
        self.index = index                      # index of the token in the document 0..n
        self.dep = dep                          # the name of the SRL dependency
        self.tag = tag                          # penn tag, ucase
        self.ancestor_list = ancestor_list      # dependency tree parent list
        self.synid = -1                         # synset id (default -1, not set)


# simple json encoder / decoder
class JsonSystem(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, Token):
            return {'text': obj.text, 'index': obj.index, 'synid': obj.synid,
                    'tag': obj.tag, 'dep': obj.dep, 'list': obj.ancestor_list}
        return json.JSONEncoder.default(self, obj)


# the text parser
class Parser:

    # cleanup text to ASCII
    def cleanup_text(self, data) -> str:
        try:
            return data.decode("utf-8")
        except:
            text = []
            for ch in data:
                if 32 <= ch <= 255:
                    text.append(chr(ch))
                else:
                    text.append(" ")
            return ''.join(text)


    # convert from spacy to the above Token format for each sentence
    def convert_sentence(self, sent) -> List[Token]:
        sentence = []
        for token in sent:
            ancestors = []
            for an in token.ancestors:
                ancestors.append(an.i)
            text = str(token)
            sentence.append(Token(text, token.i, token.tag_, token.dep_, ancestors))
        filtered_sentence = []
        for token in sentence:
            if token.text != " ":
                filtered_sentence.append(token)
        return filtered_sentence

    # convert a document to a set of entity tagged, pos tagged, and dependency parsed entities
    def parse_document(self, text) -> List[List[Token]]:
        doc = en_nlp(text)
        sentence_list = []
        for sent in doc.sents:
            sentence = self.convert_sentence(sent)
            sentence_list.append(sentence)
        return sentence_list
