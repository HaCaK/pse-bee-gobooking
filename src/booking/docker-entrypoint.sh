#!/bin/sh

set -e

# Wait for DB
if [ -n "$DB_CONNECT" ]; then
    ./wait-for-it.sh "$DB_CONNECT" -t 20
fi

# Wait for property
if [ -n "$PROPERTY_CONNECT" ]; then
  ./wait-for-it.sh "$PROPERTY_CONNECT" -t 20
fi

# Run the main container command.
exec "$@"
