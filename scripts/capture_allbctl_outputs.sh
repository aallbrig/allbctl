#!/bin/bash

# Script to build allbctl at historical commits and capture command output
# 
# Usage: 
#   ./capture_allbctl_outputs.sh [OPTIONS]
#
# Options:
#   -l, --limit N         Number of commits to process (default: 10)
#   -c, --command CMD     Command to run with allbctl (default: "status")
#                         For multiple words, quote them: "config list"
#   -d, --output-dir DIR  Directory to save outputs (default: "data/outputs")
#   -h, --help           Show this help message
#
# Examples:
#   ./capture_allbctl_outputs.sh
#   ./capture_allbctl_outputs.sh --limit 5
#   ./capture_allbctl_outputs.sh -l 10 -c "config list"
#   ./capture_allbctl_outputs.sh --limit 20 --command status --output-dir results
#
# Output files are named: HEAD~N.<commit_hash>.<command>.output.txt
# Where command spaces are replaced with underscores (e.g., "config list" -> "config_list")

set -e

# Default values
NUM_COMMITS=10
COMMAND="status"
OUTPUT_DIR="data/outputs"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -l|--limit)
      NUM_COMMITS="$2"
      shift 2
      ;;
    -c|--command)
      COMMAND="$2"
      shift 2
      ;;
    -d|--output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    -h|--help)
      grep '^#' "$0" | grep -v '#!/bin/bash' | sed 's/^# //' | sed 's/^#//'
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Create sanitized command name for filenames (replace spaces with underscores)
COMMAND_FILENAME=$(echo "${COMMAND}" | tr ' ' '_')

echo "Building allbctl for last ${NUM_COMMITS} commits..."
echo "Command: allbctl ${COMMAND}"
echo "Output directory: ${OUTPUT_DIR}"
echo ""

# Create output directory if it doesn't exist
mkdir -p "${OUTPUT_DIR}"

# Get commit hashes dynamically
mapfile -t COMMIT_DATA < <(git --no-pager log -${NUM_COMMITS} --pretty=format:"%h" HEAD)

# Get current branch/commit to return to later
CURRENT_REF=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_REF" = "HEAD" ]; then
  CURRENT_REF=$(git rev-parse HEAD)
fi

# Stash any uncommitted changes (except this script)
SCRIPT_PATH=$(realpath "$0")
SCRIPT_REL_PATH=$(git ls-files --full-name "$SCRIPT_PATH" 2>/dev/null || echo "")

echo "Stashing any uncommitted changes (excluding this script)..."

# Get all dirty files
DIRTY_FILES=$(git diff --name-only && git diff --cached --name-only)

# Filter out this script and stash the rest
if [ -n "$DIRTY_FILES" ]; then
  FILES_TO_STASH=""
  while IFS= read -r file; do
    if [ "$file" != "$SCRIPT_REL_PATH" ]; then
      FILES_TO_STASH="$FILES_TO_STASH $file"
    fi
  done <<< "$DIRTY_FILES"
  
  if [ -n "$FILES_TO_STASH" ]; then
    git stash push -m "capture_allbctl_outputs.sh temporary stash" $FILES_TO_STASH > /dev/null 2>&1
    STASHED=$?
  else
    STASHED=1  # Nothing to stash
  fi
else
  STASHED=1  # Nothing to stash
fi

# Process each commit
PREV_NORMALIZED_OUTPUT=""
PREV_OUTPUT_FILE=""
RANGE_START_NUM=""
RANGE_START_HASH=""
RANGE_HASHES=()

