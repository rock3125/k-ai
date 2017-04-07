#!/usr/bin/python3

import uuid
import os
import sys

from datetime import datetime

from flask import Flask
from flask import request
from flask_cors import CORS

from kai.parser.parser import Parser, JsonSystem


# setup the spacy parser
parser = Parser()


app = Flask(__name__)
CORS(app, resources={r"/parse/*": {"origins": "*"}})


@app.route('/')
def index():
    return "parser service layer"


# curl -H "Content-Type: plain/text" -X POST --data "@test.txt" http://localhost:9000/parse
@app.route('/parse', methods=['POST'])
def parse():
    t1 = datetime.now()
    text = parser.cleanup_text(request.data)
    sentence_list = parser.parse_document(text)

    num_tokens = 0
    for sentence in sentence_list:
        num_tokens += len(sentence)

    delta = datetime.now() - t1
    return JsonSystem().encode({"processing_time": int(delta.total_seconds() * 1000),
                                "sentence_list": sentence_list,
                                "num_tokens": num_tokens,
                                "num_sentences": len(sentence_list)
                                })


# non gunicorn use - debug
if __name__ == "__main__":
    app.run(host="0.0.0.0",
            port=9000,
            debug=True,
            use_reloader=False)
