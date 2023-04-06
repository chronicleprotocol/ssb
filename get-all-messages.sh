#!/usr/bin/env bash
#
# SSB Tools
#     Copyright (C) 2023 Chronicle Labs, Inc.
#
#     This program is free software: you can redistribute it and/or modify
#     it under the terms of the GNU Affero General Public License as published
#     by the Free Software Foundation, either version 3 of the License, or
#     (at your option) any later version.
#
#     This program is distributed in the hope that it will be useful,
#     but WITHOUT ANY WARRANTY; without even the implied warranty of
#     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#     GNU Affero General Public License for more details.
#
#     You should have received a copy of the GNU Affero General Public License
#     along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

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
