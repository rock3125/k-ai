#!/bin/bash

if [ $# -ne 1 ]
  then
    echo "takes one parameter: version"
    echo "the version to release / build"
    exit 0
fi

RELEASE=$1

# go test k-ai/...
# process if all passed
if [ $? -eq 0 ]; then

# build go system (no unit testing)
go install k-ai/kai

# only proceed on success
if [ $? -eq 0 ]; then

  # set the version number using git
  # VERSION=`git describe --all --long | cut -d "/" -f 2 | cut -d "-" -f 1`
  sed -i -e 's/^Version.*$/Version = "'"$RELEASE"'"/g' data/properties.ini

  # commit changes
  git commit -am "release $RELEASE"
  git push

  # tag
  git tag "$RELEASE" master
  git push --tags
fi

fi

