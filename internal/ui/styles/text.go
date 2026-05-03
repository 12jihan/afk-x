package styles

import "fmt"

// FormatResource renders a resource name-value pair with fixed-width alignment
// so values don't shift as they grow.
func FormatResource(name string, value float64) string {
	return fmt.Sprintf("%-20s %10.2f", name, value)
}
