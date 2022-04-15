#!/bin/bash
set -e

# first arg is side bridge, for example: eth or bsc
# second arg is network type, for example: integr or test or main

SIDE=${1:-eth}
TYPE=${2:-test}

AMB_NET="amb_$TYPE"
SIDE_NET="${SIDE}_$TYPE"

echo "Deploying on $AMB_NET and $SIDE_NET"

yarn hardhat deploy --network ${AMB_NET} --tags tokens
yarn hardhat deploy --network ${SIDE_NET} --tags tokens

yarn hardhat deploy --network ${AMB_NET} --tags bridges
yarn hardhat deploy --network ${SIDE_NET} --tags bridges

yarn hardhat deploy --network ${AMB_NET} --tags bridges  # setSideBridge to newly eth deployment

yarn hardhat deploy --network ${AMB_NET} --tags tokens_add_bridges
yarn hardhat deploy --network ${SIDE_NET} --tags tokens_add_bridges
