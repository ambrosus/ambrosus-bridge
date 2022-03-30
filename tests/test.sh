source .env.relay

#cd ../contracts
#yarn run integr_deploy

cd ../relay
go build ./cmd/bridge
export AMB_PRIVATE_KEY ETH_PRIVATE_KEY
./cmd/bridge/bridge

#cd ../contracts
#yarn run integr_test

