#!/bin/bash
git remote add upstream git@gitlab.com:letsboot/core/kubernetes-course.git
git checkout master ||exit 1
git fetch upstream  ||exit 1
git pull upstream master ||exit 1
git push origin master ||exit 1