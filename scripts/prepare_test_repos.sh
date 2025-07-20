#!/usr/bin/env bash

set -e
processed_repos=" "
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
    if ! echo "$processed_repos" | grep -q " $repo_name "; then
        if [ ! -d "$repo_dir" ]; then
            git clone $url $repo_dir
        else
            cd $repo_dir
            git fetch --all
            cd - > /dev/null
        fi
        processed_repos+="$repo_name "
    fi
    cd $repo_dir
    git checkout -q $branch
    git reset -q --hard origin/$branch
    cd - > /dev/null
    submodule_dir="$submodules_dir/$repo_name/$submodule"
    mkdir -p $submodule_dir
    cp -r "$repo_dir/." $submodule_dir
done
