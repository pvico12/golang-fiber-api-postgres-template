#!/bin/bash

docker-compose -f tests/postgres.docker-compose.yml up -d
sleep 5 # wait for postgres to start
TESTING_MODE=true go test ./tests
docker-compose -f tests/postgres.docker-compose.yml down
