#!/bin/bash

killall -9 CavaSvrProxy

APP_LOG_DIR=./log
if [ ! -d $APP_LOG_DIR ];then
    mkdir -p $APP_LOG_DIR
fi

mkdir -p cava3/upload
chmod -R 755 cava3/upload
mkdir -p cava3/outputs/log

./CavaSvrProxy > $APP_LOG_DIR/run.log 2>&1 &

pidof CavaSvrProxy > server.pid
echo "CavaSvrProxy start succeed."
