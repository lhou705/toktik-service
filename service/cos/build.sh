#!/usr/bin/env bash
RUN_NAME="cos"
mkdir -p "output/bin"
cp "./script/bootstrap.sh" "output/"
chmod +x "./output/bootstrap.sh"

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o "output/bin/${RUN_NAME}"
else
    go test -c -covermode=set -o "output/bin/${RUN_NAME}" -coverpkg="./..."
fi
