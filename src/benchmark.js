const { simulateBatteryVoltageWithFailure, simulateDailyTemperature, simulateTemperatureWithControl } = require("./generator/main")
const { addDataInfluxDb, queryDataInfluxDb } = require("./runners/influx")
const Benchmark = require('benchmark');
const addSuite = new Benchmark.Suite;
const querySuite = new Benchmark.Suite;


const NUM_SAMPLES = 144
const NUM_SIMULATIONS = 10
const NUM_STATIONS = 100

function generateData() {
  let result = []
  let date = new Date()
  for (let stationId = 0; stationId < NUM_STATIONS; stationId++) {
    let simulations = []
    date = new Date(date.getTime() - 86400000)
    for (let simulation = 0; simulation < NUM_SIMULATIONS; simulation++) {
      const batterySimulation = simulateBatteryVoltageWithFailure(14, NUM_SAMPLES)
      const dailySimulation = simulateDailyTemperature(25, NUM_SAMPLES)
      const heatingSimulation = simulateTemperatureWithControl(25, NUM_SAMPLES, 'heating', 30);
      const coolingSimulation = simulateTemperatureWithControl(25, NUM_SAMPLES, 'cooling', -15);

      let samples = []
      for (let i = 0; i < NUM_SAMPLES; i++) {
        const b = batterySimulation[i];
        const d = dailySimulation[i];
        const h = heatingSimulation[i];
        const c = coolingSimulation[i];

        samples.push({
          stationId,
          timestamp: date,
          voltage: b,
          daily: d,
          heating: h,
          cooling: c
        })
        date = new Date(date.getTime() - 10 * 60000)
      }
      simulations.push(samples)
    }
    result.push(simulations)
  }
  return result
}

const data = generateData()


// Mock functions to represent data insertion and querying
// Replace these with your actual database operation functions
const addDataInflux = async () => {
  for (let stationId = 0; stationId < NUM_STATIONS; stationId++) {
    for (let simulation = 0; simulation < NUM_SIMULATIONS; simulation++) {
      for (let i = 0; i < NUM_SAMPLES; i++) {
        console.log("before1")
        await addDataInfluxDb(data[stationId][simulation][i])
        console.log("after1")
      }
    }
  }
};
const addDataTimescale = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 2000);
    setTimeout(resolve, delay);
  });
};
const addDataQuest = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 2000);
    setTimeout(resolve, delay);
  });
};
const addDataPostgresql = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 2000);
    setTimeout(resolve, delay);
  });
};
const addDataMongo = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 2000);
    setTimeout(resolve, delay);
  });
};

const queryDataInflux = async () => {
  for (let stationId = 0; stationId < NUM_STATIONS; stationId++) {
    await queryDataInfluxDb(stationId)
  }
};
const queryDataTimescale = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 1000);
    setTimeout(resolve, delay);
  });
};
const queryDataQuest = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 1000);
    setTimeout(resolve, delay);
  });
};
const queryDataPostgresql = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 1000);
    setTimeout(resolve, delay);
  });
};
const queryDataMongo = async () => {
  return new Promise(resolve => {
    const delay = Math.floor(Math.random() * 1000);
    setTimeout(resolve, delay);
  });
};


// Adding data benchmark
addSuite
  .add('InfluxDB - Add Data', {
    defer: true,
    fn: async (deferred) => {
      console.log("before")
      await addDataInflux()
      console.log("after")
      deferred.resolve();
    }
  })
  .add('TimescaleDB - Add Data', async () => {
    await addDataTimescale();
  })
  .add('QuestDB - Add Data', async () => {
    await addDataQuest();
  })
  .add('PostgreSQL - Add Data', async () => {
    await addDataPostgresql();
  })
  .add('MongoDB - Add Data', async () => {
    await addDataMongo();
  })
  .on('cycle', (event) => {
    console.log(String(event.target));
  })
  .on('complete', function() {
    console.log('Fastest add operation is ' + this.filter('fastest').map('name'));
  }).run({ async: true });

querySuite
  .add('InfluxDB - Query Data', {
    defer: true,
    fn: async (deferred) => {
      await queryDataInfluxDb()
      deferred.resolve();
    }
  })
  .add('TimescaleDB - Query Data', async () => {
    await queryDataTimescale();
  })
  .add('QuestDB - Query Data', async () => {
    await queryDataQuest();
  })
  .add('PostgreSQL - Query Data', async () => {
    await queryDataPostgresql();
  })
  .add('MongoDB - Query Data', async () => {
    await queryDataMongo();
  })
  .on('cycle', (event) => {
    console.log(String(event.target));
  })
  .on('complete', function() {
    console.log('Fastest query operation is ' + this.filter('fastest').map('name'));
  }).run({ async: true });
