#!/bin/bash
set -e

# TODO: remove this after update users' relays; needed for moving to new version of the script
[ -d ${PWD}/share ] && mv share share_eth


SIDE_NET=$1
MPC_MEID=$2
MPC_ACCESS_TOKEN=$3
MPC_PARTY_IDS=$4
MPC_THRESHOLD=$5
MPC_KEYGEN_URL=$6

SIDE_NET_LOWERCASE="${SIDE_NET,,}"
SIDE_NET_UPPERCASE="${SIDE_NET^^}"

RELAY_CONTAINER_NAME="${SIDE_NET_LOWERCASE}-relay"
KEYGEN_CONTAINER_NAME="${SIDE_NET_LOWERCASE}-relay-keygen"

SHARE_DIR="${PWD}/share_${SIDE_NET_LOWERCASE}"
SHARE_PATH="${SHARE_DIR}/share_${MPC_MEID}"


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

read -sp "${SIDE_NET} network private key: " SIDE_PRIVATE_KEY
while true;
do
    echo -e "\n"
    if [ ${#SIDE_PRIVATE_KEY} -ne 64 ];
        then read -sp 'Key length should be 64 characters, type again: ' SIDE_PRIVATE_KEY;
        else break
    fi
done


echo -e "\nPlease enter your token (issued by the bridge developers)"

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
docker rm -f $RELAY_CONTAINER_NAME
set -e

IMAGE=ghcr.io/ambrosus/ambrosus-bridge
STAGE=${STAGE:-prod}
if [ $STAGE == "prod" ]; then
  TAG=latest
elif [ $STAGE == "test" ]; then
  TAG=dev
elif [ $STAGE == "dev" ]; then
  TAG=dev
fi

docker pull $IMAGE:$TAG

echo "Starting keygen..."
if [ -f $SHARE_PATH ]; then
  echo "Share already exist"
else
  docker run -it --rm \
  --name $KEYGEN_CONTAINER_NAME \
  -v $SHARE_DIR:/app/shared \
  --entrypoint '/bin/sh' \
  $IMAGE:$TAG \
  -c "go run ./cmd/mpc_keygen -url $MPC_KEYGEN_URL -meID $MPC_MEID -partyIDs '$MPC_PARTY_IDS' -threshold $MPC_THRESHOLD -accessToken $MPC_ACCESS_TOKEN -shareDir ./shared"
fi

echo "Starting relay..."
docker run -d \
--name $RELAY_CONTAINER_NAME \
--restart unless-stopped \
-v $SHARE_DIR:/app/shared \
-e STAGE=$STAGE \
-e NETWORK="${SIDE_NET_LOWERCASE}-untrustless" \
-e NETWORKS_AMB_PRIVATEKEY=$AMB_PRIVATE_KEY \
-e "NETWORKS_${SIDE_NET_UPPERCASE}_PRIVATEKEY"=$SIDE_PRIVATE_KEY \
-e EXTERNALLOGGER_TELEGRAM_TOKEN=$EXTERNALLOGGER_TELEGRAM_TOKEN \
-e SUBMITTERS_AMBTOSIDE_MPC_MEID=$MPC_MEID \
-e SUBMITTERS_AMBTOSIDE_MPC_SHAREPATH="shared/share_$MPC_MEID" \
-e SUBMITTERS_AMBTOSIDE_MPC_ACCESSTOKEN=$MPC_ACCESS_TOKEN \
-e SUBMITTERS_SIDETOAMB_MPC_MEID=$MPC_MEID \
-e SUBMITTERS_SIDETOAMB_MPC_SHAREPATH="shared/share_$MPC_MEID" \
-e SUBMITTERS_SIDETOAMB_MPC_ACCESSTOKEN=$MPC_ACCESS_TOKEN \
$IMAGE:$TAG >> /dev/null

sleep 10
docker logs $RELAY_CONTAINER_NAME
