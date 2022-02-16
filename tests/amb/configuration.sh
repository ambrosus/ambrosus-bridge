#!/bin/bash

BASE_PATH=$2
NETWORK_ID=$(printf '0x%x' $3)
NETWORK_NAME=$4
FILEPATH_GENESIS=${BASE_PATH}/config/genesis.json
FILEPATH_ACCOUNT_PASSWORD=${BASE_PATH}/config/account.password

if [ "$1" == "yes" ]; then \

    echo "Deleting data dir..." && \
    rm -rf ${BASE_PATH}/data/*  && \
    rm -rf ${BASE_PATH}/config/account.*  && \
    rm -f ${BASE_PATH}/config/config.toml && \
    rm -f ${BASE_PATH}/config/genesis.json && \

    echo ${ACCOUNT_PASSWORD} > ${BASE_PATH}/config/account.password  && \

    ${BASE_PATH}/openethereum --keys-path=${BASE_PATH}/data/keys --base-path=${BASE_PATH}/data account new --password ${BASE_PATH}/config/account.password > ${BASE_PATH}/config/account.address && \

    ACCOUNT_ADDRESS=$(cat ${BASE_PATH}/config/account.address) && \
    sed -e "s|0x0000000000000000000000000000000000000005|"${ACCOUNT_ADDRESS}"|g" -e "s|NETWORK_ID|"${NETWORK_ID}"|g" -e "s|NETWORK_NAME|"${NETWORK_NAME}"|g" ${BASE_PATH}/config/genesis.json.example > ${BASE_PATH}/config/genesis.json && \
    sed -e "s|0x0000000000000000000000000000000000000005|"${ACCOUNT_ADDRESS}"|g" -e "s|FILEPATH_GENESIS|${FILEPATH_GENESIS}|g" -e "s|FILEPATH_ACCOUNT_PASSWORD|${FILEPATH_ACCOUNT_PASSWORD}|g" ${BASE_PATH}/config/config.toml.example  >> ${BASE_PATH}/config/config.toml && \

    cat ${BASE_PATH}/data/keys/ethereum/"$(ls -1rt ${BASE_PATH}/data/keys/ethereum | tail -n1)" > ${BASE_PATH}/config/account.private.json && \
    rm -rf ${BASE_PATH}/data/keys/ethereum;
fi

${BASE_PATH}/openethereum \
    --config ${BASE_PATH}/config/config.toml \
    --base-path=${BASE_PATH}/data account import ${BASE_PATH}/config/account.private.json;
