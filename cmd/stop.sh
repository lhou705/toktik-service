#!/usr/bin/env bash
base=$(pwd)
mkdir -p log


if [ "$1" == "message" ] || [ "$1" == "all" ]; then
  bash $(pwd)/cmd/items/message.sh stop
fi

if [ "$1" == "user" ] || [ "$1" == "all" ]; then
  bash $(pwd)/cmd/items/user.sh stop
fi

if [ "$1" == "video" ] || [ "$1" == "all" ]; then
  bash $(pwd)/cmd/items/video.sh stop
fi

if [ "$1" == "web" ] || [ "$1" == "all" ]; then
  bash $(pwd)/cmd/items/web.sh stop
fi
