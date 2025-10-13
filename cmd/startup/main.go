package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

type StartupManager struct {
	projectRoot  string
	backendCmd   *exec.Cmd
	frontendCmd  *exec.Cmd
	dbRunning    bool
	envVars      map[string]string
}

func main() {
	manager := &StartupManager{
		envVars: make(map[string]string),
	}

	// Get project root directory
	projectRoot, err := findProjectRoot()
	if err != nil {
		printError("Failed to find project root: %v", err)
		os.Exit(1)
	}
	manager.projectRoot = projectRoot

	printHeader("Template Store - Automated Startup")
	fmt.Println()

	// Load environment variables
	printStep("Loading environment variables...")
	if err := manager.loadEnv(); err != nil {
		printWarning("Could not load .env file: %v", err)
		printInfo("Continuing with system environment variables...")
	} else {
		printSuccess("Environment variables loaded")
	}

	// Check prerequisites
	printStep("Checking prerequisites...")
	if err := manager.checkPrerequisites(); err != nil {
		printError("Prerequisites check failed: %v", err)
		os.Exit(1)
	}
	printSuccess("All prerequisites met")

	// Start PostgreSQL
	printStep("Starting PostgreSQL database...")
	if err := manager.startPostgres(); err != nil {
		printError("Failed to start PostgreSQL: %v", err)
		os.Exit(1)
	}
	manager.dbRunning = true
	printSuccess("PostgreSQL is running")

	// Wait for database to be ready
	printStep("Waiting for database to be ready...")
	if err := manager.waitForDatabase(); err != nil {
		printError("Database failed to become ready: %v", err)
		manager.cleanup()
		os.Exit(1)
	}
	printSuccess("Database is ready")

	// Build backend
	printStep("Building backend server...")
	if err := manager.buildBackend(); err != nil {
		printError("Failed to build backend: %v", err)
		manager.cleanup()
		os.Exit(1)
	}
	printSuccess("Backend built successfully")

	// Start backend
	printStep("Starting backend server...")
	if err := manager.startBackend(); err != nil {
		printError("Failed to start backend: %v", err)
		manager.cleanup()
		os.Exit(1)
	}
	printSuccess("Backend server started on port %s", manager.getEnv("PORT", "8080"))

	// Wait for backend to be ready
	printStep("Waiting for backend to be ready...")
	if err := manager.waitForBackend(); err != nil {
		printError("Backend failed to start: %v", err)
		manager.cleanup()
		os.Exit(1)
	}
	printSuccess("Backend is ready")

	// Seed database
	printStep("Seeding database with sample data...")
	if err := manager.seedDatabase(); err != nil {
		printWarning("Failed to seed database: %v", err)
		printInfo("You may need to seed manually using the API endpoints")
	} else {
		printSuccess("Database seeded successfully")
	}

	// Start frontend
	printStep("Starting frontend server...")
	if err := manager.startFrontend(); err != nil {
		printError("Failed to start frontend: %v", err)
		manager.cleanup()
		os.Exit(1)
	}
	printSuccess("Frontend server started on port 3000")

	// Open browser
	printStep("Opening browser...")
	if err := manager.openBrowser("http://localhost:3000"); err != nil {
		printWarning("Could not open browser automatically: %v", err)
		printInfo("Please open http://localhost:3000 manually")
	} else {
		printSuccess("Browser opened")
	}

	fmt.Println()
	printHeader("Application is running!")
	fmt.Println()
	printInfo("🌐 Frontend: http://localhost:3000")
	printInfo("🔧 Backend:  http://localhost:%s", manager.getEnv("PORT", "8080"))
	printInfo("🗄️  Database: localhost:5432")
	fmt.Println()
	printInfo("Press Ctrl+C to stop all services")
	fmt.Println()

	// Wait for interrupt signal
	manager.waitForShutdown()
}

func findProjectRoot() (string, error) {
	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for go.mod file to identify project root
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root (looking for go.mod)")
}

