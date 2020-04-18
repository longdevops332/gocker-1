package images

import "fmt"

// Return file size in human readable format
func sizeOfFmt(num float64) string {
	formats := [4]string{"", "Ki", "Mi", "Gi"}
	step := 1024.0
	for _, f := range formats {
		if num < step {
			return fmt.Sprintf("%f%sB", num, f)
		}
		step /= 1024
	}
	return ""
}
