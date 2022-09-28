#!/usr/bin/env bash

set -eo pipefail

if [ "${TRAVIS:-false}" == "true" ]; then
  cat > ~/.netrc <<EOF
machine github.ibm.com
    login ${GITHUB_USER}
    password ${GITHUB_TOKEN}
EOF
  chmod 600 ~/.netrc
else
  echo "I do not appear to be running in a Travis environment. Exiting." >&2
  exit 1
fi

echo "Done."
