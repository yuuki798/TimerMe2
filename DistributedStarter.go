// run_services.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	dirs := []string{"./gateway", "./task_service"}

	for _, dir := range dirs {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			if err := runCommand(dir); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}(dir)
	}

	wg.Wait()
}
func runCommand(dir string) error {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command in %s: %w", dir, err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command in %s finished with error: %w", dir, err)
	}

	return nil
}
