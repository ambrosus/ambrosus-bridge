rm -rf ./amb/apollo0/chains
rm -rf ./amb/apollo1/chains
rm -rf ./amb/apollo2/chains

rm -rf ./eth/chain

docker-compose rm -f
docker volume prune -f
docker-compose up