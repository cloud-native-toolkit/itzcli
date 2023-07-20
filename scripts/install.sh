#!/usr/bin/env bash

set -ou pipefail

ITZ_INSTALL_HOME=${ITZ_INSTALL_HOME:-/usr/local}
ITZ_INSTALL_BIN_DIR=${ITZ_INSTALL_HOME}/bin
ITZ_INSTALL_VER=${ITZ_INSTALL_VER:-v0.1.28}

echo "Installing itz to ${ITZ_INSTALL_BIN_DIR}..."

assert_installed() {
  BIN=$(command -v "$1")
  if [[ -z "${BIN:-}" ]]; then
    echo "Could not find \"${1}\" on the path; make sure that \"${1}\" is installed." >&2
    exit 1
  fi
}

# Check to make sure that curl is installed...
assert_installed curl
assert_installed tar

# Download the binary to the /tmp folder for the OS
INSTALL_OS=$(uname -s)
if [[ "${INSTALL_OS:-none}" == "Linux" ]]; then
  ITZ_RELEASE_URL=https://github.com/cloud-native-toolkit/itzcli/releases/download/${ITZ_INSTALL_VER}/itzcli-linux-amd64.tar.gz
elif [[ "${INSTALL_OS:-none}" == "Darwin" ]]; then
  ITZ_RELEASE_URL=https://github.com/cloud-native-toolkit/itzcli/releases/download/${ITZ_INSTALL_VER}/itzcli-darwin-amd64.tar.gz
else
    echo "${INSTALL_OS} is currently not supported for installing ITZ with this script."
    exit 1
fi

curl -sS -L -o /tmp/itzcli.tar.gz "${ITZ_RELEASE_URL}"

# Check to make sure the itz bin dir is on the path, and if it's not, tell the
# user to all it to the path in their favorite shell.
if [[ ! -d "${ITZ_INSTALL_BIN_DIR}" ]]; then
  echo "Install bin directory \"${ITZ_INSTALL_BIN_DIR}\" does not exist; creating..." >&2
  sudo mkdir -p "${ITZ_INSTALL_BIN_DIR}"
fi

(cd /tmp && tar xzf /tmp/itzcli.tar.gz)
sudo mv /tmp/itzcli "${ITZ_INSTALL_BIN_DIR}"
sudo mv /tmp/itz "${ITZ_INSTALL_BIN_DIR}"

ON_PATH=$(echo "${PATH}" | grep -c "${ITZ_INSTALL_BIN_DIR}")
if [[ ${ON_PATH} -eq 0 ]]; then
  echo "The directory \"${ITZ_INSTALL_BIN_DIR}\" is not found on your path." >&2
  echo "Make sure to add it to your path before running \"itz\", for example:" >&2
  echo "    PATH=${ITZ_INSTALL_BIN_DIR}:\${PATH}" >&2
fi

echo "Install successful!" >&2
echo -n "Version: "
"${ITZ_INSTALL_BIN_DIR}"/itz version
exit $?
