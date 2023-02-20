#!/usr/bin/env bash
base=$(pwd)
mkdir -p bin
if [ "$1" == "cos" ] || [ "$1" == "all" ]; then
  rm bin/cos
  echo "进入cos目录"
  cd service/cos
  echo "删除cos编译文件"
  rm -rf output
  echo "已完成cos模块的清理"
  cd $base
fi

if [ "$1" == "message" ] || [ "$1" == "all" ]; then
  rm bin/message
  echo "进入message目录"
  cd service/message
  echo "删除message编译文件"
  rm -rf output
  echo "已完成message模块的清理"
  cd $base
fi

if [ "$1" == "user" ] || [ "$1" == "all" ]; then
  rm bin/user
  echo "进入user目录"
  cd service/user
  echo "删除user编译文件"
  rm -rf output
  echo "已完成user模块的清理"
  cd $base
fi

if [ "$1" == "video" ] || [ "$1" == "all" ]; then
  rm bin/video
  echo "进入video目录"
  cd service/video
  echo "删除video编译文件"
  rm -rf output
  echo "已完成video模块的清理"
  cd $base
fi

if [ "$1" == "web" ] || [ "$1" == "all" ]; then
  rm bin/web
  echo "已完成web模块的清理"
  cd $base
fi