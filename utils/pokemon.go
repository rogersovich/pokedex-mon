package utils

import (
	"fmt"
	"math"
	"strings"
)

func roundToDecimalPlaces(val float64, decimalPlaces int) float64 {
	// Calculate the power of 10
	powerOfTen := math.Pow(10, float64(decimalPlaces))

	// Multiply, round, then divide back
	rounded := math.Round(val*powerOfTen) / powerOfTen
	return rounded
}

func CalcCaptureRate(baseCatchRate int, ballMultiplier float64, statusMultiplier float64) float64 {
	if ballMultiplier < 0 {
		ballMultiplier = 0
	}
	if statusMultiplier < 0 {
		statusMultiplier = 0
	}

	a := math.Floor(float64(baseCatchRate) * ballMultiplier * statusMultiplier / 3.0)

	if a < 0 {
		a = 0
	}

	percentage := (a / 255.0) * 100.0

	return roundToDecimalPlaces(percentage, 1)
}

func ConvertGrowthRate(growthRate string) string {
	if growthRate == "" {
		return ""
	}

	// Split the input string by the hyphen
	parts := strings.Split(growthRate, "-")

	// Create a slice to hold the capitalized parts
	capitalizedParts := make([]string, len(parts))

	for i, part := range parts {
		// Capitalize the first letter of each part using our helper function
		capitalizedParts[i] = CapitalizeFirstLetter(part)
	}

	// Join the capitalized parts with a space
	return strings.Join(capitalizedParts, " ")
}

func CalcGenderDistribution(genderRate int) GenderDistributionResult {
	var femaleChance float64
	var maleChance float64

	if genderRate == -1 {
		// Pokémon tanpa gender (genderless)
		femaleChance = 0.0
		maleChance = 0.0
	} else {
		// Menghitung peluang betina dalam bentuk desimal (e.g., 1/8 = 0.125)
		femaleChance = float64(genderRate) / 8.0

		// Memastikan peluang berada dalam rentang 0.0 hingga 1.0
		if femaleChance < 0.0 {
			femaleChance = 0.0
		} else if femaleChance > 1.0 {
			femaleChance = 1.0
		}

		// Menghitung peluang jantan
		maleChance = 1.0 - femaleChance
	}

	// Mengkonversi peluang desimal ke persentase string (misal: "12.5%")
	return GenderDistributionResult{
		Female: fmt.Sprintf("%.1f%%", femaleChance*100.0),
		Male:   fmt.Sprintf("%.1f%%", maleChance*100.0),
	}
}

func CalcEggCycles(eggCycles int) int {
	cycles := eggCycles * 255

	return cycles
}

func CalcMinStat(baseStat int, typeStat string) int {
	var minStat = 0

	if typeStat == "hp" {
		minStat = baseStat*2 + 110
	} else {
		minStat = int(float64(baseStat*2+5) * 0.9)
	}

	return minStat
}

func CalcMaxStat(baseStat int, typeStat string) int {
	var maxStat = 0

	if typeStat == "hp" {
		maxStat = baseStat*2 + 204
	} else {
		maxStat = int(float64(baseStat*2+99) * 1.1)
	}

	return maxStat
}
