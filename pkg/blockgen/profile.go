package blockgen

import (
	"bytes"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-io/thanosbench/pkg/seriesgen"
	"gopkg.in/yaml.v2"
	"io"
	"text/template"
	"time"
)

type ProfileSpecification struct {
	Type            GenType                   `yaml:"type"`
	Characteristics seriesgen.Characteristics `yaml:"characteristics"`
	Labels          labels.Labels             `yaml:"labels"`
}

type Profile struct {
	Ranges           []int                `yaml:"ranges"`
	RolloutInterval  int                  `yaml:"rolloutInterval"`
	Targets          int                  `yaml:"targets"`
	MetricsPerTarget int                  `yaml:"metricsPerTarget"`
	Specification    ProfileSpecification `yaml:"specification"`
}

func (p Profile) TimeRanges() (result []time.Duration) {
	for _, r := range p.Ranges {
		result = append(result, time.Duration(r)*time.Hour)
	}
	return result
}

func (p Profile) ChurnInterval() time.Duration {
	return time.Duration(p.RolloutInterval) * time.Hour
}

func (p Profile) seriesLabels(index int) labels.Labels {
	data := map[string]any{"index": index}
	seriesLabels := labels.Labels{}
	for _, label := range p.Specification.Labels {
		tmpl, err := template.New(label.Name).Parse(label.Value)
		if err != nil {
			seriesLabels = append(seriesLabels, label)
		} else {
			result := bytes.Buffer{}
			err := tmpl.Execute(&result, data)
			if err != nil {
				seriesLabels = append(seriesLabels, label)
			} else {
				seriesLabels = append(seriesLabels, labels.Label{Name: label.Name, Value: result.String()})
			}
		}
	}
	return seriesLabels
}

func (p Profile) SeriesSpec() SeriesSpecFn {
	return func(targets int, index int, nextRolloutTime string) SeriesSpec {
		return SeriesSpec{
			Targets:         p.Targets,
			Type:            p.Specification.Type,
			Characteristics: p.Specification.Characteristics,
			Labels:          p.seriesLabels(index),
		}
	}
}

func NewProfile(reader io.Reader, profile *Profile) error {
	decoder := yaml.NewDecoder(reader)
	return decoder.Decode(&profile)
}
