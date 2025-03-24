package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Read directory path from flag, default = current directory
	dir := flag.String("path", ".", "Directory to watch")

	flag.Usage = func() {
		fmt.Println("Go File Watcher — A CLI tool to monitor file changes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  file-watcher [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  -path <dir>   Directory to watch (default: current dir)")
		fmt.Println("  -help         Show this help message")
	}
	flag.Parse()

	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Add directory
	if err := watcher.Add(*dir); err != nil {
		log.Fatalf("Failed to watch directory %s: %v", *dir, err)
	}
	fmt.Printf("Watching directory: %s\n", *dir)

	// Setup logging: both stdout and file
	logFile, err := os.OpenFile("tasks.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)

	// Handle events & errors
	for {
		select {
		case event := <-watcher.Events:
			logger.Println("Event:", event)
		case err := <-watcher.Errors:
			logger.Println("Error:", err)
		}
	}
}
