package krong

import (
	"errors"
	"expvar"
	"fmt"
	"strings"
	"sync"

	"github.com/karasz/krong/parser"
	"github.com/robfig/cron/v3"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewInt("scheduler_started")

	// ErrScheduleParse is the error returned when the schedule parsing fails.
	ErrScheduleParse = errors.New("can't parse job schedule")
)

// Scheduler scheduler instance
type Scheduler struct {
	mu          sync.RWMutex
	Cron        *cron.Cron
	started     bool
	EntryJobMap sync.Map
}

// NewScheduler creates a new Scheduler instance
func NewScheduler() *Scheduler {
	schedulerStarted.Set(0)
	return &Scheduler{
		Cron:        nil,
		started:     false,
		EntryJobMap: sync.Map{},
	}
}

// Start the cron scheduler, adding its corresponding jobs and
// executing them on time.
func (s *Scheduler) Start(jobs []*Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Cron != nil {
		// Creating a new cron is risky if not nil because the previous invocation is dirty
		return fmt.Errorf("cron is already configured, can not start scheduler")
	}

	s.Cron = cron.New(cron.WithParser(parser.NewParser()))

	for _, job := range jobs {
		if err := s.AddJob(job); err != nil {
			return err
		}
	}
	s.Cron.Start()
	s.started = true
	schedulerStarted.Set(1)

	return nil
}

// Stop stops the scheduler effectively not running any job.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		s.Cron.Stop()
		s.started = false
		// expvars
		cronInspect.Do(func(kv expvar.KeyValue) {
			kv.Value = nil
		})
	}
	schedulerStarted.Set(0)
}

// Restart the scheduler
func (s *Scheduler) Restart(jobs []*Job) {
	s.Stop()
	s.ClearCron()
	s.Start(jobs)
}

// Clear cron separately, this can only be called when agent will be stop.
func (s *Scheduler) ClearCron() {
	s.Cron = nil
}

// Started will safely return if the scheduler is started or not
func (s *Scheduler) Started() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.started
}

// GetEntry returns a scheduler entry from a snapshot in
// the current time, and whether or not the entry was found.
func (s *Scheduler) GetEntry(jobName string) (cron.Entry, bool) {
	for _, e := range s.Cron.Entries() {
		j, _ := e.Job.(*Job)
		if j.Name == jobName {
			return e, true
		}
	}
	return cron.Entry{}, false
}

// AddJob Adds a job to the cron scheduler
func (s *Scheduler) AddJob(job *Job) error {
	// Check if the job is already set and remove it if exists
	if _, ok := s.EntryJobMap.Load(job.Name); ok {
		s.RemoveJob(job)
	}

	// If Timezone is set on the job, and not explicitly in its schedule,
	// AND its not a descriptor (that don't support timezones), add the
	// timezone to the schedule so robfig/cron knows about it.
	schedule := job.Schedule
	if job.Timezone != "" &&
		!strings.HasPrefix(schedule, "@") &&
		!strings.HasPrefix(schedule, "TZ=") &&
		!strings.HasPrefix(schedule, "CRON_TZ=") {
		schedule = "CRON_TZ=" + job.Timezone + " " + schedule
	}

	id, err := s.Cron.AddJob(schedule, job)
	if err != nil {
		return err
	}
	s.EntryJobMap.Store(job.Name, id)

	cronInspect.Set(job.Name, job)

	return nil
}

// RemoveJob removes a job from the cron scheduler
func (s *Scheduler) RemoveJob(job *Job) {

	if v, ok := s.EntryJobMap.Load(job.Name); ok {
		s.Cron.Remove(v.(cron.EntryID))
		s.EntryJobMap.Delete(job.Name)

		cronInspect.Delete(job.Name)
	}
}
