#!/bin/bash
if [ $# -eq 1 ]
  then omxplayer $1
elif [ $# -eq 2 ]
  then omxplayer $(youtube-dl -g -f best $1)
fi
