#!/bin/bash

# little bit rough - refine this later
killall gunicorn

if [ -f /opt/kai/spacy/server.py ]; then
    cd /opt/kai/spacy/
fi

# start with large timeout allowing for large texts to be parsed
gunicorn --bind 0.0.0.0:9000 --timeout 120 --threads 1 server:app &

exit 0

