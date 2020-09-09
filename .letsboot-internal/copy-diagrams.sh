#!/bin/bash

if [ ! -d "./assets" ]; then \
    echo "run in course folder" 1>&2
    exit 1
fi 

cp -v /Volumes/GoogleDrive/Shared\ drives/Letsboot/kubernetes/training-diagrams/* ./assets/