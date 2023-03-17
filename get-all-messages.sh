#!/usr/bin/env bash
set -euo pipefail
cd "$(cd "${0%/*}" && pwd)"

_CMD=(
	go run ./ssb-rpc-client/
	-H "$1"
	-P ${2:-8008}
	-c ${3:-local.config.json}
	-s ${4:-local.secret.json}
)

TS=$("${_CMD[@]}" log --keys --limit 1 | jq -c '.timestamp | floor')

NOW="$(date +%s)000"
STEP=1000

echo "since $TS" >&2
echo "until $NOW" >&2
echo " step $STEP" >&2

for ((i=TS;i<NOW;i=i+STEP)); do
	echo "--gt $((i-1)) --lt $((i+STEP))" >&2

	"${_CMD[@]}" log --keys --gt "$((i-1))" --lt "$((i+STEP))"
done