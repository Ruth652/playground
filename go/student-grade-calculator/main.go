package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Student Grade Calculator Console App

- Prompts for student name and number of subjects
- Collects each subject's name and grade with validation
- Calculates the average grade
- Prints a neatly formatted report with subjects, grades, and average
*/

func calculateAverage(grades []float64) float64 {
	if len(grades) == 0 {
		return 0
	}
	sum := 0.0
	for _, grade := range grades {
		sum += grade
	}
	return sum / float64(len(grades))
}

// Handles all input and prints report
func prompt() {
	reader := bufio.NewReader(os.Stdin)

	// Read student name
	fmt.Print("Enter student name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Read number of subjects
	numSubjects := 0
	for numSubjects <= 0 {
		fmt.Print("How many subjects? ")
		nStr, _ := reader.ReadString('\n')
		nStr = strings.TrimSpace(nStr)
		n, err := strconv.Atoi(nStr)

		if err != nil || n <= 0 {
			fmt.Println("Please enter a valid positive integer for number of subjects.")
			continue
		}
		numSubjects = n
	}

	subjects := make([]string, 0, numSubjects)
	grades := make([]float64, 0, numSubjects)
	gradesMap := make(map[string]float64)

	// Read subject names and grades
	for i := 0; i < numSubjects; i++ {
		fmt.Printf("Subject %d name: ", i+1)
		subj, _ := reader.ReadString('\n')
		subj = strings.TrimSpace(subj)

		if subj == "" {
			subj = fmt.Sprintf("Subject%d", i+1)
		}

		subjects = append(subjects, subj)

		// Read grade with validation
		var grade float64
		for {
			fmt.Printf("Grade for %s (0-100): ", subj)
			gStr, _ := reader.ReadString('\n')
			gStr = strings.TrimSpace(gStr)
			g, err := strconv.ParseFloat(gStr, 64)

			if err != nil {
				fmt.Println("Please enter a valid number for the grade.")
				continue
			}
			if g < 0 || g > 100 {
				fmt.Println("Grade must be between 0 and 100.")
				continue
			}
			grade = g
			break
		}
		grades = append(grades, grade)
		gradesMap[subj] = grade
	}

	// Calculate average
	avg := calculateAverage(grades)

	// Print report
	fmt.Println("\n--- Student Grade Report ---")
	fmt.Printf("Student: %s\n\n", name)
	fmt.Println("Subjects and grades:")

	for _, s := range subjects {
		fmt.Printf(" - %-20s : %.2f\n", s, gradesMap[s])
	}

	fmt.Printf("\nAverage Grade: %.2f\n", avg)
}

func main() {
	prompt()
}
