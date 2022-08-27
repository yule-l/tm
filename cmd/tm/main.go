package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yule-l/tm/internal/cmd"
)

func main() {
	ctx := context.Background()
	exitOnError(cmd.NewCli().ExecuteContext(ctx))
}

func exitOnError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Выполнение команды завершилось с ошибкой: %q", err)
	os.Exit(1)
}
