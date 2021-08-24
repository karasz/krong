package krong

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/karasz/krong/isset"
	"github.com/karasz/krong/parser"
)

// Job is a cron job
type Job struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	DisplayName  string     `json:"display_name"`
	Timezone     string     `json:"timezone"`
	Schedule     string     `json:"schedule"`
	Owner        int        `json:"owner"`
	SuccessCount int        `json:"success_count"`
	ErrorCount   int        `json:"error_count"`
	LastSuccess  time.Time  `json:"last_success"`
	LastError    time.Time  `json:"last_error"`
	Disabled     isset.Bool `json:"disabled"`
	Retries      int        `json:"retries"`
	Concurrency  string     `json:"concurency"`
	Status       string     `json:"status"`
	Next         time.Time  `json:"next"`
	Ephemeral    isset.Bool `json:"ephemeral"`
	ExpiresAt    time.Time  `json:"expires_at"`
	WebHook      WebHook    `json:"webhook"`
	Agent        Agent      `json:"agent"`
	Type         string     `json:"type"`
}

func NewJob() *Job {
	return &Job{
		ID:           0,
		Name:         "",
		DisplayName:  "",
		Timezone:     time.Local.String(),
		Schedule:     "",
		Owner:        0,
		SuccessCount: 0,
		ErrorCount:   0,
		LastSuccess:  time.Time{},
		LastError:    time.Time{},
		Disabled:     isset.Bool{},
		Retries:      0,
		Concurrency:  "",
		Status:       "",
		Next:         time.Time{},
		Ephemeral:    isset.Bool{},
		ExpiresAt:    time.Time{},
		WebHook:      WebHook{},
		Agent:        Agent{},
		Type:         "",
	}
}
func (j *Job) Run() {
	switch j.Type {
	case "webhook":
	case "agent":
	default:
		fmt.Println(time.Now())
	}
}
func (j *Job) IsValid() error {
	var isLetter = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	var isConcurencyAllowForbid = regexp.MustCompile(`^allow$|^forbid$`).MatchString

	if j.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if !isLetter(j.Name) {
		return fmt.Errorf("%s name contains illegal characters", j.Name)
	}

	if j.Schedule != "" {
		if _, err := parser.Parse(j.Schedule); err != nil {
			return fmt.Errorf("%s", err)
		}
	}
	if !isConcurencyAllowForbid(j.Concurrency) && j.Concurrency != "" {
		return fmt.Errorf("%s concurency is illegal", j.Concurrency)
	}
	if _, err := time.LoadLocation(j.Timezone); err != nil {
		return err
	}
	return nil
}

func (j *Job) GetNext() (time.Time, error) {
	if j.Schedule != "" {
		s, err := parser.Parse(j.Schedule)
		if err != nil {
			return time.Time{}, err
		}
		return s.Next(time.Now()), nil
	}

	return time.Time{}, nil
}

func (j *Job) String() string {
	s, err := json.Marshal(j)
	if err != nil {
		return err.Error()
	}
	return string(s)
}
