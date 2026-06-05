#!/bin/sh

echo "argv1=$1 argv2=$2 argv3=$3" >&2

if [ "$1" = "git" ] &&
   [ "$2" = "remote" ] &&
   [ "$3" = "get" ]; then
    echo "get-url"
fi
