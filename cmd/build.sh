#!/usr/bin/env bash
base=$(pwd)
mkdir -p bin
if [ "$1" == "cos" ] || [ "$1" == "all" ]; then
  echo "进入cos目录"
  cd service/cos
  echo "编译cos模块"
  sh build.sh
  cp output/bin/cos $base/bin/toktik_cos
  echo "已完成cos模块的编译"
  cd $base
fi

if [ "$1" == "user" ] || [ "$1" == "all" ]; then
  echo "进入user目录"
  cd service/user
  echo "编译user模块"
  sh build.sh
  cp output/bin/user $base/bin/toktik_user
  echo "已完成user模块的编译"
  cd $base
fi

if [ "$1" == "message" ] || [ "$1" == "all" ]; then
  echo "进入message目录"
  cd service/message
  echo "编译message模块"
  sh build.sh
  cp output/bin/message $base/bin/toktik_message
  echo "已完成message模块的编译"
  cd $base
fi

if [ "$1" == "video" ] || [ "$1" == "all" ]; then
  echo "进入video目录"
  cd service/video
  echo "编译video模块"
  sh build.sh
  cp output/bin/video $base/bin/toktik_video
  echo "已完成video模块的编译"
  cd $base
fi

if [ "$1" == "web" ] || [ "$1" == "all" ]; then
  echo "进入web目录"
  cd service/web
  echo "编译web模块"
  go build -o web
  mv web $base/bin/toktik_web
  echo "已完成web模块的编译"
  cd $base
fi

