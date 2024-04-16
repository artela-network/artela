#!/bin/bash

DATA_DIR="$HOME/.artelad"

sed -i 's/127.0.0.1:8545/0.0.0.0:8545/g' $DATA_DIR/config/app.toml
sed -i 's/timeout_commit = \"5s\"/timeout_commit = \"1s\"/g' $DATA_DIR/config/config.toml
sed -i 's/"extra_eips": \[\]/"extra_eips": \[3855\]/g' $DATA_DIR/config/genesis.json

echo "starting artela node $i in background ..."
./artelad start \
--minimum-gas-prices=0.0001art \
--api.enable \
--json-rpc.api eth,txpool,personal,net,debug,web3,miner \
--api.enable \
&>>$DATA_DIR/node.log & disown

echo "started artela node"
tail -f /dev/null