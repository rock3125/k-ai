# KAI spaCy parser

A micro-service around the spaCy parser for KAI.

## Required software
See [requirements.txt](./requirements.txt) for this project's main requirements.
To install spaCy's data first install sputnik:
```
$ sudo pip3 install sputnik
```
then from the Python3 REPL or equivalent run:
* NB. this will download ~1.3GB of data
```
import sputnik
import spacy

sputnik.install('spacy', spacy.about.__version__, 'en', data_path='/opt/spacy')
```

## Distribution build and installation

from the root of the repository run:
```
cd spacy/
python3 setup.py bdist_wheel
```
this builds a wheel distribution for this project in the spacy/dist/ folder.

Installation and un-installation are then trivial:
```
sudo pip3 install spacy/dist/kai-parser-1.0-py3-none-any.whl
```

and

```
sudo pip3 uninstall kai-parser
```
