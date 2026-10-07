package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Stats struct {
	Count    int     `json:"count"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Mean     float64 `json:"mean"`
	Median   float64 `json:"median"`
	Variance float64 `json:"variance"`
	StdDev   float64 `json:"std_dev"`
	P25      float64 `json:"p25"`
	P75      float64 `json:"p75"`
	IQR      float64 `json:"iqr"`
}

func CalculateStats(numbers []float64) (Stats, error) {
	n := len(numbers)
	if n == 0 {
		return Stats{}, fmt.Errorf("dataset is empty")
	}

	data := make([]float64, n)
	copy(data, numbers)
	sort.Float64s(data)

	minVal := data[0]
	maxVal := data[n-1]

	var sum float64
	for _, num := range data {
		sum += num
	}
	mean := sum / float64(n)

	var median float64
	if n%2 == 0 {
		median = (data[n/2-1] + data[n/2]) / 2.0
	} else {
		median = data[n/2]
	}

	var varianceSum float64
	for _, num := range data {
		varianceSum += math.Pow(num-mean, 2)
	}
	variance := varianceSum / float64(n)
	stdDev := math.Sqrt(variance)

	p25 := calculatePercentile(data, 25)
	p75 := calculatePercentile(data, 75)
	iqr := p75 - p25

	return Stats{
		Count:    n,
		Min:      minVal,
		Max:      maxVal,
		Mean:     mean,
		Median:   median,
		Variance: variance,
		StdDev:   stdDev,
		P25:      p25,
		P75:      p75,
		IQR:      iqr,
	}, nil
}

func calculatePercentile(sortedData []float64, percentile float64) float64 {
	index := (percentile / 100.0) * float64(len(sortedData)-1)
	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	if upper >= len(sortedData) {
		return sortedData[lower]
	}
	return sortedData[lower]*(1-weight) + sortedData[upper]*weight
}

func ExportToJSON(stats Stats, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stats)
}

func ExportToCSV(stats Stats, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Metric", "Value"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	records := [][]string{
		{"Count", strconv.Itoa(stats.Count)},
		{"Min", fmt.Sprintf("%.4f", stats.Min)},
		{"Max", fmt.Sprintf("%.4f", stats.Max)},
		{"Mean", fmt.Sprintf("%.4f", stats.Mean)},
		{"Median", fmt.Sprintf("%.4f", stats.Median)},
		{"Variance", fmt.Sprintf("%.4f", stats.Variance)},
		{"StdDev", fmt.Sprintf("%.4f", stats.StdDev)},
		{"P25", fmt.Sprintf("%.4f", stats.P25)},
		{"P75", fmt.Sprintf("%.4f", stats.P75)},
		{"IQR", fmt.Sprintf("%.4f", stats.IQR)},
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	inputFile := flag.String("input", "data.txt", "Path to input file")
	format := flag.String("format", "text", "Output format: text, json, csv")
	outputFile := flag.String("output", "", "Path to output file")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening file '%s': %v\n", *inputFile, err)
		os.Exit(1)
	}
	defer file.Close()

	var numbers []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		numbers = append(numbers, val)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	stats, err := CalculateStats(numbers)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(*format) {
	case "json":
		if *outputFile != "" {
			if err := ExportToJSON(stats, *outputFile); err != nil {
				fmt.Printf("Error exporting JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Results exported to JSON: %s\n", *outputFile)
		} else {
			data, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Println(string(data))
		}
	case "csv":
		if *outputFile != "" {
			if err := ExportToCSV(stats, *outputFile); err != nil {
				fmt.Printf("Error exporting CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Results exported to CSV: %s\n", *outputFile)
		} else {
			fmt.Println("Error: -output path is required for CSV format")
			os.Exit(1)
		}
	default:
		fmt.Printf("--- Scientific Data Analysis Report ---\n")
		fmt.Printf("Elements Processed: %d\n", stats.Count)
		fmt.Printf("Min Value:          %.4f\n", stats.Min)
		fmt.Printf("Max Value:          %.4f\n", stats.Max)
		fmt.Printf("Mean:               %.4f\n", stats.Mean)
		fmt.Printf("Median:             %.4f\n", stats.Median)
		fmt.Printf("Variance:           %.4f\n", stats.Variance)
		fmt.Printf("Standard Deviation: %.4f\n", stats.StdDev)
		fmt.Printf("25th Percentile:    %.4f\n", stats.P25)
		fmt.Printf("75th Percentile:    %.4f\n", stats.P75)
		fmt.Printf("IQR:                %.4f\n", stats.IQR)
	}
}
