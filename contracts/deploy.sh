#!/bin/bash
set -e

# first arg is side bridge, for example: eth or bsc
# second arg is network type, for example: integr or test or main

SIDE=${1:-eth}
TYPE=${2:-dev}

AMB_NET="$TYPE/amb"
SIDE_NET="$TYPE/$SIDE"
BRIDGE_TAG="bridges_$SIDE"

echo "Deploying on $AMB_NET and $SIDE_NET"

yarn hardhat deploy --network ${AMB_NET} --tags tokens
yarn hardhat deploy --network ${SIDE_NET} --tags tokens

yarn hardhat deploy --network ${AMB_NET} --tags ${BRIDGE_TAG}
yarn hardhat deploy --network ${SIDE_NET} --tags ${BRIDGE_TAG}

yarn hardhat deploy --network ${AMB_NET} --tags ${BRIDGE_TAG}  # setSideBridge to newly eth deployment

yarn hardhat deploy --network ${AMB_NET} --tags tokens_add_bridges
yarn hardhat deploy --network ${SIDE_NET} --tags tokens_add_bridges
