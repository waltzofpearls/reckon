package prom

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Metric struct {
	Name   string
	Labels map[string]string
	Values SamplePairs

	labelNames []string
}

type SamplePair struct {
	Time  time.Time
	Value float64
}

type SamplePairs []SamplePair

func (m *Metric) LabelNames() []string {
	if len(m.labelNames) == 0 && len(m.Labels) > 0 {
		for name := range m.Labels {
			m.labelNames = append(m.labelNames, name)
		}
	}
	return m.labelNames
}

func (m *Metric) LabelValuesFor(names []string) []string {
	var values []string
	for _, name := range names {
		values = append(values, m.Labels[name])
	}
	return values
}

func (m *Metric) Key() string {
	lpStrings := make([]string, 0, len(m.labelNames))
	for _, lp := range m.labelNames {
		lpStrings = append(
			lpStrings,
			fmt.Sprintf("%s=%q", lp, m.Labels[lp]),
		)
	}
	return fmt.Sprintf(
		"Metric{name: %q, labels: {%s}}",
		m.Name,
		strings.Join(lpStrings, ","),
	)
}

func (m *Metric) Start() time.Time {
	if len(m.Values) == 0 {
		return time.Time{}
	}
	return m.Values[0].Time
}

func (m *Metric) End() time.Time {
	if len(m.Values) == 0 {
		return time.Time{}
	}
	return m.Values[len(m.Values)-1].Time
}

func (m *Metric) Equal(to Metric) bool {
	return m.Name == to.Name && twoStringMapsAreEqual(m.Labels, to.Labels)
}

func (m *Metric) Append(data Metric, rollingWindow time.Duration) error {
	if !m.Equal(data) {
		return errors.New("new and original metrics are not equal")
	}

	existingDataStart := m.Start()
	existingDataEnd := m.End()
	newDataStart := data.Start()
	newDataEnd := data.End()

	if (newDataStart.After(existingDataStart) || newDataStart.Equal(existingDataStart)) &&
		(newDataEnd.Before(existingDataEnd) || newDataEnd.Equal(existingDataEnd)) {
		// noop if the entire new data set is a part of existing data set
		return nil
	}
	if (existingDataStart.After(newDataStart) || existingDataStart.Equal(newDataStart)) &&
		(existingDataEnd.Before(newDataEnd) || existingDataEnd.Equal(newDataEnd)) {
		// noop if the entire existing data set is a part of new data set
		return nil
	}

	existingDataValues := make(SamplePairs, len(m.Values))
	newDataValues := make(SamplePairs, len(data.Values))
	updatedValues := make(SamplePairs, 0)

	copy(existingDataValues, m.Values)
	copy(newDataValues, data.Values)

	if existingDataEnd.Before(newDataStart) || existingDataEnd.Equal(newDataStart) {
		// no intersection between existing and new data, and new data set is newer than existing data set
		updatedValues = append(existingDataValues, newDataValues...)
	} else if newDataEnd.Before(existingDataStart) || newDataEnd.Equal(existingDataStart) {
		// no intersection between existing and new data, and existing data set is newer than new data set
		updatedValues = append(newDataValues, existingDataValues...)
	} else {
		// has intersection between existing and new data, remove overlapping data and concatenate the rest
		var sliceFrom int
		for i, pair := range newDataValues {
			if pair.Time.After(existingDataEnd) {
				sliceFrom = i
				break
			}
		}
		updatedValues = append(existingDataValues, newDataValues[sliceFrom:]...)
	}

	rollingWindowStart := updatedValues[len(updatedValues)-1].Time.Add(-rollingWindow)

	var sliceFrom int
	for i, pair := range updatedValues {
		if pair.Time.After(rollingWindowStart) || pair.Time.Equal(rollingWindowStart) {
			sliceFrom = i
			break
		}
	}

	m.Values = updatedValues[sliceFrom:]
	return nil
}

func (s SamplePair) IsZero() bool {
	return s.Time.IsZero() && s.Value == 0
}

func (s SamplePairs) Nearest(nearest time.Time) SamplePair {
	if len(s) == 0 {
		return SamplePair{}
	}
	current := s[0]
	for _, pair := range s {
		if math.Abs(float64(nearest.Sub(pair.Time))) < math.Abs(float64(nearest.Sub(current.Time))) {
			current = pair
		}
	}
	return current
}

func twoStringMapsAreEqual(a, b map[string]string) bool {
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if vv, exists := b[k]; !exists || v != vv {
			return false
		}
	}
	return true
}
