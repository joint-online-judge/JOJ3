#!/usr/bin/env bash

set -xe
DIRNAME=`dirname -- "$0"`
# cd to make CopyInCwd work
cd $DIRNAME
./../../build/joj3
cat ./joj3_result.json
rm -f ./joj3_result.json
cd -
