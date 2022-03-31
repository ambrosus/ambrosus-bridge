source .env.relay

#cd ../contracts
#yarn run integr_deploy

cd ../relay
go build ./cmd/bridge
export CONFIG_PATH DEBUG AMB_PRIVATE_KEY ETH_PRIVATE_KEY
./bridge

#cd ../contracts
#yarn run integr_test

