package integrations

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	imodels "github.com/influxdb/influxdb/models"
	"github.com/influxdb/kapacitor"
	"github.com/influxdb/kapacitor/clock"
	"github.com/influxdb/kapacitor/wlog"
)

func TestBatch_Derivative(t *testing.T) {

	var script = `
batch
	.query('''
		SELECT sum("value") as "value"
		FROM "telegraf"."default".packets
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))
	.derivative('value')
	.httpOut('TestBatch_Derivative')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "packets",
				Tags:    nil,
				Columns: []string{"time", "value"},
				Values: [][]interface{}{
					{
						time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 2, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 4, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 6, 0, time.UTC),
						0.5,
					},
				},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_Derivative", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(21 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_Derivative")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}

func TestBatch_DerivativeUnit(t *testing.T) {

	var script = `
batch
	.query('''
		SELECT sum("value") as "value"
		FROM "telegraf"."default".packets
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))
	.derivative('value')
		.unit(2s)
	.httpOut('TestBatch_Derivative')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "packets",
				Tags:    nil,
				Columns: []string{"time", "value"},
				Values: [][]interface{}{
					{
						time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						1.0,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 2, 0, time.UTC),
						1.0,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 4, 0, time.UTC),
						1.0,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 6, 0, time.UTC),
						1.0,
					},
				},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_Derivative", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(21 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_Derivative")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}

func TestBatch_DerivativeN(t *testing.T) {

	var script = `
batch
	.query('''
		SELECT sum("value") as "value"
		FROM "telegraf"."default".packets
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))
	.derivative('value')
	.httpOut('TestBatch_Derivative')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "packets",
				Tags:    nil,
				Columns: []string{"time", "value"},
				Values: [][]interface{}{
					{
						time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 2, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 4, 0, time.UTC),
						-501.0,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 6, 0, time.UTC),
						0.5,
					},
				},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_DerivativeNN", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(21 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_Derivative")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}
func TestBatch_DerivativeNN(t *testing.T) {

	var script = `
batch
	.query('''
		SELECT sum("value") as "value"
		FROM "telegraf"."default".packets
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))
	.derivative('value')
		.nonNegative()
	.httpOut('TestBatch_Derivative')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "packets",
				Tags:    nil,
				Columns: []string{"time", "value"},
				Values: [][]interface{}{
					{
						time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 2, 0, time.UTC),
						0.5,
					},
					{
						time.Date(1971, 1, 1, 0, 0, 6, 0, time.UTC),
						0.5,
					},
				},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_DerivativeNN", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(21 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_Derivative")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}

func TestBatch_SimpleMR(t *testing.T) {

	var script = `
batch
	.query('''
		SELECT mean("value")
		FROM "telegraf"."default".cpu_usage_idle
		WHERE "host" = 'serverA'
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s), 'cpu')
	.mapReduce(influxql.count('value'))
	.window()
		.period(20s)
		.every(20s)
	.mapReduce(influxql.sum('count'))
	.httpOut('TestBatch_SimpleMR')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "cpu_usage_idle",
				Tags:    map[string]string{"cpu": "cpu-total"},
				Columns: []string{"time", "sum"},
				Values: [][]interface{}{[]interface{}{
					time.Date(1971, 1, 1, 0, 0, 28, 0, time.UTC),
					10.0,
				}},
			},
			{
				Name:    "cpu_usage_idle",
				Tags:    map[string]string{"cpu": "cpu0"},
				Columns: []string{"time", "sum"},
				Values: [][]interface{}{[]interface{}{
					time.Date(1971, 1, 1, 0, 0, 28, 0, time.UTC),
					10.0,
				}},
			},
			{
				Name:    "cpu_usage_idle",
				Tags:    map[string]string{"cpu": "cpu1"},
				Columns: []string{"time", "sum"},
				Values: [][]interface{}{[]interface{}{
					time.Date(1971, 1, 1, 0, 0, 28, 0, time.UTC),
					10.0,
				}},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_SimpleMR", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(30 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_SimpleMR")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}

func TestBatch_Join(t *testing.T) {

	var script = `
var cpu0 = batch
	.query('''
		SELECT mean("value")
		FROM "telegraf"."default".cpu_usage_idle
		WHERE "cpu" = 'cpu0'
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))

var cpu1 = batch
	.query('''
		SELECT mean("value")
		FROM "telegraf"."default".cpu_usage_idle
		WHERE "cpu" = 'cpu1'
''')
		.period(10s)
		.every(10s)
		.groupBy(time(2s))

cpu0.join(cpu1)
	.as('cpu0', 'cpu1')
	.mapReduce(influxql.count('cpu0.mean'))
	.window()
		.period(20s)
		.every(20s)
	.mapReduce(influxql.sum('count'))
	.httpOut('TestBatch_Join')
`

	er := kapacitor.Result{
		Series: imodels.Rows{
			{
				Name:    "cpu_usage_idle",
				Columns: []string{"time", "sum"},
				Values: [][]interface{}{[]interface{}{
					time.Date(1971, 1, 1, 0, 0, 28, 0, time.UTC),
					10.0,
				}},
			},
		},
	}

	clock, et, errCh, tm := testBatcher(t, "TestBatch_Join", script)
	defer tm.Close()

	// Move time forward
	clock.Set(clock.Zero().Add(30 * time.Second))
	// Wait till the replay has finished
	if e := <-errCh; e != nil {
		t.Error(e)
	}
	// Wait till the task is finished
	if e := et.Err(); e != nil {
		t.Error(e)
	}

	// Get the result
	output, err := et.GetOutput("TestBatch_Join")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(output.Endpoint())
	if err != nil {
		t.Fatal(err)
	}

	// Assert we got the expected result
	result := kapacitor.ResultFromJSON(resp.Body)
	if eq, msg := compareResults(er, result); !eq {
		t.Error(msg)
	}
}

// Helper test function for batcher
func testBatcher(t *testing.T, name, script string) (clock.Setter, *kapacitor.ExecutingTask, <-chan error, *kapacitor.TaskMaster) {
	if testing.Verbose() {
		wlog.SetLevel(wlog.DEBUG)
	} else {
		wlog.SetLevel(wlog.OFF)
	}

	// Create task
	task, err := kapacitor.NewBatcher(name, script, dbrps)
	if err != nil {
		t.Fatal(err)
	}

	// Load test data
	var allData []io.ReadCloser
	var data io.ReadCloser
	for i := 0; err == nil; {
		f := fmt.Sprintf("%s.%d.brpl", name, i)
		data, err = os.Open(path.Join("data", f))
		if err == nil {
			allData = append(allData, data)
			i++
		}
	}
	if len(allData) == 0 {
		t.Fatal("could not find any data files for", name)
	}

	// Use 1971 so that we don't get true negatives on Epoch 0 collisions
	c := clock.New(time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC))
	r := kapacitor.NewReplay(c)

	// Create a new execution env
	tm := kapacitor.NewTaskMaster(logService)
	tm.HTTPDService = httpService
	tm.Open()

	//Start the task
	et, err := tm.StartTask(task)
	if err != nil {
		t.Fatal(err)
	}

	// Replay test data to executor
	batches := tm.BatchCollectors(name)
	errCh := r.ReplayBatch(allData, batches, false)

	t.Log(string(et.Task.Dot()))
	return r.Setter, et, errCh, tm
}