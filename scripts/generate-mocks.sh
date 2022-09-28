#!/usr/bin/env bash
# generate-mocks.sh
# Generate the mock test objects for the go testing.

set -o pipefail

if [ -z "${GOROOT}" ]; then
  GOROOT=${HOME}/go
  PATH=${GOROOT}/bin:$PATH
fi

if [ "$#" -ne 1 ]; then
  echo "must supply a base directory in which to create mocks" >&2
  exit 1
fi

MOCKERY_FOUND=$(command -v mockery)

if [ ! -x "${MOCKERY_FOUND-}" ]; then
  go install github.com/vektra/mockery/v2@latest
fi

if [ ! -d "${1}" ]; then
  echo "expected a directory name for argument: \"${1}\"" >&2
  exit 1
fi

(cd $1 && mockery --dir pkg --recursive --all --keeptree)

exit $?