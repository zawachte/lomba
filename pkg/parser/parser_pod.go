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

func (cp *containerParser) Parse(input string) (time.Time, model.LabelSet, error) {
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

	labels := model.LabelSet{}
	labels["log_type"] = "journal"

	return tm, labels, nil
}
