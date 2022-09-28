#!/usr/bin/env bash
# write-netrc.sh
#
# The purpose of this script is to generate the .netrc file in the Travis
# environment because it is required for "go get" to work with private modules.
#
# It uses the GITHUB_USER and GITHUB_TOKEN environment variables, so they need
# to be configured correctly in the build for this script to work properly.

set -eo pipefail

# We only want to run this script if it is running in Travis
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
