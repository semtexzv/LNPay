#!/usr/bin/env bash
set -x
source .env
docker-compose up --build -d
sleep 10

LN_NET_ARG=""
BT_NET_ARG=""
if [[ -v NETWORK ]]; then
  LN_NET_ARG="--network ${NETWORK}"
  BT_NET_ARG="--${NETWORK}"
fi



docker cp lnd:/root/.lnd/data/chain/bitcoin/${NETWORK}/admin.macaroon ${LND_ADMIN_MACAROON}
docker cp lnd:/root/.lnd/tls.cert ./${LND_TLS_CERT}

ALICE_IDENT=$(docker exec alice bash -c "lncli ${LN_NET_ARG} getinfo" | jq '.identity_pubkey' -r)
ALICE_IP=$(docker inspect alice | jq -r ".[0].NetworkSettings.Networks.${DOCKER_NETWORK}.IPAddress")
docker exec lnd bash -c "lncli ${LN_NET_ARG} connect ${ALICE_IDENT}@${ALICE_IP}"


LND_IDENT=$(docker exec lnd bash -c "lncli ${LN_NET_ARG} getinfo" | jq '.identity_pubkey' -r)
LND_IP=$(docker inspect lnd | jq -r ".[0].NetworkSettings.Networks.${DOCKER_NETWORK}.IPAddress")
docker exec alice bash -c "lncli ${LN_NET_ARG} connect ${LND_IDENT}@${LND_IP}"


docker exec btcd bash -c "./start-btcctl.sh ${BT_NET_ARG} generate 500"