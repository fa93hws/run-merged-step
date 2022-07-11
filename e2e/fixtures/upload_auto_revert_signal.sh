#!/usr/bin/env bash
set -eu -o pipefail

echo -n $1 > "${AUTO_REVERT_OUTPUT_PATH}"
