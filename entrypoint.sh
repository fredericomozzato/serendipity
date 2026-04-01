#!/bin/bash
set -e

# Remove stale server PID (prevents crash on container restart)
rm -f /app/tmp/pids/server.pid

exec "$@"
