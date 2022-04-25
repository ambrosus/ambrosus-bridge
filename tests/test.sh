source .env.relay

#cd ../contracts
#yarn run integr_deploy

cd ../relay
go build ./cmd/bridge
export CONFIG_PATH=configs/integr
./bridge

#cd ../contracts
#yarn run integr_test

