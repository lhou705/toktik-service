#!/usr/bin/env bash
base=$(pwd)
mkdir -p log
if [ "$1" == "cos" ] || [ "$1" == "all" ]; then
  nohup ./bin/cos -config=$base/config/cos.config.json >> log/cos.log &
fi

if [ "$1" == "message" ] || [ "$1" == "all" ]; then
  nohup ./bin/message -config=$base/config/message.config.json >> log/message.log &
fi

if [ "$1" == "user" ] || [ "$1" == "all" ]; then
  nohup ./bin/user -config=$base/config/user.config.json >> log/user.log &
fi

if [ "$1" == "video" ] || [ "$1" == "all" ]; then
  nohup ./bin/video -config=$base/config/video.config.json >> log/video.log &
fi

if [ "$1" == "web" ] || [ "$1" == "all" ]; then
  nohup ./bin/web -config=$base/config/web.config.json >> log/web.log &
fi