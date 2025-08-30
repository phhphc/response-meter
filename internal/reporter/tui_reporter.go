package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/phhphc/response-meter/internal/meter"
)

type TUIReporter struct{}

func NewTUIReporter() meter.Reporter {
	return &TUIReporter{}
}

func (t TUIReporter) Update(s meter.Stats) error {
	t.resetScreen()
	t.displayDashboard(s)
	return nil
}

func (t TUIReporter) resetScreen() {
	fmt.Print("\033[2J\033[H")
}

func (t TUIReporter) displayDashboard(s meter.Stats) {
	width := 80

	title := "Response Meter"
	t.printCentered(title, width)

	// Calculate all-time metrics
	totalCount := 0
	for _, v := range s.TotalCounts {
		totalCount += v
	}

	var totalRequestRate float64
	if s.TotalDuration.Seconds() > 0 {
		totalRequestRate = float64(totalCount) / s.TotalDuration.Seconds()
	}

	// Calculate last period metrics
	lastCount := 0
	for _, v := range s.LastPeriodCounts {
		lastCount += v
	}

	var lastRequestRate float64
	if s.LastPeriodDuration.Seconds() > 0 {
		lastRequestRate = float64(lastCount) / s.LastPeriodDuration.Seconds()
	}

	// Display all-time metrics
	fmt.Println()
	t.printSectionHeader("ALL TIME", width)
	fmt.Println()
	t.displayMetrics(s.TotalDuration, totalRequestRate, totalCount, s.TotalCounts)

	// Display last period metrics
	fmt.Println()
	t.printSectionHeader("LAST PERIOD", width)
	fmt.Println()
	t.displayMetrics(s.LastPeriodDuration, lastRequestRate, lastCount, s.LastPeriodCounts)

	// Footer
	fmt.Println()
	fmt.Print("Press Ctrl+C to exit.")
}

func (t TUIReporter) displayMetrics(duration time.Duration, requestRate float64, totalCount int, counts map[string]int) {
	fmt.Printf("Duration: %-20s Requests/sec: %-10.1f Total Requests: %s\n", duration.Round(time.Second), requestRate, t.formatNumber(totalCount))
	fmt.Println("Response Distribution:")
	if totalCount > 0 {
		for response, count := range counts {
			displayResponse := t.formatResponse(response, 30)
			percentage := float64(count) / float64(totalCount) * 100
			progressBar := t.createProgressBar(percentage, 30)
			fmt.Printf("  %-30s %6.1f%% %s %s\n",
				displayResponse,
				percentage,
				progressBar,
				t.formatNumber(count))
		}
	} else {
		fmt.Println("  No data yet...")
	}
}

func (t TUIReporter) printCentered(text string, width int) {
	padding := max((width-len(text))/2, 0)
	fmt.Printf("%s%s\n", strings.Repeat(" ", padding), text)
}

func (t TUIReporter) printSectionHeader(title string, width int) {
	padding := max((width-len(title)-2)/2, 0)
	line := strings.Repeat("-", padding) + " " + title + " " + strings.Repeat("-", padding)
	for len(line) < width {
		line += "-"
	}
	fmt.Println(line)
}

func (t TUIReporter) formatResponse(response string, maxLen int) string {
	if len(response) <= maxLen {
		return response
	}
	return response[:maxLen-3] + "..."
}

func (t TUIReporter) createProgressBar(percentage float64, width int) string {
	filled := int(min(percentage, 100.0) / 100.0 * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

func (t TUIReporter) formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	result := ""
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}
