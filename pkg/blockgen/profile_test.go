package blockgen

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
	"github.com/thanos-io/thanosbench/pkg/seriesgen"
	"strings"
	"testing"
	"time"
)

func TestNewProfile(t *testing.T) {
	t.Run("should be able to create new profile", func(t *testing.T) {
		profileStr := `
---
ranges: [2,2]
rolloutInterval: 1
targets: 1
metricsPerTarget: 1
specification:
  type: 'COUNTER'
  characteristics:
    jitter: 1
    scrapeInterval: 1
    changeInterval: 1
    max: 5
    min: 1
  labels:
    one: 1
    two: 2
`
		expected := Profile{
			Ranges:           []int{2, 2},
			RolloutInterval:  1,
			Targets:          1,
			MetricsPerTarget: 1,
			Specification: ProfileSpecification{
				Type: "COUNTER",
				Characteristics: seriesgen.Characteristics{
					Jitter:         1,
					ScrapeInterval: 1,
					ChangeInterval: 1,
					Max:            5,
					Min:            1,
				},
				Labels: labels.Labels{
					{Name: "one", Value: "1"},
					{Name: "two", Value: "2"},
				},
			},
		}

		reader := strings.NewReader(profileStr)
		actual := Profile{}
		err := NewProfile(reader, &actual)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestProfile_TimeRanges(t *testing.T) {
	profile := Profile{Ranges: []int{4, 5}}
	expected := []time.Duration{4 * time.Hour, 5 * time.Hour}

	actual := profile.TimeRanges()

	assert.Equal(t, expected, actual)
}

func TestProfile_ChurnInterval(t *testing.T) {
	profile := Profile{RolloutInterval: 2}
	expected := 2 * time.Hour

	actual := profile.ChurnInterval()
	assert.Equal(t, expected, actual)
}

func TestProfile_SeriesSpec(t *testing.T) {
	t.Run("it should return series spec with rendered value", func(t *testing.T) {
		profile := Profile{
			Ranges:           []int{2, 2},
			RolloutInterval:  1,
			Targets:          1,
			MetricsPerTarget: 1,
			Specification: ProfileSpecification{
				Type: "COUNTER",
				Characteristics: seriesgen.Characteristics{
					Jitter:         1,
					ScrapeInterval: 1,
					ChangeInterval: 1,
					Max:            5,
					Min:            1,
				},
				Labels: labels.Labels{
					{Name: "name", Value: "app-{{.index}}"},
				},
			},
		}
		expected := SeriesSpec{
			Targets: 1,
			Type:    "COUNTER",
			Characteristics: seriesgen.Characteristics{
				Jitter:         1,
				ScrapeInterval: 1,
				ChangeInterval: 1,
				Max:            5,
				Min:            1,
			},
			Labels: labels.Labels{
				{Name: "name", Value: "app-1"},
			},
		}

		seriesSpecFn := profile.SeriesSpec()
		actual := seriesSpecFn(5, 1, "next")

		assert.Equal(t, expected, actual)
	})
}
