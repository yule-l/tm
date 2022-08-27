package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/yule-l/tm/pkg/tm"
)

func NewDoCommand() *cobra.Command {
	in := &tm.ParamsIn{}
	var scriptFilePath string
	cmd := &cobra.Command{
		Use:   "do",
		Short: "Do tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			doer, err := (&tm.Factory{}).NewDoer(in)
			if err != nil {
				return err
			}

			ctx := context.Background()
			doer.Do(ctx, func(ctx context.Context, task string) error {
				cmd := exec.Command(scriptFilePath, task)
				cmd.Env = os.Environ()
				cmd.Stdout = os.Stdout
				return cmd.Run()
			})
			return nil
		},
	}

	cmd.Flags().StringVarP(&in.TasksFilePath, "input", "i", "tasks.txt", "Tasks file")
	cmd.Flags().StringVarP(&scriptFilePath, "script", "s", "./do.sh", "Script file")
	cmd.Flags().BoolVarP(&in.Parallel, "parallel", "p", true, "Enable parallel mode")
	cmd.Flags().BoolVarP(&in.Force, "force", "f", false, "Do tasks even if they are already done")
	cmd.Flags().Uint8Var(&in.MaxRetries, "max-retries", 10, "Max retries")

	return cmd
}
