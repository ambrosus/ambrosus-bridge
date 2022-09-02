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

read -sp 'Ambrosus private key: ' amb_private_key
while true;
do
    echo -e "\n"
    if [ ${#amb_private_key} -ne 64 ];
        then read -sp 'Key length should be 64 characters, type again: ' amb_private_key ;
        else break
    fi
done

set +e
docker rm -f eth-relay
set -e

echo "Starting relay..."
docker run -d \
--name eth-relay \
--restart unless-stopped \
-e STAGE=prod \
-e NETWORK=eth-untrustless \
-e NETWORKS_AMB_PRIVATEKEY=$amb_private_key \
ghcr.io/ambrosus/ambrosus-bridge >> /dev/null

sleep 10
docker logs eth-relay
