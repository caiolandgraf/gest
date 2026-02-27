package main

import "fmt"

// Calculator performs basic arithmetic operations.
type Calculator struct{}

func (c Calculator) Add(a, b float64) float64      { return a + b }
func (c Calculator) Subtract(a, b float64) float64 { return a - b }
func (c Calculator) Multiply(a, b float64) float64 { return a * b }
func (c Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}
