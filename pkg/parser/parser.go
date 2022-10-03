package parser

import (
	"time"

	"github.com/prometheus/common/model"
)

type Parser interface {
	Parse(string, map[string]string) (time.Time, model.LabelSet, error)
}
