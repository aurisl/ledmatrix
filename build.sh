#!/bin/bash

if [[ -z "$1" ]]
  then
    echo "Missing render mode"
    exit
fi

#software | hardware
RENDER_MODE=$1

go build -i -tags '"$RENDER_MODE"' .
