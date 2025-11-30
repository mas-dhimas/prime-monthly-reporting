package lib

import (
	"fmt"
)

func FormatSeconds(seconds int) string {
	if seconds <= 0 {
		return "0 seconds"
	}

	days := seconds / 86400
	seconds %= 86400

	hours := seconds / 3600
	seconds %= 3600

	minutes := seconds / 60
	seconds %= 60

	parts := []string{}

	if days > 0 {
		if days == 1 {
			parts = append(parts, fmt.Sprintf("%d day", days))
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, fmt.Sprintf("%d hour", hours))
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, fmt.Sprintf("%d minute", minutes))
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	// Only add seconds if duration is less than 1 hour
	if days == 0 && hours == 0 && seconds > 0 {
		if seconds == 1 {
			parts = append(parts, fmt.Sprintf("%d second", seconds))
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	return joinParts(parts)
}

func joinParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, p := range parts[1:] {
		result += " " + p
	}
	return result
}
