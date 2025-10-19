package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const serverAddr = "localhost:41114"

var (
	program string
	path    string
	circuit string
)

func startDigital(path string, circuit string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	if strings.HasSuffix(path, ".jar") {
		cmd = exec.Command("java", "-jar", path, "-circuit", circuit)
	} else {
		cmd = exec.Command(path, "-circuit", circuit)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Digital: %w", err)
	}

	log.Printf("Started Digital (PID: %d)", cmd.Process.Pid)
	return cmd, nil
}

func main() {
	flag.StringVar(&program, "program", "", "path to .hex file to be run")
	flag.StringVar(&program, "p", "", "path to .hex file to be run")
	flag.StringVar(&path, "digital", "", "path to Digital executable (.jar or .exe)")
	flag.StringVar(&path, "d", "", "path to Digital executable (.jar or .exe)")
	flag.StringVar(&circuit, "circuit", "", "path to Digital circuit (.dig)")
	flag.StringVar(&circuit, "c", "", "path to Digital circuit (.dig)")
	flag.Parse()

	if program == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current directory:", err)
		}
		program = filepath.Join(cwd, "example", "Conway32.hex")
	}

	if path == "" {
		path = os.Getenv("DIGITAL_PATH")
	}

	if path == "" {
		log.Fatal("Digital path not specified. Please either:\n" +
			"  1. Set the DIGITAL_PATH environment variable, or\n" +
			"  2. Use the -d flag: -d /path/to/Digital.jar")
	}

	if circuit == "" {
		circuit = os.Getenv("CIRCUIT_PATH")
	}

	if circuit == "" {
		log.Fatal("Circuit path not specified. Please either:\n" +
			"  1. Set the CIRCUIT_PATH environment variable, or\n" +
			"  2. Use the -c flag: -c /path/to/Circuit.dig")
	}

	cmd, err := startDigital(path, circuit)
	if err != nil {
		log.Fatal(err)
	}

	// Clean up process on interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down Digital...")
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
		os.Exit(0)
	}()

	log.Println("Waiting for Digital server to start...")
	time.Sleep(3 * time.Second)

	tcpClient := New(serverAddr, 5*time.Second)

	command := []string{"start:", program}
	start := strings.Join(command, "")

	response, err := tcpClient.SendCommand(start)
	if err != nil {
		log.Fatal("Failed to send start command:", err)
	}
	fmt.Println("Start response:", string(response))

	time.Sleep(1 * time.Second)

	a := app.New()
	w := a.NewWindow("Conway's Game of Life")
	gridSize := 32
	cellSize := float32(20)

	cells := InitGUI(w, gridSize, cellSize)

	go func() {
		// read memory addresses: 0x2000-0x23FF
		readout := "output:8192:9215"
		alive := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		dead := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		for {
			response, err := tcpClient.SendCommand(readout)
			if err != nil {
				log.Printf("Error: %v", err)
				time.Sleep(10 * time.Millisecond)
				continue
			}

			gridData, err := Grid(response)
			if err != nil {
				log.Printf("Error: %v", err)
				time.Sleep(10 * time.Millisecond)
				continue
			}

			fyne.DoAndWait(func() {
				UpdateGUI(cells, gridData, alive, dead)
			})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	w.ShowAndRun()

	// Clean up process when window closes
	if cmd != nil && cmd.Process != nil {
		log.Println("Shutting down Digital...")
		cmd.Process.Kill()
	}
}
