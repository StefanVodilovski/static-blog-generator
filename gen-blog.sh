#!/bin/bash

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        --input)
            INPUT_FOLDER="$2"
            shift
            shift
            ;;
        --output)
            OUTPUT_FOLDER="$2"
            shift
            shift
            ;;
        --titleName)
            TITLE_NAME="$2"
            shift
            shift
            ;;
        --posts-per-page)
            POSTS_PER_PAGE="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Check if required arguments are provided
if [ -z "$INPUT_FOLDER" ] || [ -z "$OUTPUT_FOLDER" ] || [ -z "$TITLE_NAME" ]; then
    echo "Usage: gen-blog --input $INPUT_FOLDER --output $OUTPUT_FOLDER --titleName $TITLE_NAME [--posts-per-page $POSTS_PER_PAGE]"
    exit 1
fi
echo "$INPUT_FOLDER"
echo $OUTPUT_FOLDER
echo $TITLE_NAME

# Run Docker command
docker run -v $INPUT_FOLDER:/app/input -v $OUTPUT_FOLDER:/app/output gen-blog generate --input /app/input --output /app/output --title $TITLE_NAME ${POSTS_PER_PAGE:+--posts-per-page "$POSTS_PER_PAGE"}