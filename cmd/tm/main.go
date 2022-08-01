package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yule-l/tm"
)

func main() {
	ctx := context.Background()
	_ = newCli().ExecuteContext(ctx)
}

type doParams struct {
	parallel   bool
	input      string
	script     string
	force      bool
	maxRetries uint8
}

func newCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tm",
		Short: "The simplest task manager",
	}

	cmd.AddCommand(newDoCommand())

	return cmd
}

func newDoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "do",
		Short: "Do tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			params, err := parseFlags(cmd)
			if err != nil {
				return err
			}

			doer, err := getDoer(params)
			if err != nil {
				return err
			}

			ctx := context.Background()
			doer.Do(ctx, func(ctx context.Context, task string) error {
				cmd := exec.Command(params.script, task)
				cmd.Env = os.Environ()
				cmd.Stdout = os.Stdout
				return cmd.Run()
			})
			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "tasks.txt", "Tasks file")
	cmd.Flags().StringP("script", "s", "./do.sh", "Script file")
	cmd.Flags().BoolP("parallel", "p", true, "Enable parallel mode")
	cmd.Flags().BoolP("force", "f", false, "Do tasks even if they are already done")
	cmd.Flags().Uint8("max-retries", 10, "Max retries")

	return cmd
}

func getDoer(params *doParams) (tm.Doer, error) {
	taskManager, err := tm.NewTasksManager(tm.Config{
		Force:      params.force,
		FilePath:   params.input,
		MaxRetries: params.maxRetries,
	})
	if err != nil {
		return nil, err
	}
	if params.parallel {
		return tm.NewDefaultParallelDoer(taskManager), nil
	}
	return tm.NewDefaultOrderedDoer(taskManager), nil
}

func parseFlags(cmd *cobra.Command) (*doParams, error) {
	// check parallel
	parallel, err := cmd.Flags().GetBool("parallel")
	if err != nil {
		return nil, err
	}

	// check input file
	// todo check read permissions
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}
	inputExists, err := exists(input)
	if err != nil {
		return nil, err
	}
	if !inputExists {
		return nil, fmt.Errorf("input file does not exist: %s", input)
	}

	// check script file
	// todo check execute permissions
	script, err := cmd.Flags().GetString("script")
	if err != nil {
		return nil, err
	}
	script = strings.TrimSpace(script)
	if script == "" {
		return nil, fmt.Errorf("script is empty")
	}
	scriptExists, err := exists(script)
	if err != nil {
		return nil, err
	}
	if !scriptExists {
		return nil, fmt.Errorf("script file does not exist: %s", script)
	}

	// check force
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return nil, err
	}

	// check max retries
	maxRetries, err := cmd.Flags().GetUint8("max-retries")
	if err != nil {
		return nil, err
	}

	return &doParams{
		parallel:   parallel,
		input:      input,
		script:     script,
		force:      force,
		maxRetries: maxRetries,
	}, nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
