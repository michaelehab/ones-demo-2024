#!/bin/bash

DIR=$PWD
CMD=../../cmd

PATH=/usr/bin:/bin:$PATH

# Kill all logistics-* process
function cleanup {
  pkill ones-demo
}

cd $CMD/transitor
exec -a ones-demo-transitor ./transitor-demo  -cfg ./res/config-mqtt.json &
cd $DIR
sleep 1

cd $CMD/mutator
exec -a ones-demo-mutator ./mutator-demo  -cfg ./res/config-mqtt.json &
cd $DIR
sleep 1

cd $CMD/creator
exec -a ones-demo-creator ./creator-demo  -cfg ./res/config-mqtt.json &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done