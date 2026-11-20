#!/bin/sh
# Records real responses for a new test: runs ./mz through a recording proxy,
# saves each response into e2e/testdata/mz/<name>/ with a routes.txt of
# Route(...) lines to paste into a test, and scrubs personal values listed in
# e2e/scrub.local.
#   e2e/record.sh <name> <mz args...>
# Uses your real credentials from the environment or ./.env. Only GETs and the
# login are forwarded, so nothing changes on the site.
set -eu
cd "$(dirname "$0")/.."

[ $# -ge 2 ] || { echo "usage: e2e/record.sh <name> <mz args...>" >&2; exit 2; }
name=$1
shift
dir=e2e/testdata/mz/$name
port=${RECORD_PORT:-8090}
if [ -e "$dir" ]; then echo "$dir exists; delete it to re-record" >&2; exit 1; fi

e2e/bin/fakemz record -dir "$dir" -ref "$name" -addr "127.0.0.1:$port" &
pid=$!
trap 'kill $pid 2>/dev/null' EXIT
i=0
until curl -sf -o /dev/null "http://127.0.0.1:$port/_/health"; do
	i=$((i + 1))
	[ $i -lt 50 ] || { echo "recorder did not start" >&2; exit 1; }
	sleep 0.1
done

now=$(date -Iseconds)
MZ_URL="http://127.0.0.1:$port/" ./mz "$@" || echo "(mz exited with $?)"

if [ -f e2e/scrub.local ]; then
	e2e/bin/fakemz scrub -map e2e/scrub.local "$dir"
else
	echo "warning: no e2e/scrub.local; $dir contains your account data" >&2
fi
echo
echo "recorded $dir at $now (set mz.Env[\"MZ_NOW\"] if the command depends on the date)."
echo "Routes to paste into a test ($dir/routes.txt):"
cat "$dir/routes.txt"
