#!/bin/bash

KEY="mykey"
KEY2="mykey2"
KEY3="mykey3"
KEY4="mykey4"

CHAINID="artela_9000-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="debug"
# trace evm
TRACE="--trace"
# TRACE=""

export PATH=./:./build:$PATH

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
    echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
    exit 1
}

# remove existing daemon and client
rm -rf ~/.artelad*

echo ./cmd/artelad/artelad config keyring-backend $KEYRING
./cmd/artelad/artelad config keyring-backend $KEYRING
echo ./cmd/artelad/artelad config chain-id $CHAINID
./cmd/artelad/artelad config chain-id $CHAINID

# if $KEY exists it should be deleted
./cmd/artelad/artelad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO
./cmd/artelad/artelad keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO
./cmd/artelad/artelad keys add $KEY3 --keyring-backend $KEYRING --algo $KEYALGO
./cmd/artelad/artelad keys add $KEY4 --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for artela (Moniker can be anything, chain-id must be an integer)
echo ./cmd/artelad/artelad init $MONIKER --chain-id $CHAINID
./cmd/artelad/artelad init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to aphoton
cat $HOME/.artelad/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aphoton"' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json
cat $HOME/.artelad/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aphoton"' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json
cat $HOME/.artelad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aphoton"' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json
cat $HOME/.artelad/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aphoton"' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json

# Set gas limit in genesis
cat $HOME/.artelad/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="20000000"' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json
cat $HOME/.artelad/config/genesis.json | jq '.app_state["evm"]["params"]["extra_eips"]=[3855]' >$HOME/.artelad/config/tmp_genesis.json && mv $HOME/.artelad/config/tmp_genesis.json $HOME/.artelad/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
./cmd/artelad/artelad add-genesis-account $KEY 100000000000000000000000000aphoton --keyring-backend $KEYRING
./cmd/artelad/artelad add-genesis-account $KEY2 100000000000000000000000000aphoton --keyring-backend $KEYRING
./cmd/artelad/artelad add-genesis-account $KEY3 100000000000000000000000000aphoton --keyring-backend $KEYRING
./cmd/artelad/artelad add-genesis-account $KEY4 100000000000000000000000000aphoton --keyring-backend $KEYRING
echo ./cmd/artelad/artelad add-genesis-account $KEY 100000000000000000000000000aphoton --keyring-backend $KEYRING


# Sign genesis transaction
./cmd/artelad/artelad gentx $KEY 1000000000000000000000aphoton --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
./cmd/artelad/artelad collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
./cmd/artelad/artelad validate-genesis

# disable produce empty block and enable prometheus metrics
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.artelad/config/config.toml
    sed -i '' 's/prometheus = false/prometheus = true/' $HOME/.artelad/config/config.toml
    sed -i '' 's/prometheus-retention-time = 0/prometheus-retention-time  = 1000000000000/g' $HOME/.artelad/config/app.toml
    sed -i '' 's/enabled = false/enabled = true/g' $HOME/.artelad/config/app.toml
else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.artelad/config/config.toml
    sed -i 's/prometheus = false/prometheus = true/' $HOME/.artelad/config/config.toml
    sed -i 's/prometheus-retention-time  = "0"/prometheus-retention-time  = "1000000000000"/g' $HOME/.artelad/config/app.toml
    sed -i 's/enabled = false/enabled = true/g' $HOME/.artelad/config/app.toml
fi

if [[ $1 == "pending" ]]; then
    echo "pending mode is on, please wait for the first block committed."
    if [[ $OSTYPE == "darwin"* ]]; then
        sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.artelad/config/config.toml
        sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.artelad/config/config.toml
    else
        sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.artelad/config/config.toml
        sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.artelad/config/config.toml
    fi
fi
