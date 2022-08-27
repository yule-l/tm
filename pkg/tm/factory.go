package tm

type Factory struct {
}

// ParamsIn is a configuration for doers
type ParamsIn struct {
	// Force will truncate file with tasks results
	Force bool

	// TasksFilePath must contain tasks file path
	TasksFilePath string

	// MaxRetries maximum number of attempts to complete the task
	// If number of attempts reaches MaxRetries, task will be marked as completed with errors
	MaxRetries uint8

	// Parallel will enable parallel mode
	Parallel bool
}

func (f *Factory) NewDoer(in *ParamsIn) (Doer, error) {
	taskManager, err := NewTasksManager(Config{
		Force:      in.Force,
		FilePath:   in.TasksFilePath,
		MaxRetries: in.MaxRetries,
	})
	if err != nil {
		return nil, err
	}
	if in.Parallel {
		return NewDefaultParallelDoer(taskManager), nil
	}
	return NewDefaultOrderedDoer(taskManager), nil
}
