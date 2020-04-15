#!/usr/bin/env bash

set -e
set -o pipefail # Only exit with zero if all commands of the pipeline exit successfully

[[ -z ${DBL_BOT_ID} ]] && echo "DBL_BOT_ID not defined" && exit 1
[[ -z ${DBL_BOT_TOKEN} ]] && echo "DBL_BOT_TOKEN not defined" && exit 1

SCRIPT_PATH=$(readlink -f "${0}")
SCRIPT_DIR=$(dirname "${SCRIPT_PATH}")

COMMIT=$(git rev-parse --short HEAD)

REPO_YMLS="${SCRIPT_DIR}/../deployments/kubernetes"

DEPLOYMENT_YML="${REPO_YMLS}/deployment.yml"
VARIABLIZED_DEPLOYMENT_YML="/tmp/deployment.yml"

setup() {
  cp "${DEPLOYMENT_YML}" "${VARIABLIZED_DEPLOYMENT_YML}"
}

applyValues() {
  sed -i "s|{COMMIT}|${COMMIT}|g" "${VARIABLIZED_DEPLOYMENT_YML}"
  sed -i "s|{DBL_BOT_ID}|${DBL_BOT_ID}|g" "${VARIABLIZED_DEPLOYMENT_YML}"
  sed -i "s|{DBL_BOT_TOKEN}|${DBL_BOT_TOKEN}|g" "${VARIABLIZED_DEPLOYMENT_YML}"
}

deploy() {
  kubectl apply -f "${VARIABLIZED_DEPLOYMENT_YML}"
  kubectl -n ephemeral-roles rollout status --timeout 120s deployment/ephemeral-roles-informer
}

cleanup() {
  rm -f "${VARIABLIZED_DEPLOYMENT_YML}"
}

trap cleanup EXIT

setup
applyValues
deploy
