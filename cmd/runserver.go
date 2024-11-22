/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/conf"
	"github.com/spf13/cobra"
)

var cmd *exec.Cmd

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start a lightweight Web server for development.",
	Run: func(cmd *cobra.Command, args []string) {
		noreload, _ := cmd.Flags().GetBool("noreload")
			if noreload {
				startGorimServer() // Start the Echo server directly
			} else {
				runServerWithHotReload() // Start the watcher for hot reload
			}
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
	runserverCmd.Flags().Bool("noreload", false, "Run the server without hot reload")
}

func runServerWithHotReload() {
    // Create a new file watcher
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // Watch the current directory and its subdirectories
    if err := watchRecursive(".", watcher); err != nil {
        log.Fatal(err)
    }

    // Start the server for the first time
    startServer()

    // Signal channel to listen for termination (CTRL+C) and exit cleanly
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    // Channel to handle file events
    done := make(chan bool)

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                // Check if the file is modified, renamed, or created
                if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Create == fsnotify.Create {
                    fmt.Println("Detected file change:", event.Name)

                    // If a new directory is created, start watching it recursively
                    fileInfo, err := os.Stat(event.Name)
                    if err == nil && fileInfo.IsDir() {
                        if err := watchRecursive(event.Name, watcher); err != nil {
                            log.Printf("Error watching new directory: %s\n", err)
                        }
                    }

                    restartServer() // Restart the server on file change
                }

            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Error:", err)

            case <-signalChan:
                fmt.Println("Received interrupt, shutting down...")
                stopServer()
                done <- true
            }
        }
    }()

    <-done
}

// Start the server by running `go run main.go`
func startServer() {
    cmd = exec.Command("go", "run", "main.go", "runserver", "--noreload")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Start()
    if err != nil {
        fmt.Printf("Error starting server: %s\n", err)
        return
    }

    // fmt.Printf("Server running with PID %d\n", cmd.Process.Pid)
}

// Stop the running server by killing the process
func stopServer() {
    if cmd != nil && cmd.Process != nil {
        fmt.Println("Stopping server...")
		killProcessAndChildren(cmd.Process.Pid)

        cmd = nil
        // Ensure the port is released
        time.Sleep(1 * time.Second)
    }
}

// Restart the server by stopping and then starting it again
func restartServer() {
    stopServer()
    // Check if the port is still in use
    startServer()
}

// Recursively watch directories and subdirectories
func watchRecursive(dir string, watcher *fsnotify.Watcher) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            err = watcher.Add(path)
            if err != nil {
                return fmt.Errorf("failed to watch directory %s: %v", path, err)
            }
            // fmt.Printf("Watching directory: %s\n", path)
        }
        return nil
    })
}

// killProcessAndChildren forcefully kills the process and all of its children
func killProcessAndChildren(pid int) {
    // Kill the process tree starting from the PID
    // fmt.Printf("Killing process and children with PID %d\n", pid)

    process := exec.Command("pkill", "-TERM", "-P", fmt.Sprint(pid))
    process.Stdout = os.Stdout
    process.Stderr = os.Stderr
    if err := process.Run(); err != nil {
        fmt.Printf("Error killing process tree: %s\n", err)
    }
	cmd.Wait()
}

func printBanner(version string, env string, host string) {
	// Current timestamp
	now := time.Now().Format("January 02, 2006 - 15:04:05")

	// Banner content
	fmt.Println("System check identified no issues (0 silenced).")
	fmt.Println(now)
	fmt.Printf("Gorim version %s using env '%s'\n", version, env)
	fmt.Printf("Starting development server at %s\n", host)
	fmt.Println("Quit the server with CONTROL-C.")
}

// startGorimServer starts the Echo server directly without hot reload.
func startGorimServer() {
	server := conf.GorimServer.(*gorim.Server)
    address := fmt.Sprintf("%s:%d", conf.HOST, conf.PORT)
    versionNumber := "v1.1.0"
    printBanner(versionNumber, conf.ENV_PATH, conf.HOST)

    go func() {
        if err := server.Start(address); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    fmt.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Failed to gracefully shut down server: %v", err)
    }
    fmt.Println("Server stopped.")
}
