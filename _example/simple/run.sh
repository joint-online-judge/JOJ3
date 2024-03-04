#!/usr/bin/env bash

set -xe
DIRNAME=`dirname -- "$0"`
# cd to make CopyInCwd work
cd $DIRNAME
./../../build/joj3
cd -
