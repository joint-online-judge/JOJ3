#!/usr/bin/env bash

set -ex
tmp_dir=${1:-./tmp}
JOJ3=$(git rev-parse --show-toplevel)/build/joj3
command=${2:-$JOJ3}
submodules_dir="$tmp_dir/submodules"
submodules=$(git config --file .gitmodules --get-regexp path | awk '{ print $2 }')
for submodule in $submodules; do
    url=$(git config --file .gitmodules --get-regexp "submodule.$submodule.url" | awk '{ print $2 }')
    repo_name=$(echo $url | rev | cut -d'/' -f 1 | rev | cut -d'.' -f 1)
    submodule_dir="$submodules_dir/$repo_name/$submodule"
    cd $submodule_dir
    eval "$command"
    if [[ $command == $JOJ3 ]]; then
        if [ -f "./expected.json" ]; then
            mv -f "joj3_result.json" "expected.json"
        fi
    fi
    cd -
done
