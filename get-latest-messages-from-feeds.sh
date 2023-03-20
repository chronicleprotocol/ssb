#!/usr/bin/env bash
# pass the feed ids to stdin

set -xeuo pipefail
cd "$(cd "${0%/*}" && pwd)"

_CMD=(
	go run ./ssb-rpc-client/
	-H "$1"
	-P ${2:-8008}
	-c ${3:-local.config.json}
	-s ${4:-local.secret.json}
)

while read -r _line; do
	"${_CMD[@]}" user --keys --reverse --limit 1 "$_line"
done

# | jq -c '{key}*.value|{key,author,sequence}*.content|{key,author,sequence,type,version}'