#!/bin/bash

FUNCTION_NAME=$1

echo "Building function: $FUNCTION_NAME"

go build -o build/$FUNCTION_NAME \
functions/$FUNCTION_NAME/main.go

echo "Build complete"