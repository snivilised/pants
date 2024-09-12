#!/bin/bash

# Define paths relative to the root directory
ROOT_DIR="$(git rev-parse --show-toplevel)"
COVERAGE_FILE="$ROOT_DIR/coverage.out"
EXCLUSIONS_FILE="$ROOT_DIR/scripts/coverage-exclusion-list.txt"

# Check if required files exist
if [ ! -f "$COVERAGE_FILE" ]; then
    echo "Error: Coverage file not found at $COVERAGE_FILE"
    exit 1
fi

if [ ! -f "$EXCLUSIONS_FILE" ]; then
    echo "Error: Exclusions file not found at $EXCLUSIONS_FILE"
    exit 1
fi

# Create a temporary file
TEMP_FILE=$(mktemp)

# Process the exclusions
while IFS= read -r line || [[ -n "$line" ]]; do
    # Escape special characters in the line for use in sed
    escaped_line=$(echo "$line" | sed 's/[\/&]/\\&/g')
    
    # Remove matching lines from the coverage file
    sed "/${escaped_line}/d" "$COVERAGE_FILE" > "$TEMP_FILE"
    
    # Replace the original file with the modified content
    mv "$TEMP_FILE" "$COVERAGE_FILE"
done < "$EXCLUSIONS_FILE"

echo "Coverage exclusions have been applied successfully."
