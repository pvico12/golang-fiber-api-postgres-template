#!/bin/bash

# This script runs a swag command to generate Swagger documentation for the Go project.
# It dumps the output in the ./ci-expected/docs folder.
# It then checks if the documents in the ./docs folder are up to date and equal to the ones in the ./ci-expected/docsdocs folder.

# If the documents are not up to date, it will print a message and exit with a non-zero status.

# Check if the docs directory exists
if [ -d "./docs" ]; then
    echo "Docs directory exists."
else
    echo "Docs directory does not exist."
    exit 1
fi

# Generate Expected Docs
swag init -o ./ci-expected/docs
if [ $? -eq 0 ]; then
    echo "Expected documentation generated successfully."
else
    echo "Failed to generate Swagger documentation."
    exit 1
fi

# Check if the ci-expected/docs directory exists
if [ -d "./ci-expected/docs" ]; then
    echo "Expected docs directory exists."
else
    echo "Expected docs directory does not exist."
    exit 1
fi

# Check if the ci-expected/docs directory contains the expected files
if [ -f "./ci-expected/docs/swagger.json" ] && [ -f "./ci-expected/docs/swagger.yaml" ]; then
    echo "Expected docs directory contains the expected files."
else
    echo "Expected docs directory does not contain the expected files."
    exit 1
fi

# Compare the docs directory with the expected-docs directory
if diff -qr ./docs ./ci-expected/docs; then
    echo "Docs are up to date."
else
    echo "Docs are not up to date. Please run the script again to update them."
    exit 1
fi

# Check if the docs directory contains the expected files
if [ -f "./docs/swagger.json" ] && [ -f "./docs/swagger.yaml" ]; then
    echo "Docs directory contains the expected files."
else
    echo "Docs directory does not contain the expected files."
    exit 1
fi

# Function to extract routes from a given file
extract_routes_from_file() {
    local file=$1
    local basefilename=$(basename "$file" .go)

    # Check if file exists
    if [[ ! -f "$file" ]]; then
        echo "Error: Input file '$file' not found."
        exit 1
    fi

    # Parse the file to extract Method and Route
    while IFS= read -r line; do
        # Check if line contains api.Get, api.Post, api.Put, api.Delete
        if [[ $line =~ api\.(Get|Post|Put|Delete)\( ]]; then
            # Extract the method (Get, Post, Put, Delete)
            method=$(echo "$line" | grep -oP 'api\.\K(Get|Post|Put|Delete)' | tr '[:lower:]' '[:upper:]')
            
            # Try to extract route from current line first
            if [[ $line =~ \([\s]*\"([^\"]+)\" ]]; then
                route=$(echo "$line" | grep -oP '\([\s]*\"\K[^\"]+')
            else
                # If not found, read next line and try to extract route
                read -r next_line
                if [[ $next_line =~ \"([^\"]+)\" ]]; then
                    route=$(echo "$next_line" | grep -oP '\"\K[^\"]+')
                fi
            fi
            

            # add base filename to the route
            route=$(echo "$route" | sed "s|^|/$basefilename|")
            
            # If both method and route are found, print in desired format
            if [[ -n "$method" && -n "$route" ]]; then
                echo "$method $route"
            fi
        fi
    done < "$file"
}

# Function to extract documented routes from a given file
extract_documented_routes_from_file() {
    local file=$1
    grep -E '//\s*@Router\s+[^ ]+\s+\[.*\]' "$file" | while read -r line; do
        # Extract the route and method
        route=$(echo "$line" | sed -E 's/.*@Router\s+([^ ]+).*/\1/')
        method=$(echo "$line" | sed -E 's/.*\[(.*)\].*/\1/' | tr '[:lower:]' '[:upper:]')
        # Print in the desired format
        echo "$method $route"
    done
}

# Compare routes in a pair of files
compare_routes_in_files() {
    local router_file=$1
    local service_file=$2

    echo "Comparing routes in $router_file and $service_file..."

    router_routes=$(extract_routes_from_file "$router_file")
    service_routes=$(extract_documented_routes_from_file "$service_file")

    echo "Router routes: $router_routes"
    echo "Service routes: $service_routes"

    mismatched_routes=$(comm -3 <(echo "$router_routes" | sort) <(echo "$service_routes" | sort))

    if [ -n "$mismatched_routes" ]; then
        echo "Mismatched routes found between $router_file and $service_file:"
        echo "$mismatched_routes"
        exit 1
    else
        echo "All routes match between $router_file and $service_file."
    fi
}

# Pull all filenames in the routers and services directories
router_files=$(find ./routers -maxdepth 1 -type f -name "*.go")
service_files=$(find ./services -maxdepth 1 -type f -name "*.go")

# Compare files with matching filenames
for router_file in $router_files; do
    filename=$(basename "$router_file")
    service_file=$(find ./services -maxdepth 1 -type f -name "$filename")

    if [ -n "$service_file" ]; then
        compare_routes_in_files "$router_file" "$service_file"
    else
        echo "No matching service file found for $router_file"
        exit 1
    fi
done