#!/bin/bash

export DAEMON_NAME=artelad
export DAEMON_HOME=$HOME/.artelad
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=false
export DAEMON_DOWNLOAD_MUST_HAVE_CHECKSUM=false
#export DAEMON_DATA_BACKUP_DIR=$HOME/.artelad_bak
export DAEMON_PREUPGRADE_MAX_RETRIES=3

DIR_NAME=$1
HEIGHT=$2

cosmovisor add-upgrade $DIR_NAME ./build/artelad

sleep 6s

./build/artelad tx gov submit-legacy-proposal software-upgrade $DIR_NAME \
    --node "tcp://localhost:26657" \
    --title "upgrade $DIR_NAME" \
    --description "First step - upgrade test" \
    --upgrade-height $HEIGHT \
    --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/artela-network/artela/releases/download/v0.4.1-beta/artelad_update_Linux_amd64.zip?checksum=sha256:15eff4e0767ddad4cc20fd8751ff87ffd22c72d1c8d355c915e426156a409e17"}}' \
    --deposit 30000000aart \
    --from mykey2 \
    --no-validate \
    --yes

./build/artelad tx gov deposit 1 30000000aart --from mykey --yes
./build/artelad tx gov vote 1 yes --from mykey --yes

./build/artelad  query auth module-account gov



#./build/artelad query tx 00BE99132CE0C2E741842A342EE63DE8353EB9A6327B280995599BA05CA00FDF
#./build/artelad query gov params
#./build/artelad query gov proposal 1
#./build/artelad query upgrade plan

./artelad tx gov submit-legacy-proposal software-upgrade v047rc7 \
    --node "tcp://localhost:26657" \
    --title "upgrade v047rc7" \
    --description "upgrade v047rc7" \
    --upgrade-height 2448458 \
    --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/artela-network/artela/releases/download/v0.4.1-beta/artelad_update_Linux_amd64.zip?checksum=sha256:15eff4e0767ddad4cc20fd8751ff87ffd22c72d1c8d355c915e426156a409e17"}}' \
    --deposit 30000000aart \
    --from node0 \
    --yes


./build/artelad tx gov submit-proposal draft_proposal.json --from mykey --keyring-backend test --chain-id=artela_11820-1
#./build/artelad query gov proposal 1
#./build/artelad query upgrade plan
./build/artelad tx gov deposit 1 30000000aart --from mykey --yes
./build/artelad tx gov vote 1 yes --from mykey --yes

./build/artelad  query auth module-account gov

./build/artelad  query gov tally 1