func (m *StartupManager) loadEnv() error {
	envPath := filepath.Join(m.projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	// Cache common env vars
	m.envVars["PORT"] = os.Getenv("PORT")
	m.envVars["DB_HOST"] = os.Getenv("DB_HOST")
	m.envVars["DB_PORT"] = os.Getenv("DB_PORT")
	m.envVars["DB_USER"] = os.Getenv("DB_USER")
	m.envVars["DB_PASSWORD"] = os.Getenv("DB_PASSWORD")
	m.envVars["DB_NAME"] = os.Getenv("DB_NAME")

	return nil
}

func (m *StartupManager) getEnv(key, defaultValue string) string {
	if val, ok := m.envVars[key]; ok && val != "" {
		return val
	}
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func (m *StartupManager) checkPrerequisites() error {
	// Check Docker
	if err := exec.Command("docker", "--version").Run(); err != nil {
		return fmt.Errorf("Docker is not installed or not in PATH")
	}

	// Check Docker Compose
	if err := exec.Command("docker-compose", "--version").Run(); err != nil {
		return fmt.Errorf("Docker Compose is not installed or not in PATH")
	}

	// Check if Docker daemon is running
	if err := exec.Command("docker", "info").Run(); err != nil {
		return fmt.Errorf("Docker daemon is not running. Please start Docker Desktop")
	}

	// Check Go
	if err := exec.Command("go", "version").Run(); err != nil {
		return fmt.Errorf("Go is not installed or not in PATH")
	}

	return nil
}

func (m *StartupManager) startPostgres() error {
	// Check if postgres container is already running
	cmd := exec.Command("docker", "ps", "--filter", "name=postgres", "--format", "{{.Names}}")
	cmd.Dir = m.projectRoot
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "postgres") {
		printInfo("PostgreSQL container already running")
		return nil
	}

	// Start PostgreSQL using docker-compose
	cmd = exec.Command("docker-compose", "up", "-d", "postgres")
	cmd.Dir = m.projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker-compose failed: %w", err)
	}

	return nil
}

func (m *StartupManager) waitForDatabase() error {
	dbHost := m.getEnv("DB_HOST", "localhost")
	dbPort := m.getEnv("DB_PORT", "5432")
	dbUser := m.getEnv("DB_USER", "postgres")
	dbPassword := m.getEnv("DB_PASSWORD", "postgres")
	dbName := m.getEnv("DB_NAME", "template_store")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		db, err := sql.Open("pgx", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				db.Close()
				return nil
			}
			db.Close()
		}

		if i < maxAttempts-1 {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
	}

	fmt.Println()
	return fmt.Errorf("database did not become ready within %d seconds", maxAttempts)
}

func (m *StartupManager) buildBackend() error {
	cmd := exec.Command("go", "build", "-o", "main", "cmd/server/main.go")
	cmd.Dir = m.projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}

func (m *StartupManager) startBackend() error {
	mainPath := filepath.Join(m.projectRoot, "main")

	cmd := exec.Command(mainPath)
	cmd.Dir = m.projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start backend: %w", err)
	}

	m.backendCmd = cmd
	return nil
}

func (m *StartupManager) waitForBackend() error {
	port := m.getEnv("PORT", "8080")
	url := fmt.Sprintf("http://localhost:%s/health", port)

	maxAttempts := 30
	client := &http.Client{Timeout: 2 * time.Second}

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		if i < maxAttempts-1 {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
	}

	fmt.Println()
	return fmt.Errorf("backend did not become ready within %d seconds", maxAttempts)
}

func (m *StartupManager) seedDatabase() error {
	port := m.getEnv("PORT", "8080")
	client := &http.Client{Timeout: 10 * time.Second}

	// Seed categories first (templates reference categories)
	categoriesURL := fmt.Sprintf("http://localhost:%s/api/v1/categories/seed", port)
	resp, err := client.Post(categoriesURL, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("categories seed returned status %d", resp.StatusCode)
	}

	// Seed templates
	templatesURL := fmt.Sprintf("http://localhost:%s/api/v1/templates/seed", port)
	resp, err = client.Post(templatesURL, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to seed templates: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("templates seed returned status %d", resp.StatusCode)
	}

	return nil
}

func (m *StartupManager) startFrontend() error {
	// Check if port 3000 is already in use
	if err := m.killProcessOnPort("3000"); err != nil {
		printWarning("Could not check/kill process on port 3000: %v", err)
	}

	webDir := filepath.Join(m.projectRoot, "web")

	// Start a simple Go HTTP server for the frontend
	go func() {
		fs := http.FileServer(http.Dir(webDir))
		http.Handle("/", fs)

		printInfo("Frontend server listening on :3000")
		if err := http.ListenAndServe(":3000", nil); err != nil {
			printError("Frontend server error: %v", err)
		}
	}()

	// Give it a moment to start
	time.Sleep(500 * time.Millisecond)

	// Test if frontend is accessible
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:3000")
	if err != nil {
		return fmt.Errorf("frontend server not responding: %w", err)
	}
	resp.Body.Close()

	return nil
}

