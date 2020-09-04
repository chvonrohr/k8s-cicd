#!/bin/bash

read -p "Are you sure? " -n 1 -r
echo -e "\n"
if [[ ! $REPLY =~ ^[Yy]$ ]]
    rm -i project-start/build/ci/*
    rm -i project-start/build/package/*.Dockerfile
    rm -i project-start/deployments/*.yaml
    rm -i project-start/deployments/*/*.yaml
then
    exit 1
fi
echo -e "\n"