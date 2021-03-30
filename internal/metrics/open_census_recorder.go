package metrics

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

const AddItemLatencyMeasureName = "cart_add_item_latency"

type Recorder interface {
	Record(start int64, measureName string)
}

type OpenCensusRecorder struct {
	latencyMeasures []*stats.Float64Measure
	views           []*view.View
}

func NewOpenCensusRecorder() (OpenCensusRecorder, error) {
	addItemLatencyS := stats.Float64(AddItemLatencyMeasureName, "The latency in seconds for the complete validation process", "s")

	views := []*view.View{
		{
			Name:        "cart/latency",
			Description: "Latency of adding an item",
			Measure:     addItemLatencyS,
			Aggregation: view.Distribution(0, 10, 50, 100, 200, 400, 800, 1000, 1400, 2000, 5000, 10000, 15000),
		},
		{
			Name:        "cart/counts",
			Description: "Counts of items added",
			Measure:     addItemLatencyS,
			Aggregation: view.Count(),
		},
	}

	if err := view.Register(views...); err != nil {
		return OpenCensusRecorder{}, err
	}

	measures := []*stats.Float64Measure{addItemLatencyS}
	return OpenCensusRecorder{
		latencyMeasures: measures,
		views:           views,
	}, nil
}

func (t OpenCensusRecorder) Record(start int64, measureName string) {
	s := time.Unix(start, 0)

	totalTime := time.Since(s)

	measure := t.measureByName(measureName)
	if measure == nil {
		return
	}

	go stats.Record(context.Background(), measure.M(totalTime.Seconds()))
}

func (t OpenCensusRecorder) measureByName(measurementName string) *stats.Float64Measure {
	for _, measure := range t.latencyMeasures {
		if measure.Name() == measurementName {
			return measure
		}
	}

	return nil
}
