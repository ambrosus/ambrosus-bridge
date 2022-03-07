FROM ethereum/client-go:v1.10.13

ARG ACCOUNT_PASSWORD

COPY geth-genesis.json /tmp

RUN geth init /tmp/geth-genesis.json \
    && rm -f ~/.ethereum/geth/nodekey \
    && echo ${ACCOUNT_PASSWORD} > /tmp/password \
    && geth account new --password /tmp/password \
    && geth account new --password /tmp/password \
    && geth account new --password /tmp/password \
    && LOCAL_ADDRESS=$(geth account new --password /tmp/password | grep -Eo "0x[0-9a-zA-Z]+$") \
    && echo LOCAL_ADDRESS=$LOCAL_ADDRESS \
    #&& sed -i "s/LOCAL_ADDRESS=.*/LOCAL_ADDRESS=${LOCAL_ADDRESS}/" /.env \
    && rm -f /tmp/password

ENTRYPOINT ["geth"]
