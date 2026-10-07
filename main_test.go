package main

import (
	"math"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	data := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	stats, err := CalculateStats(data)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if stats.Count != 5 {
		t.Errorf("Expected count 5, got %d", stats.Count)
	}

	if stats.Min != 10.0 || stats.Max != 50.0 {
		t.Errorf("Expected Min 10.0 and Max 50.0, got Min %.2f and Max %.2f", stats.Min, stats.Max)
	}

	if stats.Mean != 30.0 {
		t.Errorf("Expected mean 30.0, got %f", stats.Mean)
	}

	if stats.Median != 30.0 {
		t.Errorf("Expected median 30.0, got %f", stats.Median)
	}

	expectedStdDev := 14.1421
	if math.Abs(stats.StdDev-expectedStdDev) > 0.001 {
		t.Errorf("Expected stddev ~%.4f, got %.4f", expectedStdDev, stats.StdDev)
	}
}

func TestEmptyDataset(t *testing.T) {
	_, err := CalculateStats([]float64{})
	if err == nil {
		t.Errorf("Expected error for empty dataset, got nil")
	}
}