func (m *StartupManager) killProcessOnPort(port string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		// Use lsof to find the process using the port
		cmd = exec.Command("lsof", "-ti", fmt.Sprintf(":%s", port))
		output, err := cmd.Output()
		if err != nil {
			// No process found on this port, which is fine
			return nil
		}

		pid := strings.TrimSpace(string(output))
		if pid != "" {
			printInfo("Found process %s using port %s, killing it...", pid, port)
			killCmd := exec.Command("kill", "-9", pid)
			if err := killCmd.Run(); err != nil {
				return fmt.Errorf("failed to kill process %s: %w", pid, err)
			}
			printSuccess("Killed process %s on port %s", pid, port)
			// Give the OS time to release the port
			time.Sleep(500 * time.Millisecond)
		}
	case "windows":
		// Use netstat to find the process
		cmd = exec.Command("netstat", "-ano")
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, fmt.Sprintf(":%s", port)) && strings.Contains(line, "LISTENING") {
				fields := strings.Fields(line)
				if len(fields) > 4 {
					pid := fields[len(fields)-1]
					printInfo("Found process %s using port %s, killing it...", pid, port)
					killCmd := exec.Command("taskkill", "/F", "/PID", pid)
					if err := killCmd.Run(); err != nil {
						return fmt.Errorf("failed to kill process %s: %w", pid, err)
					}
					printSuccess("Killed process %s on port %s", pid, port)
					time.Sleep(500 * time.Millisecond)
					break
				}
			}
		}
	}

	return nil
}

func (m *StartupManager) openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

func (m *StartupManager) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println()
	printHeader("Shutting down...")
	fmt.Println()

	m.cleanup()

	printSuccess("All services stopped")
	os.Exit(0)
}

func (m *StartupManager) cleanup() {
	// Stop backend
	if m.backendCmd != nil && m.backendCmd.Process != nil {
		printStep("Stopping backend server...")
		if err := m.backendCmd.Process.Signal(os.Interrupt); err != nil {
			// If interrupt fails, force kill
			m.backendCmd.Process.Kill()
		}
		// Wait for process to exit
		m.backendCmd.Wait()
		printSuccess("Backend stopped")
	}

	// Frontend will stop when the program exits (it's a goroutine)
	printStep("Stopping frontend server...")
	printSuccess("Frontend stopped")

	// Stop PostgreSQL
	if m.dbRunning {
		printStep("Stopping PostgreSQL...")
		cmd := exec.Command("docker-compose", "stop", "postgres")
		cmd.Dir = m.projectRoot
		cmd.Run() // Ignore errors during cleanup
		printSuccess("PostgreSQL stopped")
	}
}

// Printing utilities
func printHeader(message string) {
	fmt.Printf("%s%s═══════════════════════════════════════════════════%s\n", colorCyan, "", colorReset)
	fmt.Printf("%s  %s%s\n", colorCyan, message, colorReset)
	fmt.Printf("%s═══════════════════════════════════════════════════%s\n", colorCyan, colorReset)
}

func printStep(format string, args ...interface{}) {
	fmt.Printf("%s▶ %s%s\n", colorBlue, fmt.Sprintf(format, args...), colorReset)
}

func printSuccess(format string, args ...interface{}) {
	fmt.Printf("%s✓ %s%s\n", colorGreen, fmt.Sprintf(format, args...), colorReset)
}

func printError(format string, args ...interface{}) {
	fmt.Printf("%s✗ Error: %s%s\n", colorRed, fmt.Sprintf(format, args...), colorReset)
}

func printWarning(format string, args ...interface{}) {
	fmt.Printf("%s⚠ Warning: %s%s\n", colorYellow, fmt.Sprintf(format, args...), colorReset)
}

func printInfo(format string, args ...interface{}) {
	fmt.Printf("%sℹ %s%s\n", colorCyan, fmt.Sprintf(format, args...), colorReset)
}
