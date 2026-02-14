package model

import "time"

type PhaseStatus string

const (
	StatusTodo       PhaseStatus = "todo"
	StatusInProgress PhaseStatus = "in_progress"
	StatusDone       PhaseStatus = "done"
	StatusBlocked    PhaseStatus = "blocked"
)

type Phase struct {
	ID          int         `json:"id" yaml:"id"`
	Day         int         `json:"day" yaml:"day"`
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description" yaml:"description"`
	Status      PhaseStatus `json:"status" yaml:"status"`
	DependsOn   []int       `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	CreatedAt   time.Time   `json:"created_at" yaml:"-"`
	UpdatedAt   time.Time   `json:"updated_at" yaml:"-"`
}

type Task struct {
	ID        int       `json:"id" yaml:"id"`
	PhaseID   int       `json:"phase_id" yaml:"phase_id"`
	Name      string    `json:"name" yaml:"name"`
	Done      bool      `json:"done" yaml:"done"`
	CreatedAt time.Time `json:"created_at" yaml:"-"`
	UpdatedAt time.Time `json:"updated_at" yaml:"-"`
}

type PlanImport struct {
	ProjectName string  `json:"project_name" yaml:"project_name"`
	Phases      []Phase `json:"phases" yaml:"phases"`
}

type ProjectStats struct {
	TotalPhases      int
	TodoPhases       int
	InProgressPhases int
	DonePhases       int
	BlockedPhases    int
	TotalTasks       int
	DoneTasks        int
	PercentComplete  float64
}
