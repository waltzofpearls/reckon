package prom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricLabelNames(t *testing.T) {
	tests := []struct {
		name  string
		given *Metric
		want  []string
	}{
		{
			name:  "no labels",
			given: &Metric{},
			want:  []string{},
		},
		{
			name: "labelNames already set",
			given: &Metric{
				Labels: map[string]string{
					"label1": "value1",
				},
				labelNames: []string{"label1"},
			},
			want: []string{"label1"},
		},
		{
			name: "no labels",
			given: &Metric{
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
					"label3": "value3",
				},
			},
			want: []string{"label1", "label2", "label3"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.LabelNames()
			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestMetricLabelValuesFor(t *testing.T) {
	tests := []struct {
		name        string
		givenMetric *Metric
		givenNames  []string
		want        []string
	}{
		{
			name: "empty names given",
			givenMetric: &Metric{
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
					"label3": "value3",
				},
			},
			givenNames: []string{},
			want:       []string{},
		},
		{
			name:        "no labels",
			givenMetric: &Metric{},
			givenNames:  []string{"label1", "label2", "label3"},
			want:        []string{"", "", ""},
		},
		{
			name: "given a subset of names",
			givenMetric: &Metric{
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
					"label3": "value3",
				},
			},
			givenNames: []string{"label1"},
			want:       []string{"value1"},
		},
		{
			name: "happy path",
			givenMetric: &Metric{
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
					"label3": "value3",
				},
			},
			givenNames: []string{"label1", "label2", "label3"},
			want:       []string{"value1", "value2", "value3"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.givenMetric.LabelValuesFor(test.givenNames)
			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestMetricKey(t *testing.T) {
	tests := []struct {
		name  string
		given *Metric
		want  string
	}{
		{
			name:  "no name no labels",
			given: &Metric{},
			want:  `Metric{name: "", labels: {}}`,
		},
		{
			name: "has name but no labels",
			given: &Metric{
				Name: "test_metric_name",
			},
			want: `Metric{name: "test_metric_name", labels: {}}`,
		},
		{
			name: "has both name and labels",
			given: &Metric{
				Name: "test_metric_name",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
				labelNames: []string{"label1", "label2"},
			},
			want: `Metric{name: "test_metric_name", labels: {label1="value1",label2="value2"}}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.Key()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricStart(t *testing.T) {
	now := time.Now()
	twoSecondsAgo := now.Add(-2 * time.Second)
	oneSecondAgo := now.Add(-1 * time.Second)

	tests := []struct {
		name  string
		given *Metric
		want  time.Time
	}{
		{
			name:  "no values",
			given: &Metric{},
			want:  time.Time{},
		},
		{
			name: "has values",
			given: &Metric{
				Values: SamplePairs{
					{Time: twoSecondsAgo, Value: 25.0},
					{Time: oneSecondAgo, Value: 26.0},
				},
			},
			want: twoSecondsAgo,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.Start()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricEnd(t *testing.T) {
	now := time.Now()
	twoSecondsAgo := now.Add(-2 * time.Second)
	oneSecondAgo := now.Add(-1 * time.Second)

	tests := []struct {
		name  string
		given *Metric
		want  time.Time
	}{
		{
			name:  "no values",
			given: &Metric{},
			want:  time.Time{},
		},
		{
			name: "has values",
			given: &Metric{
				Values: SamplePairs{
					{Time: twoSecondsAgo, Value: 25.0},
					{Time: oneSecondAgo, Value: 26.0},
				},
			},
			want: oneSecondAgo,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.End()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricEqual(t *testing.T) {
	tests := []struct {
		name       string
		givenLeft  *Metric
		givenRight Metric
		want       bool
	}{
		{
			name: "not equal - different name",
			givenLeft: &Metric{
				Name: "test_metric_1",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			givenRight: Metric{
				Name: "test_metric_2",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			want: false,
		},
		{
			name: "not equal - different labels",
			givenLeft: &Metric{
				Name: "test_metric_1",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			givenRight: Metric{
				Name: "test_metric_1",
				Labels: map[string]string{
					"label3": "value3",
					"label4": "value4",
				},
			},
			want: false,
		},
		{
			name:       "equal - both empty",
			givenLeft:  &Metric{},
			givenRight: Metric{},
			want:       true,
		},
		{
			name: "equal - both not empty",
			givenLeft: &Metric{
				Name: "test_metric_1",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			givenRight: Metric{
				Name: "test_metric_1",
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			want: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.givenLeft.Equal(test.givenRight)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricAppend(t *testing.T) {
	now := time.Now()
	fiveSecondsAgo := now.Add(-5 * time.Second)
	tenSecondsAgo := now.Add(-10 * time.Second)
	fifteenSecondsAgo := now.Add(-15 * time.Second)
	twentySecondsAgo := now.Add(-20 * time.Second)
	twentyFiveSecondsAgo := now.Add(-25 * time.Second)
	thirtySecondsAgo := now.Add(-30 * time.Second)
	thirtyFiveSecondsAgo := now.Add(-35 * time.Second)

	metricName := "test_metric_1"
	metricLabels := map[string]string{
		"label1": "value1",
		"label2": "value2",
		"label3": "value3",
	}
	metricValues := SamplePairs{
		{Time: thirtySecondsAgo, Value: 1300.0},
		{Time: twentyFiveSecondsAgo, Value: 1250.0},
		{Time: twentySecondsAgo, Value: 1200.0},
		{Time: fifteenSecondsAgo, Value: 1150.0},
	}

	tests := []struct {
		name         string
		givenMetric  *Metric
		givenNewData Metric
		givenWindow  time.Duration
		wantError    error
		wantValues   SamplePairs
	}{
		{
			name: "not appending the same metric",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name: "a_different_metric",
				Labels: map[string]string{
					"something": "else",
				},
				Values: SamplePairs{
					{Time: now, Value: 1000.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantError:   errCannotAppendDifferentMetric,
			wantValues:  metricValues,
		},
		{
			name: "new data is a part of existing data",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: SamplePairs{
					{Time: twentyFiveSecondsAgo, Value: 1250.0},
					{Time: twentySecondsAgo, Value: 1200.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantValues:  metricValues,
		},
		{
			name: "existing data is a part of new data",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: SamplePairs{
					{Time: thirtyFiveSecondsAgo, Value: 1350.0},
					{Time: thirtySecondsAgo, Value: 1300.0},
					{Time: twentyFiveSecondsAgo, Value: 1250.0},
					{Time: twentySecondsAgo, Value: 1200.0},
					{Time: fifteenSecondsAgo, Value: 1150.0},
					{Time: tenSecondsAgo, Value: 1100.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantValues: SamplePairs{
				{Time: thirtySecondsAgo, Value: 1300.0},
				{Time: twentyFiveSecondsAgo, Value: 1250.0},
				{Time: twentySecondsAgo, Value: 1200.0},
				{Time: fifteenSecondsAgo, Value: 1150.0},
				{Time: tenSecondsAgo, Value: 1100.0},
			},
		},
		{
			name: "no intersection and new data is newer",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: SamplePairs{
					{Time: tenSecondsAgo, Value: 1100.0},
					{Time: fiveSecondsAgo, Value: 1050.0},
					{Time: now, Value: 1000.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantValues: SamplePairs{
				{Time: twentySecondsAgo, Value: 1200.0},
				{Time: fifteenSecondsAgo, Value: 1150.0},
				{Time: tenSecondsAgo, Value: 1100.0},
				{Time: fiveSecondsAgo, Value: 1050.0},
				{Time: now, Value: 1000.0},
			},
		},
		{
			name: "no intersection and existing data is newer",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: SamplePairs{
					{Time: thirtyFiveSecondsAgo, Value: 1350.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantValues: SamplePairs{
				{Time: thirtyFiveSecondsAgo, Value: 1350.0},
				{Time: thirtySecondsAgo, Value: 1300.0},
				{Time: twentyFiveSecondsAgo, Value: 1250.0},
				{Time: twentySecondsAgo, Value: 1200.0},
				{Time: fifteenSecondsAgo, Value: 1150.0},
			},
		},
		{
			name: "has intersection between existing and new data",
			givenMetric: &Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: metricValues,
			},
			givenNewData: Metric{
				Name:   metricName,
				Labels: metricLabels,
				Values: SamplePairs{
					{Time: twentySecondsAgo, Value: 1200.0},
					{Time: fifteenSecondsAgo, Value: 1150.0},
					{Time: tenSecondsAgo, Value: 1100.0},
					{Time: fiveSecondsAgo, Value: 1050.0},
				},
			},
			givenWindow: 20 * time.Second,
			wantValues: SamplePairs{
				{Time: twentyFiveSecondsAgo, Value: 1250.0},
				{Time: twentySecondsAgo, Value: 1200.0},
				{Time: fifteenSecondsAgo, Value: 1150.0},
				{Time: tenSecondsAgo, Value: 1100.0},
				{Time: fiveSecondsAgo, Value: 1050.0},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.givenMetric.Append(test.givenNewData, test.givenWindow)
			assert.Equal(t, test.wantError, err)
			assert.ElementsMatch(t, test.wantValues, test.givenMetric.Values)
		})
	}
}

func TestSamplePairIsZero(t *testing.T) {
	tests := []struct {
		name  string
		given SamplePair
		want  bool
	}{
		{
			name:  "zero sample pair",
			given: SamplePair{},
			want:  true,
		},
		{
			name: "not a zero sample pair",
			given: SamplePair{
				Time:  time.Now(),
				Value: 25.0,
			},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.IsZero()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestSamplePairsNearest(t *testing.T) {
	now := time.Now()

	twentySecondsAgo := now.Add(-20 * time.Second)
	fifteenSecondsAgo := now.Add(-15 * time.Second)
	tenSecondsAgo := now.Add(-10 * time.Second)
	fiveSecondAgo := now.Add(-5 * time.Second)

	tests := []struct {
		name       string
		givenPairs SamplePairs
		givenTime  time.Time
		want       SamplePair
	}{
		{
			name:       "empty forecasts",
			givenPairs: SamplePairs{},
			givenTime:  now.Add(-9 * time.Second),
			want:       SamplePair{},
		},
		{
			name: "only 1 item in forecasts",
			givenPairs: SamplePairs{
				{
					Time:  fiveSecondAgo,
					Value: 25.0,
				},
			},
			givenTime: now.Add(-9 * time.Second),
			want: SamplePair{
				Time:  fiveSecondAgo,
				Value: 25.0,
			},
		},
		{
			name: "multiple forecasts",
			givenPairs: SamplePairs{
				{
					Time:  twentySecondsAgo,
					Value: 28.0,
				},
				{
					Time:  fifteenSecondsAgo,
					Value: 27.0,
				},
				{
					Time:  tenSecondsAgo,
					Value: 26.0,
				},
				{
					Time:  fiveSecondAgo,
					Value: 25.0,
				},
			},
			givenTime: now.Add(-9 * time.Second),
			want: SamplePair{
				Time:  tenSecondsAgo,
				Value: 26.0,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.givenPairs.Nearest(test.givenTime)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTwoStringMapsAreEqual(t *testing.T) {
	tests := []struct {
		name       string
		givenLeft  map[string]string
		givenRight map[string]string
		want       bool
	}{
		{
			name:       "both nil maps are equal",
			givenLeft:  nil,
			givenRight: nil,
			want:       true,
		},
		{
			name: "length of two maps are not equal",
			givenLeft: map[string]string{
				"label1": "value1",
			},
			givenRight: map[string]string{
				"label1": "value1",
				"label2": "value2",
			},
			want: false,
		},
		{
			name: "two maps have different items",
			givenLeft: map[string]string{
				"label1": "value1",
				"label2": "value2",
			},
			givenRight: map[string]string{
				"label3": "value3",
				"label4": "value4",
			},
			want: false,
		},
		{
			name: "two maps are equal",
			givenLeft: map[string]string{
				"label1": "value1",
				"label2": "value2",
				"label3": "value3",
			},
			givenRight: map[string]string{
				"label1": "value1",
				"label2": "value2",
				"label3": "value3",
			},
			want: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := twoStringMapsAreEqual(test.givenLeft, test.givenRight)
			assert.Equal(t, test.want, got)
		})
	}
}
