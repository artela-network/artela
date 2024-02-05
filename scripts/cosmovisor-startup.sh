#!/bin/bash

export DAEMON_NAME=artelad
export DAEMON_HOME=$HOME/.artelad

cosmovisor init ./build/artelad
cosmovisor run start --log_level debug
