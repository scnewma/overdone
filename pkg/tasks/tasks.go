package tasks

type Task struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Content   string `json:"content"`
}

func (t *Task) complete() {
	t.Completed = true
}

type Repository interface {
	Get(int) (Task, error)
	All() ([]Task, error)
	Save(Task) error
	NextID() int
}
