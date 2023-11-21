const { InfluxDB, Point } = require('@influxdata/influxdb-client')

const token = "zoAUTXUsH2U1BrQgQwU8OQmF__b9g6UqFnxirpVLWAauGORjeiBVMsiuDFRFiNa53PVGNXa5Ia4L3DDvbEm0_w=="
const url = 'http://localhost:8086'
const org = `zbyju`
const bucket = `benchmark`
const client = new InfluxDB({ url, token })
const writeClient = client.getWriteApi(org, bucket);
// Function to add data
async function addDataInfluxDb(data) {
  const { stationId, timestamp, voltage, daily, heating, cooling } = data;

  const point = new Point('stationData')
    .tag('stationId', stationId)
    .floatField('voltage', voltage)
    .floatField('outerTemperature', daily)
    .floatField('heatingTemperature', heating)
    .floatField('coolingTemperature', cooling);

  writeClient.writePoint(point)
  writeClient.flush()
    .then(() => {
      console.log('WRITE FINISHED')
    })
    .catch(err => console.log(err))
}

// Function to query data from the last week for a given stationId
async function queryDataInfluxDb(stationId) {
  const queryClient = client.getQueryApi(org)
  const fluxQuery = `from(bucket: "benchmark")
 |> range(start: -10w)
 |> filter(fn: (r) => r.stationId == ${stationId})
 |> filter(fn: (r) => r._measurement == 'stationData')`

  queryClient.queryRows(fluxQuery, {
    next: (row, tableMeta) => {
      const tableObject = tableMeta.toObject(row)
      console.log(tableObject)
    },
    error: (error) => {
      console.error('\nError', error)
    },
    complete: () => {
      console.log('\nSuccess')
    },
  })
}

module.exports = { addDataInfluxDb, queryDataInfluxDb };
