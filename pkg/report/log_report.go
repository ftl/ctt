package report

import (
	"log"

	"github.com/ftl/ctt/pkg/trainer"
)

type LogReport struct{}

func NewLogReport() *LogReport {
	return &LogReport{}
}

func (r *LogReport) Add(attempt trainer.Attempt) {
	log.Printf("attempt: %+v", attempt)
}
