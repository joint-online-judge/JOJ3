#!/usr/bin/env bash

set -e
declare -A repo_names
tmp_dir=${1:-./tmp}
submodules_dir="$tmp_dir/submodules"
rm -rf $submodules_dir
mkdir -p $submodules_dir
submodules=$(git config --file .gitmodules --get-regexp path | awk '{ print $2 }')
for submodule in $submodules; do
    url=$(git config --file .gitmodules --get-regexp "submodule.$submodule.url" | awk '{ print $2 }')
    branch=$(git config --file .gitmodules --get-regexp "submodule.$submodule.branch" | awk '{ print $2 }')
    repo_name=$(echo $url | rev | cut -d'/' -f 1 | rev | cut -d'.' -f 1)
    repo_dir="$tmp_dir/$repo_name"
    if [[ ! -v repo_names["$repo_name"] ]]; then
        if [ ! -d "$repo_dir" ]; then
            git clone $url $repo_dir
        else
            cd $repo_dir
            git fetch --all
            cd -
        fi
    fi
    repo_names[$repo_name]=1
    cd $repo_dir
    git checkout -q $branch
    git reset -q --hard origin/$branch
    cd -
    submodule_dir="$submodules_dir/$repo_name/$submodule"
    mkdir -p $submodule_dir
    cp -rT $repo_dir $submodule_dir
done
