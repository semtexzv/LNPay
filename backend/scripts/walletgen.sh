#!/bin/bash

set -x
source .env


LN_NET_ARG=""
BT_NET_ARG=""
if [[ -v NETWORK ]]; then
  LN_NET_ARG="--network ${NETWORK}"
  BT_NET_ARG="--${NETWORK}"
fi



ALICE_IDENT=$(docker exec alice bash -c "lncli ${LN_NET_ARG} getinfo" | jq '.identity_pubkey' -r)
ALICE_ADDR=$(docker exec alice bash -c "lncli ${LN_NET_ARG} newaddress p2wkh" | jq '.address' -r)


LND_IDENT=$(docker exec lnd bash -c "lncli ${LN_NET_ARG} getinfo" | jq '.identity_pubkey' -r)
LND_ADDR=$(docker exec lnd bash -c "lncli ${LN_NET_ARG} newaddress p2wkh" | jq '.address' -r)

docker kill btcd2
docker container rm btcd2

sleep 1
docker run -d --name btcd2 --network ${DOCKER_NETWORK} btcd bash -c "MINING_ADDRESS=${ALICE_ADDR} ./start-btcd.sh --addpeer=btcd"
sleep 1
docker exec btcd2 bash -c "./start-btcctl.sh ${BT_NET_ARG} generate 1000"

docker kill btcd2
docker container rm btcd2

sleep 1
docker run -d --name btcd2 --network ${DOCKER_NETWORK} btcd bash -c "MINING_ADDRESS=${LND_ADDR} ./start-btcd.sh --addpeer=btcd"

sleep 1
docker exec btcd2 bash -c "./start-btcctl.sh ${BT_NET_ARG} generate 100"
sleep 5
docker exec lnd bash -c "lncli ${LN_NET_ARG} openchannel --node_key=${ALICE_IDENT} --local_amt=5000000"
sleep 3
docker exec btcd bash -c "./start-btcctl.sh ${BT_NET_ARG} generate 500"
