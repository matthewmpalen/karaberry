#!/bin/sh
if [ -n "$1" ]
  then omxplayer $(youtube-dl -g -f best $1)
fi
