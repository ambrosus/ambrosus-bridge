OLD_PARTY="master Kevin Lang Rory Andrey"
NEW_PARTY="master2 Andrey2 backup2 Michael2 Seth2 Od2"

SIDE_NET=$1
MPC_MEID=$2

SIDE_NET_LOWERCASE="${SIDE_NET,,}"
KEYGEN_CONTAINER_NAME="${SIDE_NET_LOWERCASE}-relay-keygen"
SHARE_DIR="${PWD}/share_${SIDE_NET_LOWERCASE}"


if [[ "$SIDE_NET_LOWERCASE" == "eth" ]]; then
  URL="ws://65.108.51.81:6555"   # eth
elif [[ "$SIDE_NET_LOWERCASE" == "bsc" ]]; then
  URL="ws://65.108.58.138:6555"  # bsc
else
  echo "Error: SIDE_NET_LOWERCASE value does not match any condition"
  exit 1
fi


if [[ "$MPC_MEID" == "Kevin" || "$MPC_MEID" == "Lang"  || "$MPC_MEID" == "Rory" ]]; then
  CMD="go run cmd/main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID $MPC_MEID"
elif [[ "$MPC_MEID" == "backup2" || "$MPC_MEID" == "Michael2" || "$MPC_MEID" == "Seth2" || "$MPC_MEID" == "Od2" ]]; then
  CMD="go run cmd/main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meIDNew $MPC_MEID"
elif [[ "$MPC_MEID" == "Andrey" ]]; then
  CMD="go run cmd/main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID Andrey -meIDNew Andrey2"
else
    echo "Error: MPC_MEID value does not match any condition"
    exit 1
fi




if [[ -z $NO_DOCKER_INSTALL ]]; then
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
fi


IMAGE=ghcr.io/ambrosus/ambrosus-bridge
docker pull $IMAGE:latest

echo "Starting reshare..."

docker run -it --rm \
  --name $KEYGEN_CONTAINER_NAME \
  -v $SHARE_DIR:/app/shared \
  --entrypoint '/bin/sh' \
  $IMAGE:latest \
  -c "$CMD"



## run on master
#go run main.go -reshare -server :6555 -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir "shared/share_master" -meID master -meIDNew master2
#
## run on old committee members
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID Kevin
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID Lang
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID Rory
#
## run on new committee members
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meIDNew backup2
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meIDNew Michael2
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meIDNew Seth2
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meIDNew Od2
#
## run on both old and new committee members
#go run main.go -reshare -url $URL -partyIDs $OLD_PARTY -threshold 5 -partyIDsNew $NEW_PARTY -thresholdNew 5 -shareDir $SHARE_DIR -meID Andrey -meIDNew Andrey2
