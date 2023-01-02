package blockgen

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/stretchr/testify/assert"
	"github.com/thanos-io/thanos/pkg/block/metadata"
	"github.com/thanos-io/thanos/pkg/model"
	"github.com/thanos-io/thanosbench/pkg/seriesgen"
	"gopkg.in/alecthomas/kingpin.v2"
	"testing"
	"time"
)

func Test_profiles_continuous(t *testing.T) {
	t.Run("it should generate metrics using continuousAppMetricSeriesSpec", func(t *testing.T) {
		externalLabels := labels.Labels{}
		flag := kingpin.Flag("max-time", "").Default("30m")
		expected := BlockSpec{
			Meta: metadata.Meta{
				BlockMeta: tsdb.BlockMeta{
					MaxTime:    7200000,
					MinTime:    1,
					Compaction: tsdb.BlockMetaCompaction{Level: 1},
					Version:    1,
				},
				Thanos: metadata.Thanos{
					Labels:     externalLabels.Map(),
					Downsample: metadata.ThanosDownsample{Resolution: 0},
					Source:     "blockgen",
				},
			},
			Series: []SeriesSpec{
				{
					Targets: 1,
					Type:    Gauge,
					MinTime: 1,
					MaxTime: 7200000,
					Characteristics: seriesgen.Characteristics{
						Max:            200000000,
						Min:            10000000,
						Jitter:         30000000,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: "continuous_app_metric0"},
					},
				},
			},
		}

		apply := continuous([]time.Duration{2 * time.Hour}, 1, 1, continuousAppMetricSeriesSpec)
		err := apply(context.Background(), *model.TimeOrDuration(flag), externalLabels, func(actual BlockSpec) error {
			assert.Equal(t, expected, actual)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("it should generate metrics using custom series spec", func(t *testing.T) {
		externalLabels := labels.Labels{}
		flag := kingpin.Flag("max-time", "").Default("30m")
		expected := BlockSpec{
			Meta: metadata.Meta{
				BlockMeta: tsdb.BlockMeta{
					MaxTime:    7200000,
					MinTime:    1,
					Compaction: tsdb.BlockMetaCompaction{Level: 1},
					Version:    1,
				},
				Thanos: metadata.Thanos{
					Labels:     externalLabels.Map(),
					Downsample: metadata.ThanosDownsample{Resolution: 0},
					Source:     "blockgen",
				},
			},
			Series: []SeriesSpec{
				{
					Targets: 1,
					Type:    Counter,
					MinTime: 1,
					MaxTime: 7200000,
					Characteristics: seriesgen.Characteristics{
						Max:            1,
						Min:            1,
						Jitter:         3,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: "custom-0"},
					},
				},
			},
		}
		customSeriesSpec := func(targets int, index int, nextRolloutTime string) SeriesSpec {
			return SeriesSpec{
				Targets: targets,
				Type:    Counter,
				Characteristics: seriesgen.Characteristics{
					Max:            1,
					Min:            1,
					Jitter:         3,
					ScrapeInterval: 15 * time.Second,
					ChangeInterval: 1 * time.Hour,
				},
				Labels: labels.Labels{
					{Name: "__name__", Value: fmt.Sprintf("custom-%d", index)},
				},
			}
		}

		apply := continuous([]time.Duration{2 * time.Hour}, 1, 1, customSeriesSpec)
		err := apply(context.Background(), *model.TimeOrDuration(flag), externalLabels, func(actual BlockSpec) error {
			assert.Equal(t, expected, actual)
			return nil
		})
		assert.NoError(t, err)
	})
}

func Test_profiles_realisticK8s(t *testing.T) {
	t.Run("it should generate metrics using k8sAppMetricSeriesSpec", func(t *testing.T) {
		externalLabels := labels.Labels{}
		flag := kingpin.Flag("max-time", "").Default("30m")
		expected := BlockSpec{
			Meta: metadata.Meta{
				BlockMeta: tsdb.BlockMeta{
					MaxTime:    7200000,
					MinTime:    1,
					Compaction: tsdb.BlockMetaCompaction{Level: 1},
					Version:    1,
				},
				Thanos: metadata.Thanos{
					Labels:     externalLabels.Map(),
					Downsample: metadata.ThanosDownsample{Resolution: 0},
					Source:     "blockgen",
				},
			},
			Series: []SeriesSpec{
				{
					Targets: 1,
					Type:    Gauge,
					MinTime: 5400000,
					MaxTime: 7200000,
					Characteristics: seriesgen.Characteristics{
						Max:            200000000,
						Min:            10000000,
						Jitter:         30000000,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: fmt.Sprintf("k8s_app_metric0")},
						{Name: "next_rollout_time", Value: "1970-01-01 01:30:00 +0000 UTC"},
					},
				},
				{
					Targets: 1,
					Type:    Gauge,
					MinTime: 1800000,
					MaxTime: 5400000,
					Characteristics: seriesgen.Characteristics{
						Max:            200000000,
						Min:            10000000,
						Jitter:         30000000,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: fmt.Sprintf("k8s_app_metric0")},
						{Name: "next_rollout_time", Value: "1970-01-01 00:30:00 +0000 UTC"},
					},
				},
				{
					Targets: 1,
					Type:    Gauge,
					MinTime: 1,
					MaxTime: 1800000,
					Characteristics: seriesgen.Characteristics{
						Max:            200000000,
						Min:            10000000,
						Jitter:         30000000,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: fmt.Sprintf("k8s_app_metric0")},
						{Name: "next_rollout_time", Value: "1969-12-31 23:30:00 +0000 UTC"},
					},
				},
			},
		}

		apply := realisticK8s([]time.Duration{2 * time.Hour}, 1*time.Hour, 1, 1, k8sAppMetricSeriesSpec)
		err := apply(context.Background(), *model.TimeOrDuration(flag), externalLabels, func(actual BlockSpec) error {
			assert.Equal(t, expected, actual)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("it should generate metrics using custom series spec", func(t *testing.T) {
		externalLabels := labels.Labels{}
		flag := kingpin.Flag("max-time", "").Default("30m")
		expected := BlockSpec{
			Meta: metadata.Meta{
				BlockMeta: tsdb.BlockMeta{
					MaxTime:    7200000,
					MinTime:    1,
					Compaction: tsdb.BlockMetaCompaction{Level: 1},
					Version:    1,
				},
				Thanos: metadata.Thanos{
					Labels:     externalLabels.Map(),
					Downsample: metadata.ThanosDownsample{Resolution: 0},
					Source:     "blockgen",
				},
			},
			Series: []SeriesSpec{
				{
					Targets: 1,
					Type:    Counter,
					MinTime: 5400000,
					MaxTime: 7200000,
					Characteristics: seriesgen.Characteristics{
						Max:            1,
						Min:            1,
						Jitter:         3,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: "custom-0"},
						{Name: "rollout_time", Value: "1970-01-01 01:30:00 +0000 UTC"},
					},
				},
				{
					Targets: 1,
					Type:    Counter,
					MinTime: 1800000,
					MaxTime: 5400000,
					Characteristics: seriesgen.Characteristics{
						Max:            1,
						Min:            1,
						Jitter:         3,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: "custom-0"},
						{Name: "rollout_time", Value: "1970-01-01 00:30:00 +0000 UTC"},
					},
				},
				{
					Targets: 1,
					Type:    Counter,
					MinTime: 1,
					MaxTime: 1800000,
					Characteristics: seriesgen.Characteristics{
						Max:            1,
						Min:            1,
						Jitter:         3,
						ScrapeInterval: 15 * time.Second,
						ChangeInterval: 1 * time.Hour,
					},
					Labels: labels.Labels{
						{Name: "__name__", Value: "custom-0"},
						{Name: "rollout_time", Value: "1969-12-31 23:30:00 +0000 UTC"},
					},
				},
			},
		}
		customSeriesSpec := func(targets int, index int, nextRolloutTime string) SeriesSpec {
			return SeriesSpec{
				Targets: targets,
				Type:    Counter,
				Characteristics: seriesgen.Characteristics{
					Max:            1,
					Min:            1,
					Jitter:         3,
					ScrapeInterval: 15 * time.Second,
					ChangeInterval: 1 * time.Hour,
				},
				Labels: labels.Labels{
					{Name: "__name__", Value: fmt.Sprintf("custom-%d", index)},
					{Name: "rollout_time", Value: nextRolloutTime},
				},
			}
		}

		apply := realisticK8s([]time.Duration{2 * time.Hour}, 1*time.Hour, 1, 1, customSeriesSpec)
		err := apply(context.Background(), *model.TimeOrDuration(flag), externalLabels, func(actual BlockSpec) error {
			assert.Equal(t, expected, actual)
			return nil
		})
		assert.NoError(t, err)
	})
}
