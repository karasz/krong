package parser

import (
	"github.com/robfig/cron/v3"
)

// Parser is a cron parser
type Parser struct {
	parser cron.Parser
}

// NewParser creates an Parser instance
func NewParser() cron.ScheduleParser {
	return Parser{cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)}
}

// Parse parses a cron schedule specification. It accepts the cron spec with
// mandatory seconds parameter.
func (p Parser) Parse(spec string) (cron.Schedule, error) {
	switch spec {
	case "@yearly", "@annually":
		spec = "0 0 0 1 1 * *"
	case "@monthly":
		spec = "0 0 0 1 * * *"
	case "@weekly":
		spec = "0 0 0 * * 0 *"
	case "@daily":
		spec = "0 0 0 * * * *"
	case "@hourly":
		spec = "0 0 * * * * *"
	case "@minutely":
		spec = "0 * * * * *"
	}

	return p.parser.Parse(spec)
}

var standaloneParser = NewParser()

// Parse parses a cron schedule.
func Parse(spec string) (cron.Schedule, error) {
	return standaloneParser.Parse(spec)
}
