#!/bin/bash

KEY="dev0"
CHAINID="artela_11820-1"
MONIKER="mymoniker"

echo "create and add new keys"
./artelad keys add $KEY --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init artela with moniker=$MONIKER and chain-id=$CHAINID"
./artelad init $MONIKER --chain-id $CHAINID
echo "prepare genesis: Allocate genesis accounts"
./artelad add-genesis-account \
"$(./artelad keys show $KEY -a --keyring-backend test)" 1000000000000000000aart,1000000000000000000stake \
--keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./artelad gentx $KEY 1000000000000000000stake --keyring-backend test --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./artelad collect-gentxs
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./artelad validate-genesis

echo "starting artela node $i in background ..."
./artelad start --pruning=nothing --rpc.unsafe \
--keyring-backend test \
>$DATA_DIR/node.log 2>&1 & disown

echo "started artela node"
tail -f /dev/null