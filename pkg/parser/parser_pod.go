package parser

import (
	"errors"
	"strings"
	"time"

	"github.com/prometheus/common/model"
)

type containerParser struct{}

func NewContainerParser() Parser {
	return &containerParser{}
}

func (cp *containerParser) Parse(input string, labels map[string]string) (time.Time, model.LabelSet, error) {
	split_string := strings.Split(input, " ")
	if len(split_string) < 1 {
		return time.Time{}, nil, errors.New("nil log line")
	}

	likelyTime := split_string[0]

	// Declaring layout constant
	// Calling Parse() method with its parameters
	tm, err := time.Parse(time.RFC3339Nano, likelyTime)
	if err != nil {
		return time.Time{}, nil, err
	}

	labelSet := model.LabelSet{}
	labelSet["log_type"] = "container"

	for k, v := range labels {
		labelSet[model.LabelName(k)] = model.LabelValue(v)
	}

	return tm, labelSet, nil
}
