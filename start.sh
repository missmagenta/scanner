#!/bin/bash

if [ "$DB_TYPE" = "mongodb" ]; then
  docker-compose --env-file mongodb_credentials.env up -d mongo-demo-v
elif [ "$DB_TYPE" = "clickhouse" ]; then
  docker-compose --env-file clickhouse_credentials.env up -d clickhouse-demo-v
else
  echo "Invalid or unspecified DB_TYPE. Usage: DB_TYPE=mongodb ./start.sh"
fi
