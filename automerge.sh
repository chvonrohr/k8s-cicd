#!/bin/bash

# initial merge
./merge.sh

# check for changes in .md files
while true
  do sleep 2
  if [ $(find . -maxdepth 2 -name '*.md' -mtime -4s ! -name '*all*md' | wc -l) -ge 1 ]
    then ./merge.sh
    echo merged
  fi
done
