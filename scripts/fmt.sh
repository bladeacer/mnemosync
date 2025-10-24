#!/bin/bash

# Find all Go files in the current directory and its subdirectories
# that need formatting (-l), then apply the formatting (-s -w).
echo "Finding and formatting all out-of-spec Go files..."

# The 'gofmt -l' output is fed directly to 'xargs', which runs the
# 'gofmt -s -w' command for each file path.
find . -name "*.go" -print0 | xargs -0 gofmt -s -w

echo "Formatting complete."
