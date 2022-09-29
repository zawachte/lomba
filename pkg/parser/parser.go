package parser

import (
	"time"

	"github.com/prometheus/common/model"
)

type Parser interface {
	Parse(string) (time.Time, model.LabelSet, error)
}
