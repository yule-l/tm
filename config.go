package tm

import (
	"errors"
)

const DefaultMaxRetries = 5

// Config is a configuration for task manager
type Config struct {
	// Force will truncate file with tasks results
	Force bool

	// FilePath must contain tasks file path
	FilePath string

	// MaxRetries maximum number of attempts to complete the task
	// If number of attempts reaches MaxRetries, task will be marked as completed with errors
	MaxRetries uint8
}

// NewDefaultConfig returns default config
func NewDefaultConfig(filePath string) *Config {
	return &Config{
		Force:      false,
		FilePath:   filePath,
		MaxRetries: DefaultMaxRetries,
	}
}

var (
	ErrEmptyFilePath    = errors.New("FilePath can't be empty string")
	ErrMaxRetriesIsZero = errors.New("MaxRetries can't be less than 1")
)

func (c *Config) validate() error {
	if c.FilePath == "" {
		return ErrEmptyFilePath
	}
	if c.MaxRetries == 0 {
		return ErrMaxRetriesIsZero
	}
	return nil
}
