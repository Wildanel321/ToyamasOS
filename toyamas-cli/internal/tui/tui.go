package tui

import (
	"fmt"
	"strings"
	"time"
	"toyamas-cli/internal/printer"
)

type Spinner struct {
	stopChan chan bool
	msg      string
}

func StartSpinner(message string) *Spinner {
	if printer.CurrentMode == printer.ModeJSON {
		return &Spinner{}
	}

	stop := make(chan bool)
	s := &Spinner{stopChan: stop, msg: message}

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		for {
			select {
			case <-stop:
				fmt.Print("\r\033[K") // Clear line
				return
			default:
				fmt.Printf("\r\033[1;36m%s\033[0m %s...", frames[i%len(frames)], s.msg)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	return s
}

func (s *Spinner) Stop(successMsg string) {
	if printer.CurrentMode == printer.ModeJSON {
		return
	}
	if s.stopChan != nil {
		s.stopChan <- true
	}
	if successMsg != "" {
		fmt.Printf("\r\033[1;32m✔\033[0m %s\n", successMsg)
	}
}

func RenderProgressBar(current, total int, label string) {
	if printer.CurrentMode == printer.ModeJSON {
		return
	}

	width := 30
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	fmt.Printf("\r\033[1;34m%s\033[0m [%s] %d%% (%d/%d)", label, bar, int(percent*100), current, total)
	if current >= total {
		fmt.Println()
	}
}

func RenderTable(headers []string, rows [][]string) {
	if printer.CurrentMode == printer.ModeJSON {
		return
	}

	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, val := range row {
			if i < len(colWidths) && len(val) > colWidths[i] {
				colWidths[i] = len(val)
			}
		}
	}

	// Print Header
	fmt.Print("\033[1;37m┌")
	for i, w := range colWidths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐\033[0m")

	fmt.Print("\033[1;37m│")
	for i, h := range headers {
		fmt.Printf(" %-*s │", colWidths[i], h)
	}
	fmt.Println("\033[0m")

	fmt.Print("\033[1;37m├")
	for i, w := range colWidths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤\033[0m")

	// Print Rows
	for _, row := range rows {
		fmt.Print("│")
		for i, val := range row {
			if i < len(colWidths) {
				fmt.Printf(" %-*s │", colWidths[i], val)
			}
		}
		fmt.Println()
	}

	// Bottom Border
	fmt.Print("\033[1;37m└")
	for i, w := range colWidths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("┘\033[0m")
}

func ShowBanner() {
	if printer.CurrentMode == printer.ModeJSON {
		return
	}
	fmt.Println(`
\033[1;32m  _____ ______   _____  ___  ___  ___  _____ 
 |_   _/ _ \ \ / / / \   |_  _//  \/  \ /  |__/ / ___/ 
   | || (_) \ V / / _ \   | | / /\ / /\ V /|  __\___ \ 
   |_| \___/ |_|/_/   \_\ |_|/_/  /_/  \_/ |_|  /____/\033[0m
\033[1;34m  ToyamasOS Enterprise CLI v1.0.0 (Debian 13 Minimal)\033[0m
======================================================`)
}
