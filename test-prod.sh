#!/usr/bin/env bash

set -eufo pipefail

# Test language list works
NUM_LANGS=$(curl --silent https://api.iverbs.co.uk/v2/languages | jq '.data | length')
[[ "$NUM_LANGS" == "5" ]] || exit 1

# Get token and download language
TOKEN=$(curl --silent -XPOST https://api.iverbs.co.uk/v2/tokens | jq -r '.token')
[[ "$TOKEN" != "" ]] || exit 1

# Try to download French. Ensure body is over 3 million chars
LEN=$(curl --silent -H "Authorization: $TOKEN" https://api.iverbs.co.uk/v2/languages/fr | wc -c)
LEN="${LEN//[$'\t\r\n ']}"
[ "$LEN" -gt 3000000 ] || exit 1