for i in "${!COMMIT_DATA[@]}"; do
  short_hash="${COMMIT_DATA[$i]}"
  head_num=$i
  
  echo "Processing HEAD~${head_num} (${short_hash})..."
  
  # Checkout the commit
  git checkout -q HEAD~${head_num}
  
  # Recreate output directory (might be removed by checkout if untracked)
  mkdir -p "${OUTPUT_DIR}"
  
  # Clean previous build artifacts
  rm -rf bin/
  rm -f allbctl allbctl_*
  
  # Build allbctl directly with go build
  if go build -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'" -o bin/allbctl main.go > /dev/null 2>&1; then
    if [ -f "bin/allbctl" ]; then
      # Run allbctl with specified command and capture output
      raw_output=$(./bin/allbctl ${COMMAND} 2>&1)
      
      # Normalize output by removing timestamps and memory usage
      normalized_output=$(echo "$raw_output" | sed -E 's/\( [0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2} [A-Z]+ [+-][0-9]{4} \)//g' | sed -E 's/Memory:[[:space:]]+[0-9]+\.[0-9]+ [A-Z]iB \/ [0-9]+\.[0-9]+ [A-Z]iB/Memory: REDACTED/g')
      
      # Check if output is same as previous
      if [ "$normalized_output" = "$PREV_NORMALIZED_OUTPUT" ] && [ -n "$PREV_NORMALIZED_OUTPUT" ]; then
        # Same output - extend the range
        if [ -z "$RANGE_START_NUM" ]; then
          # Start a new range
          RANGE_START_NUM=$((head_num - 1))
          RANGE_START_HASH="${COMMIT_DATA[$((i - 1))]}"
          RANGE_HASHES=("${COMMIT_DATA[$((i - 1))]}")
        fi
        RANGE_HASHES+=("$short_hash")
        echo "  → Same as previous, extending range..."
        
        # Remove previous file if it exists
        if [ -f "$PREV_OUTPUT_FILE" ]; then
          rm "$PREV_OUTPUT_FILE"
        fi
      else
        # Different output or first commit
        if [ -n "$RANGE_START_NUM" ]; then
          # Save previous range
          range_end=$((i - 1))
          range_start=$RANGE_START_NUM
          commit_range="${RANGE_HASHES[0]}..${RANGE_HASHES[-1]}"
          output_file="${OUTPUT_DIR}/HEAD~${range_start}..HEAD~${range_end}.${commit_range}.${COMMAND_FILENAME}.output.txt"
          echo "$raw_output" > "$output_file"
          echo "✓ Created HEAD~${range_start}..HEAD~${range_end}.${commit_range}.${COMMAND_FILENAME}.output.txt"
          RANGE_START_NUM=""
          RANGE_HASHES=()
        fi
        
        # Save current output
        output_file="${OUTPUT_DIR}/HEAD~${head_num}.${short_hash}.${COMMAND_FILENAME}.output.txt"
        echo "$raw_output" > "$output_file"
        echo "✓ Created HEAD~${head_num}.${short_hash}.${COMMAND_FILENAME}.output.txt"
        PREV_OUTPUT_FILE="$output_file"
      fi
      
      PREV_NORMALIZED_OUTPUT="$normalized_output"
    else
      echo "✗ Failed to build at HEAD~${head_num} (${short_hash}) - binary not found"
      echo "Build failed for commit ${short_hash} - binary not created" > "${OUTPUT_DIR}/HEAD~${head_num}.${short_hash}.${COMMAND_FILENAME}.output.txt"
      PREV_NORMALIZED_OUTPUT=""
      PREV_OUTPUT_FILE=""
      RANGE_START_NUM=""
      RANGE_HASHES=()
    fi
  else
    echo "✗ Failed to build at HEAD~${head_num} (${short_hash}) - build error"
    echo "Build failed for commit ${short_hash} - compilation error" > "${OUTPUT_DIR}/HEAD~${head_num}.${short_hash}.${COMMAND_FILENAME}.output.txt"
    PREV_NORMALIZED_OUTPUT=""
    PREV_OUTPUT_FILE=""
    RANGE_START_NUM=""
    RANGE_HASHES=()
  fi
done

# Handle final range if exists
if [ -n "$RANGE_START_NUM" ]; then
  range_end=$((${#COMMIT_DATA[@]} - 1))
  range_start=$RANGE_START_NUM
  commit_range="${RANGE_HASHES[0]}..${RANGE_HASHES[-1]}"
  # Need to get the last output
  if [ -f "$PREV_OUTPUT_FILE" ]; then
    last_output=$(cat "$PREV_OUTPUT_FILE")
    rm "$PREV_OUTPUT_FILE"
    output_file="${OUTPUT_DIR}/HEAD~${range_start}..HEAD~${range_end}.${commit_range}.${COMMAND_FILENAME}.output.txt"
    echo "$last_output" > "$output_file"
    echo "✓ Created HEAD~${range_start}..HEAD~${range_end}.${commit_range}.${COMMAND_FILENAME}.output.txt"
  fi
fi

# Return to original ref
echo ""
echo "Returning to ${CURRENT_REF}..."
git checkout -q "${CURRENT_REF}"

# Restore stashed changes if any
if [ $STASHED -eq 0 ]; then
  echo "Restoring stashed changes..."
  git stash pop > /dev/null 2>&1
fi

# Rebuild at current commit
echo "Rebuilding at current commit..."
rm -rf bin/
rm -f allbctl allbctl_*
go build -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'" -o bin/allbctl main.go > /dev/null 2>&1

echo ""
echo "✓ Complete! Output files created in ${OUTPUT_DIR}/"
ls -1 "${OUTPUT_DIR}"/*.${COMMAND_FILENAME}.output.txt 2>/dev/null | wc -l | xargs echo "Total files created:"
