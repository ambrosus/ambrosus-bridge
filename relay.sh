#!/bin/bash
set -e

apt-get remove -y docker docker-engine docker.io containerd runc && \
apt-get update && \
apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

mkdir -p /etc/apt/keyrings

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

echo "Please enter your private key"

read -sp 'Ambrosus private key: ' AMB_PRIVATE_KEY
while true;
do
    echo -e "\n"
    if [ ${#AMB_PRIVATE_KEY} -ne 64 ];
        then read -sp 'Key length should be 64 characters, type again: ' AMB_PRIVATE_KEY;
        else break
    fi
done


echo "Please enter your token (issued by the bridge developers)"

read -sp 'Token: ' EXTERNALLOGGER_TELEGRAM_TOKEN
while true;
do
    echo -e "\n"
    if [ ${#EXTERNALLOGGER_TELEGRAM_TOKEN} -ne 46 ];
        then read -sp 'Token length should be 46 characters, type again: ' EXTERNALLOGGER_TELEGRAM_TOKEN;
        else break
    fi
done



set +e
docker rm -f eth-relay
set -e

IMAGE=ghcr.io/ambrosus/ambrosus-bridge
TAG=latest
STAGE=prod

docker pull $IMAGE:$TAG

echo "Starting relay..."
docker run -d \
--name eth-relay \
--restart unless-stopped \
-e STAGE=$STAGE \
-e NETWORK=eth-untrustless \
-e NETWORKS_AMB_PRIVATEKEY=$AMB_PRIVATE_KEY \
-e EXTERNALLOGGER_TELEGRAM_TOKEN=$EXTERNALLOGGER_TELEGRAM_TOKEN \
$IMAGE:$TAG >> /dev/null

sleep 10
docker logs eth-relay
