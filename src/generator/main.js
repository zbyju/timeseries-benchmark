function generateNumbers(x, n) {
  return new Array(n).fill(x);
}

function rn(min, max) {
  return Math.floor(Math.random() * (max - min) + min);
}

function ri(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function simulateBatteryVoltage(voltage, n) {
  const result = [];
  let currentVoltage = voltage;
  let count = 0;

  for (let i = 0; i < n; i++) {
    if (count === 0) {
      // Decide whether to add or subtract 0.15, or to keep the voltage the same.
      currentVoltage += (Math.random() < 0.5 ? 0.15 : -0.15) * (Math.random() < 0.5 ? 1 : 0);
      // Randomly choose a count around 20 for how long this voltage will be held
      count = Math.floor(Math.random() * 5) + 18;
    }

    result.push(currentVoltage);
    count--;
  }

  return result;
}

function simulateBatteryVoltageWithFailure(voltage, n) {
  const result = []
  let currentVoltage = voltage
  let count = 0
  let isFailed = false

  for (let i = 0; i < n; i++) {
    if (count === 0) {
      if (isFailed) {
        // After failure, reset to the original voltage
        isFailed = false
      }
      // Check for failure with a 0.5% chance
      if (Math.random() < 0.1) {
        currentVoltage = 0
        count = Math.floor(Math.random() * 10) + 25 // Around 30 samples
        isFailed = true
      } else {
        // Adjust voltage up or down by 0.15, or keep the same
        currentVoltage = voltage + (Math.random() < 0.5 ? 0.15 : -0.15) * (Math.random() < 0.5 ? 1 : 0)
        count = Math.floor(Math.random() * 5) + 18 // Around 20 samples
      }
    }

    result.push(currentVoltage);
    count--;
  }

  return result;
}


function simulateDailyTemperature(startTemp, n) {
  const result = []
  let currentTemp = startTemp
  const period = 144 // One full cycle (day) in number of samples

  for (let i = 0; i < n; i++) {
    // Oscillation component (sine wave) to simulate day-night cycle
    // The amplitude of 5 degrees Celsius is an example; adjust as needed
    let oscillation = Math.sin((2 * Math.PI / period) * i) * 0.3

    // Random walk component for small random fluctuations
    let randomWalk = (Math.random() - 0.5) * 2

    // Update the temperature
    currentTemp += randomWalk + oscillation

    result.push(currentTemp)
  }

  return result
}

function simulateTemperatureWithControl(startTemp, n, controlMode, tempChange) {
  const result = [];
  let currentTemp = startTemp;
  const period = n; // One full cycle (day) in number of samples
  let controlCountdown = 0;

  for (let i = 0; i < n; i++) {
    // Oscillation component (sine wave) for daily cycle
    let dailyCycle = Math.sin((2 * Math.PI / period) * i) * 0.1;

    // Random walk component for small random fluctuations
    let randomWalk = (Math.random() - 0.5) * 2;

    // Probability of control action (heating or cooling) based on the sine value
    let controlProbability = controlMode === 'heating' ? (1 - dailyCycle) / 50 : (1 + dailyCycle) / 50;

    // Check if control action needs to be activated
    if (Math.random() < controlProbability && controlCountdown === 0) {
      controlCountdown = Math.floor(Math.random() * 3) + 4; // Random duration around 5 samples
    }

    // Apply control action if in the respective phase
    if (controlCountdown > 0) {
      // Apply a temperature change with reduced effect of daily cycle and random walk
      currentTemp += randomWalk + dailyCycle;
      controlCountdown--;
      result.push(currentTemp + tempChange);
    } else {
      // Update the temperature with daily cycle and random walk
      currentTemp += randomWalk + dailyCycle;
      result.push(currentTemp);
    }

  }

  return result;
}
/*
const n = 144
const batterySimulation = simulateBatteryVoltageWithFailure(14, n)
const dailySimulation = simulateDailyTemperature(25, n)
const heatingSimulation = simulateTemperatureWithControl(25, n, 'heating', 30);
const coolingSimulation = simulateTemperatureWithControl(25, n, 'cooling', -15);
console.log("Battery")
batterySimulation.forEach(x => console.log(x))
console.log("Daily")
dailySimulation.forEach(x => console.log(x))
console.log("Heating")
heatingSimulation.forEach(x => console.log(x))
console.log("Cooling")
coolingSimulation.forEach(x => console.log(x)) */

module.exports = {
  simulateBatteryVoltageWithFailure,
  simulateDailyTemperature,
  simulateTemperatureWithControl
}
