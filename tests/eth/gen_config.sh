geth --datadir ./chain init ./genesis.json
#geth --datadir ./chain --password ./account-password.txt account new > ./account.txt
#geth --datadir ./chain makedag 0
geth --datadir ./chain --ethash.dagdir ./chain/ethash \
     --networkid 1234 --nodiscover  \
     --http \
     --mine --miner.threads=1 --miner.etherbase=0xa776765d341784feac83be698Ad6A1b4Ea6fCF82 --miner.gasprice 1 \
     dumpconfig > ./config.toml
