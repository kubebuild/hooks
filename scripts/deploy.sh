#!/usr/bin/env bash

ME=`basename "$0"`

usage()
{
  echo "Usage: $ME" >&2
  exit 1
}

CURRENT_SHA=`git rev-parse HEAD | cut -c1-6`

deploy_webhooks() {
  helm upgrade kube-webhooks kb/webhooks -f chart-configs/webhooks.yaml \
    --set image.tag=$CURRENT_SHA
}

deploy_webhooks
