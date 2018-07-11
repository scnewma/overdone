package todo

type Task struct {
	ID        int
	Completed bool
	Content   string
}

func (t *Task) complete() {
	t.Completed = true
}

type TaskRepository interface {
	Get(int) (*Task, error)
	All() ([]*Task, error)
	Save(*Task) error
	NextID() int
}
