#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Processor counts to test
PROCESSORS=(3 36 45 72 90 180 256)

# Problem sizes to test for -var- cases
PROBLEM_SIZES=(32 64 128 256 512 1024)

# Counters
true_positives=0
false_negatives=0
true_negatives=0
false_positives=0
compilation_errors=0

echo -e "${BOLD}${BLUE}DataRaceBench Go Translation - Race Detector Evaluation${NC}"
echo -e "${BLUE}=========================================================${NC}"
echo ""
echo -e "${BOLD}Testing with processor counts:${NC} ${PROCESSORS[*]}"
echo -e "${BOLD}Testing with problem sizes:${NC} ${PROBLEM_SIZES[*]} (for -var- cases)"
echo ""
echo -e "${BOLD}Legend:${NC}"
echo -e "  ${GREEN}TP${NC} = True Positive (race expected and detected)"
echo -e "  ${RED}FN${NC} = False Negative (race expected but not detected)" 
echo -e "  ${GREEN}TN${NC} = True Negative (no race expected and none detected)"
echo -e "  ${RED}FP${NC} = False Positive (no race expected but detected)"
echo -e "  ${YELLOW}CE${NC} = Compilation Error"
echo ""

# Table header
printf "%-40s %-15s %-10s %-6s %-8s %-8s\n" "File" "Expected" "Detected" "Result" "Procs" "Size"
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
    
    # Check if this is a variable-size test case
    is_var_case=false
    if [[ "$filename" == *"-var-"* ]]; then
        is_var_case=true
    fi
    
    # Variables to track results across all processor counts and problem sizes
    race_detected_any=false
    compilation_error_any=false
    timeout_any=false
    
    # Run with each processor count
    for procs in "${PROCESSORS[@]}"; do
        # For variable-size cases, also iterate over problem sizes
        if [[ "$is_var_case" == true ]]; then
            for size in "${PROBLEM_SIZES[@]}"; do
                # Run with race detector and capture both stdout and stderr
                output=$(GOMAXPROCS=$procs GORACE="halt_on_error=1" gtimeout 5m go run -race "$file" "$size" 2>&1)
                exit_code=$?
                
                # Check if compilation failed
                if echo "$output" | grep -q "compilation failed\|cannot find package\|undefined:"; then
                    detected="COMPILE_ERR"
                    result="${YELLOW}CE${NC}"
                    compilation_error_any=true
                # Check if race was detected
                elif echo "$output" | grep -q "WARNING: DATA RACE\|Found [0-9]* data race"; then
                    detected="RACE"
                    result="${GREEN}RACE${NC}"
                    race_detected_any=true
                # Check if program timed out
                elif [[ $exit_code -eq 124 ]]; then
                    detected="TIMEOUT"
                    result="${YELLOW}TO${NC}"
                    timeout_any=true
                # No race detected
                else
                    detected="NO_RACE"
                    result="${BLUE}NO_RACE${NC}"
                fi
                
                # Print result row
                printf "%-40s %-15s %-10s %-6s %-8s %-8s\n" "$filename" "$expected" "$detected" "$(echo -e "$result")" "$procs" "$size"
                
                # If race was detected, no need to test other combinations
                if [[ "$race_detected_any" == true ]]; then
                    break 2  # Break out of both loops
                fi
            done
        else
            # Regular case - no problem size iteration
            # Run with race detector and capture both stdout and stderr
            output=$(GOMAXPROCS=$procs GORACE="halt_on_error=1" gtimeout 5m go run -race "$file" 2>&1)
            exit_code=$?
            
            # Check if compilation failed
            if echo "$output" | grep -q "compilation failed\|cannot find package\|undefined:"; then
                detected="COMPILE_ERR"
                result="${YELLOW}CE${NC}"
                compilation_error_any=true
            # Check if race was detected
            elif echo "$output" | grep -q "WARNING: DATA RACE\|Found [0-9]* data race"; then
                detected="RACE"
                result="${GREEN}RACE${NC}"
                race_detected_any=true
            # Check if program timed out
            elif [[ $exit_code -eq 124 ]]; then
                detected="TIMEOUT"
                result="${YELLOW}TO${NC}"
                timeout_any=true
            # No race detected
            else
                detected="NO_RACE"
                result="${BLUE}NO_RACE${NC}"
            fi
            
            # Print result row
            printf "%-40s %-15s %-10s %-6s %-8s %-8s\n" "$filename" "$expected" "$detected" "$(echo -e "$result")" "$procs" "N/A"
            
            # If race was detected, no need to test other thread counts
            if [[ "$race_detected_any" == true ]]; then
                break
            fi
        fi
    done
    
    # Aggregate results across all processor counts and problem sizes for this test case
    if [[ "$expected" == "RACE" ]]; then
        if [[ "$race_detected_any" == true ]]; then
            final_result="${GREEN}TP${NC}"
            ((true_positives++))
        else
            final_result="${RED}FN${NC}"
            ((false_negatives++))
        fi
    elif [[ "$expected" == "NO_RACE" ]]; then
        if [[ "$race_detected_any" == true ]]; then
            final_result="${RED}FP${NC}"
            ((false_positives++))
        else
            final_result="${GREEN}TN${NC}"
            ((true_negatives++))
        fi
    else
        final_result="${YELLOW}??${NC}"
    fi
    
    # Print aggregated result
    printf "%-40s %-15s %-10s %-6s %-8s %-8s\n" "$filename (AGGREGATED)" "$expected" "$(if [[ "$race_detected_any" == true ]]; then echo "RACE"; else echo "NO_RACE"; fi)" "$(echo -e "$final_result")" "ALL" "ALL"
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
