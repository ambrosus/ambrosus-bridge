# yeap, using Makefile to call scripts from package.json

.PHONY: deploy_all
deploy_all:
	yarn hardhat deploy --network amb_test --tags tokens
	yarn hardhat deploy --network eth_test --tags tokens

	yarn hardhat deploy --network amb_test --tags bridges
	yarn hardhat deploy --network eth_test --tags bridges

	yarn hardhat deploy --network amb_test --tags after_deploy
	yarn hardhat deploy --network eth_test --tags after_deploy
