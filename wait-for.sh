#!/bin/sh

# wait-for.sh postgres:5432 -- your_command_here

set -e

host="$1"
shift
cmd="$@"

echo "⏳ Waiting for $host to be available..."

until nc -z -v -w30 $(echo "$host" | cut -d: -f1) $(echo "$host" | cut -d: -f2); do
  echo "❌ $host is not available yet"
  sleep 1
done

echo "✅ $host is available — starting app..."
exec $cmd
