#!/usr/bin/env bash
source .env

NET_ARG=""
if [[ -v NETWORK ]]; then
  NET_ARG="--network ${NETWORK}"
fi

docker exec alice bash -c "lncli ${NET_ARG} listchannels"

docker exec alice bash -c "lncli ${NET_ARG} pendingchannels"