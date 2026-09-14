#!/usr/bin/env bash
# Print the core.get_methods entry for one TrueNAS API method.
#
# Read-only. Runs `midclt call core.get_methods` on the box over SSH and filters locally.
#
#   TRUENAS_SSH  ssh destination, e.g. "truenas" or "admin@nas.lan" (required)
set -euo pipefail

method="${1:?usage: api-method.sh <method>}"
: "${TRUENAS_SSH:?TRUENAS_SSH is required (ssh destination)}"

ssh "$TRUENAS_SSH" midclt call core.get_methods |
	jq --exit-status --arg m "$method" '.[$m] // error("unknown method: \($m)")'
