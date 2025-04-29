#!/bin/bash

# Function to display usage
usage() {
    echo "Usage: $0 [--keep-alive | -k]"
    exit 1
}

# Parse command-line arguments
KEEP_ALIVE=false
while [[ $# -gt 0 ]]; do
    case "$1" in
        --keep-alive|-k)
            KEEP_ALIVE=true
            shift
            ;;
        *)
            usage
            ;;
    esac
done

# Start PostgreSQL container
docker compose -f tests/postgres.docker-compose.yml up -d
sleep 5 # wait for postgres to start

# Run tests
TESTING_MODE=true go test ./tests
TEST_EXIT_CODE=$?

# Stop PostgreSQL container unless --keep-alive is specified
if [[ "$KEEP_ALIVE" == false ]]; then
    docker compose -f tests/postgres.docker-compose.yml down
fi

# Exit with the test command's exit code
exit $TEST_EXIT_CODE
