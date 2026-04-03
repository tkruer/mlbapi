package mlbapi

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func asMap(value any) JSON {
	if mapped, ok := value.(map[string]any); ok {
		return mapped
	}
	return JSON{}
}

func asSlice(value any) []any {
	if slice, ok := value.([]any); ok {
		return slice
	}
	return nil
}

func nestedMap(root JSON, path ...string) JSON {
	current := any(root)
	for _, key := range path {
		next := asMap(current)
		current = next[key]
	}
	return asMap(current)
}

func nestedSlice(root JSON, path ...string) []any {
	if len(path) == 0 {
		return nil
	}
	current := root
	for _, key := range path[:len(path)-1] {
		current = asMap(current[key])
	}
	return asSlice(current[path[len(path)-1]])
}

func nestedValue(root JSON, path ...string) any {
	current := any(root)
	for _, key := range path {
		current = asMap(current)[key]
	}
	return current
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case float64:
		if math.Mod(typed, 1) == 0 {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case float32:
		if math.Mod(float64(typed), 1) == 0 {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(float64(typed), 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case TeamID:
		return strconv.Itoa(int(typed))
	case bool:
		return strconv.FormatBool(typed)
	case time.Time:
		return typed.Format(time.RFC3339)
	default:
		return fmt.Sprint(value)
	}
}

func intValue(value any) int {
	switch typed := value.(type) {
	case nil:
		return 0
	case int:
		return typed
	case int8:
		return int(typed)
	case int16:
		return int(typed)
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case TeamID:
		return int(typed)
	case uint:
		return int(typed)
	case uint8:
		return int(typed)
	case uint16:
		return int(typed)
	case uint32:
		return int(typed)
	case uint64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	case string:
		if typed == "" {
			return 0
		}
		number, err := strconv.Atoi(typed)
		if err != nil {
			return 0
		}
		return number
	default:
		return 0
	}
}

func boolValue(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, err := strconv.ParseBool(typed)
		return err == nil && parsed
	default:
		return false
	}
}

func sortedKeys(value JSON) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func center(text string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(text) >= width {
		return text
	}
	padding := width - len(text)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func left(text string, width int) string {
	if len(text) >= width {
		return text
	}
	return text + strings.Repeat(" ", width-len(text))
}

func wrapText(text string, width int) []string {
	if width <= 0 || len(text) <= width {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, len(words))
	current := words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = "    " + word
	}
	lines = append(lines, current)
	return lines
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func yearPrefix(value string) string {
	if len(value) >= 4 {
		return value[:4]
	}
	return value
}

func currentYear() int {
	return time.Now().Year()
}

func seasonFromDate(date string) int {
	if parsed, err := time.Parse("2006-01-02", date); err == nil {
		return parsed.Year()
	}
	if len(date) >= 4 {
		if year, err := strconv.Atoi(date[:4]); err == nil {
			return year
		}
		if year, err := strconv.Atoi(date[len(date)-4:]); err == nil {
			return year
		}
	}
	return currentYear()
}
