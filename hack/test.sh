#!/usr/bin/env bash

# This script runs the unit tests for the repository.
# Usage: `hack/test.sh`.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

echo $REPO_ROOT

export GOBIN="${REPO_ROOT}/hack/bin"
PATH="${GOBIN}:${PATH}"

pushd "${REPO_ROOT}" > /dev/null
  GO111MODULE=on go install github.com/onsi/ginkgo/v2/ginkgo
  ginkgo -r -p --randomize-all --randomize-suites --fail-on-pending --keep-going --race
popd > /dev/null
