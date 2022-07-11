#!/usr/bin/env bash

if [[ "${DISABLE_UPLOAD_AUTO_REVERT_SIGNAL_FILE:-}" == "true" ]]; then
  exit 0
else
  exit 1
fi