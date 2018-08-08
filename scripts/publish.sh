#!/usr/bin/env bash

ME=`basename "$0"`

usage()
{
  echo "Usage: $ME" >&2
  exit 1
}

CURRENT_SHA=`git rev-parse HEAD | cut -c1-6`

IMAGE=kubebuild/webhooks

docker build -t $IMAGE:$CURRENT_SHA .
docker tag $IMAGE:$CURRENT_SHA $IMAGE:latest

docker push $IMAGE:$CURRENT_SHA
docker push $IMAGE:latest