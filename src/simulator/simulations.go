package simulator

import (
	"math"
	"math/rand"
	"time"
)

func SimulateBatteryVoltageWithFailure(voltage float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	result := make([]float64, 0, n)
	currentVoltage := voltage
	count := 0
	isFailed := false

	for i := 0; i < n; i++ {
		if count == 0 {
			if isFailed {
				isFailed = false
			} else if rand.Float64() < 0.1 {
				currentVoltage = 0
				count = rand.Intn(10) + 25
				isFailed = true
			} else {
				change := 0.15
				if rand.Float64() >= 0.5 {
					change = -change
				}
				if rand.Float64() >= 0.5 {
					currentVoltage = voltage + change
				}
				count = rand.Intn(5) + 18
			}
		}

		result = append(result, currentVoltage)
		count--
	}

	return result
}

func SimulateDailyTemperature(startTemp float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	result := make([]float64, 0, n)
	currentTemp := startTemp
	period := float64(n)

	for i := 0; i < n; i++ {
		oscillation := math.Sin((2*math.Pi/period)*float64(i)) * 0.3
		randomWalk := (rand.Float64() - 0.5) * 2
		currentTemp += randomWalk + oscillation
		result = append(result, currentTemp)
	}

	return result
}

func SimulateTemperatureWithControl(startTemp float64, n int, controlMode string, tempChange float64) []float64 {
	rand.Seed(time.Now().UnixNano())
	result := make([]float64, 0, n)
	currentTemp := startTemp
	period := float64(n)
	controlCountdown := 0

	for i := 0; i < n; i++ {
		dailyCycle := math.Sin((2*math.Pi/period)*float64(i)) * 0.1
		randomWalk := (rand.Float64() - 0.5) * 2

		var controlProbability float64
		if controlMode == "heating" {
			controlProbability = (1 - dailyCycle) / 50
		} else {
			controlProbability = (1 + dailyCycle) / 50
		}

		if rand.Float64() < controlProbability && controlCountdown == 0 {
			controlCountdown = rand.Intn(3) + 4
		}

		if controlCountdown > 0 {
			currentTemp += randomWalk + dailyCycle
			controlCountdown--
			result = append(result, currentTemp+tempChange)
		} else {
			currentTemp += randomWalk + dailyCycle
			result = append(result, currentTemp)
		}
	}

	return result
}
