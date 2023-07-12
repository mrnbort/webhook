package task

// Task represents a single task
type Task struct {
	ID      string `yaml:"id"`
	Auth    string `yaml:"auth"`
	Command string `yaml:"command"`
}

// TaskList represents a slice of tasks
type TaskList struct {
	Tasks []Task `yaml:"tasks"`
}
