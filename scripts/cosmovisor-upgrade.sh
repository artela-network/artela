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

#mkdir -p $DAEMON_HOME/cosmovisor/upgrades/$DIR_NAME/bin
#cp ./build/$DAEMON_NAME $DAEMON_HOME/cosmovisor/upgrades/$DIR_NAME/bin

cosmovisor add-upgrade $DIR_NAME ./build/artelad

#./build/artelad tx gov submit-legacy-proposal software-upgrade $DIR_NAME --title 'upgrade' --description 'upgrade' --upgrade-height $HEIGHT --deposit 100art --no-validate --from mykey2 --node tcp://localhost:26657 --yes
#./build/artelad tx gov deposit 1 200art --from mykey2 --yes
#./build/artelad tx gov vote 1 yes --from mykey2 --yes

#./build/artelad tx gov deposit 1 200art --from mykey --yes
#./build/artelad tx gov vote 1 yes --from mykey --yes




sleep 6s

./build/artelad tx gov submit-proposal ./draft_proposal.json --from mykey --keyring-backend test -y

sleep 6s

./build/artelad tx gov deposit 1 100art --from mykey --yes

sleep 6s

./build/artelad tx gov vote 1 yes --from mykey --yes

sleep 6s

./build/artelad query gov proposal 1