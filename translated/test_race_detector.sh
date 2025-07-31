#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Counters
true_positives=0
false_negatives=0
true_negatives=0
false_positives=0
compilation_errors=0

echo -e "${BOLD}${BLUE}DataRaceBench Go Translation - Race Detector Evaluation${NC}"
echo -e "${BLUE}=========================================================${NC}"
echo ""
echo -e "${BOLD}Legend:${NC}"
echo -e "  ${GREEN}TP${NC} = True Positive (race expected and detected)"
echo -e "  ${RED}FN${NC} = False Negative (race expected but not detected)" 
echo -e "  ${GREEN}TN${NC} = True Negative (no race expected and none detected)"
echo -e "  ${RED}FP${NC} = False Positive (no race expected but detected)"
echo -e "  ${YELLOW}CE${NC} = Compilation Error"
echo ""

# Table header
printf "%-40s %-15s %-10s %-6s\n" "File" "Expected" "Detected" "Result"
echo "--------------------------------------------------------------------------------"

# Find all Go files and sort them
for file in $(find . -name "*.go" | sort); do
    filename=$(basename "$file")
    
    # Determine expected result based on filename
    if [[ "$filename" == *"-yes.go" ]]; then
        expected="RACE"
    elif [[ "$filename" == *"-no.go" ]]; then
        expected="NO_RACE"
    else
        expected="UNKNOWN"
    fi
    
    # Run with race detector and capture both stdout and stderr
    output=$(timeout 5m go run -race "$file" 2>&1)
    exit_code=$?
    
    # Check if compilation failed
    if echo "$output" | grep -q "compilation failed\|cannot find package\|undefined:"; then
        detected="COMPILE_ERR"
        result="${YELLOW}CE${NC}"
        ((compilation_errors++))
    # Check if race was detected
    elif echo "$output" | grep -q "WARNING: DATA RACE\|Found [0-9]* data race"; then
        detected="RACE"
        if [[ "$expected" == "RACE" ]]; then
            result="${GREEN}TP${NC}"
            ((true_positives++))
        else
            result="${RED}FP${NC}"
            ((false_positives++))
        fi
    # Check if program timed out
    elif [[ $exit_code -eq 124 ]]; then
        detected="TIMEOUT"
        result="${YELLOW}TO${NC}"
    # No race detected
    else
        detected="NO_RACE"
        if [[ "$expected" == "NO_RACE" ]]; then
            result="${GREEN}TN${NC}"
            ((true_negatives++))
        elif [[ "$expected" == "RACE" ]]; then
            result="${RED}FN${NC}"
            ((false_negatives++))
        else
            result="${YELLOW}??${NC}"
        fi
    fi
    
    # Print result row
    printf "%-40s %-15s %-10s %-6s\n" "$filename" "$expected" "$detected" "$(echo -e "$result")"
done

echo ""
echo -e "${BOLD}${BLUE}Summary Results:${NC}"
echo "================="
echo -e "${GREEN}True Positives (TP):${NC}  $true_positives"
echo -e "${RED}False Negatives (FN):${NC} $false_negatives" 
echo -e "${GREEN}True Negatives (TN):${NC}  $true_negatives"
echo -e "${RED}False Positives (FP):${NC} $false_positives"
echo -e "${YELLOW}Compilation Errors:${NC}   $compilation_errors"

total_tests=$((true_positives + false_negatives + true_negatives + false_positives))
total_correct=$((true_positives + true_negatives))

if [[ $total_tests -gt 0 ]]; then
    accuracy=$(echo "scale=2; $total_correct * 100 / $total_tests" | bc -l)
    echo ""
    echo -e "${BOLD}Overall Accuracy:${NC} ${accuracy}% (${total_correct}/${total_tests})"
fi

echo ""
echo -e "${BOLD}Race Detection Tool Performance:${NC}"
if [[ $((true_positives + false_negatives)) -gt 0 ]]; then
    sensitivity=$(echo "scale=2; $true_positives * 100 / ($true_positives + $false_negatives)" | bc -l)
    echo -e "Sensitivity (recall): ${sensitivity}% - ability to detect actual races"
fi

if [[ $((true_negatives + false_positives)) -gt 0 ]]; then
    specificity=$(echo "scale=2; $true_negatives * 100 / ($true_negatives + $false_positives)" | bc -l)
    echo -e "Specificity: ${specificity}% - ability to avoid false alarms"
fi

if [[ $((true_positives + false_positives)) -gt 0 ]]; then
    precision=$(echo "scale=2; $true_positives * 100 / ($true_positives + $false_positives)" | bc -l)
    echo -e "Precision: ${precision}% - accuracy of race detections"
fi

echo -e "${BLUE}=========================================================${NC}" 
