#!/usr/bin/env bash

ME=`basename "$0"`

usage()
{
  echo "Usage: $ME" >&2
  exit 1
}

CURRENT_SHA=`git rev-parse HEAD | cut -c1-6`
LONG_SHA=`git rev-parse HEAD`

deploy_webhooks() {
  helm upgrade kube-webhooks kb/webhooks -f chart-configs/webhooks.yaml \
    --set image.tag=$CURRENT_SHA
}

deploy_webhooks

source scripts/.env

curl https://build.bugsnag.com/ \
    --header "Content-Type: application/json" \
    --data "{
      \"apiKey\": \"$BUGSNAG_API_KEY\",
      \"appVersion\": \"$CURRENT_SHA\",
      \"releaseStage\": \"production\",
      \"sourceControl\": {
        \"provider\": \"github\",
        \"repository\": \"https://github.com/kubebuild/webhooks\",
        \"revision\": \"${LONG_SHA}\"
      }
    }"