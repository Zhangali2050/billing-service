#!/bin/sh
# wait-for.sh

set -e

host="$1"
shift

until nc -z $host; do
  echo "⏳ Waiting for $host..."
  sleep 1
done

echo "✅ $host is available — starting app..."
exec "$@"
