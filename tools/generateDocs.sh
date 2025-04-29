#!/bin/bash

swag init -o ./docs
# Check if the command was successful
if [ $? -eq 0 ]; then
    echo "Swagger documentation generated successfully."
else
    echo "Failed to generate Swagger documentation."
    exit 1
fi
# Check if the docs directory exists
if [ -d "./docs" ]; then
    echo "Docs directory exists."
else
    echo "Docs directory does not exist."
    exit 1
fi
# Check if the docs directory contains the expected files
if [ -f "./docs/swagger.json" ] && [ -f "./docs/swagger.yaml" ]; then
    echo "Docs directory contains the expected files."
else
    echo "Docs directory does not contain the expected files."
    exit 1
fi
