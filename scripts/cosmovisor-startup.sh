#!/bin/bash

export DAEMON_NAME=artelad
export DAEMON_HOME=$HOME/.artelad
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true

cosmovisor init ./build/artelad
cosmovisor run start --log_level debug
