#!/bin/bash

DATA_DIR="$HOME/.artelad"

sed -i 's/127.0.0.1:8545/0.0.0.0:8545/g' $DATA_DIR/config/app.toml
sed -i 's/timeout_commit = \"5s\"/timeout_commit = \"500ms\"/g' $DATA_DIR/config/config.toml
sed -i 's/"extra_eips": \[\]/"extra_eips": \[3855\]/g' $DATA_DIR/config/genesis.json

echo "starting artela node $i in background ..."
./artelad start \
--log_level debug \
--minimum-gas-prices=0.0001art \
--api.enable \
--json-rpc.api eth,txpool,personal,net,debug,web3,miner \
--api.enable \
&>>$DATA_DIR/node.log & disown

pid=$(ps ax | grep '[a]rtelad' | awk '{print $1}')

nohup dlv attach $pid --listen=:19211 --headless=true --log=true  \
  --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --accept-multiclient --api-version=2 --continue &> /dev/null &

echo "started artela node"
tail -f /dev/null