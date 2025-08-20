package console

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
	cyan   = "\033[36m"

	checkMark = "✔"
	crossMark = "✖"
	infoMark  = "➤"
)

// init runs when the package is imported
func init() {
	if runtime.GOOS == "windows" {
		enableVirtualTerminal()
	}
}

// enableVirtualTerminal enables ANSI color support on Windows 10+
func enableVirtualTerminal() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	handle := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))

	// ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	setConsoleMode.Call(uintptr(handle), uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}

func Info(message string) {
	fmt.Println(cyan, infoMark, message, reset)
}

func Success(message string) {
	fmt.Println(green, checkMark, message, reset)
}

func Error(message string) {
	fmt.Println(red, crossMark, message, reset)
}

func SetColor(color string) {
	switch color {
	case "green":
		fmt.Print(green)
	case "red":
		fmt.Print(red)
	case "yellow":
		fmt.Print(yellow)
	case "cyan":
		fmt.Print(cyan)
	default:
		fmt.Print(reset)
	}
}

func ResetColor() {
	SetColor("default")
}

// Table prints any struct in an adjustable tabular format
func Table(data interface{}) {
	v := reflect.ValueOf(data)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		fmt.Println("Table() only accepts structs")
		return
	}

	// Determine max widths
	maxFieldLen := len("Field")
	maxValLen := len("Value")

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i).Name
		value := v.Field(i).Interface()

		valStr := fmt.Sprintf("%v", value)
		if v.Field(i).Kind() == reflect.Slice {
			valStr = strings.Join(value.([]string), ", ")
		}

		if len(field) > maxFieldLen {
			maxFieldLen = len(field)
		}
		if len(valStr) > maxValLen {
			maxValLen = len(valStr)
		}
	}

	// Draw table
	fmt.Printf("┌%s┬%s┐\n", strings.Repeat("─", maxFieldLen+2), strings.Repeat("─", maxValLen+2))
	fmt.Printf("│ %-*s │ %-*s │\n", maxFieldLen, "Field", maxValLen, "Value")
	fmt.Printf("├%s┼%s┤\n", strings.Repeat("─", maxFieldLen+2), strings.Repeat("─", maxValLen+2))

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i).Name
		value := v.Field(i).Interface()

		valStr := fmt.Sprintf("%v", value)
		if v.Field(i).Kind() == reflect.Slice {
			valStr = strings.Trim(fmt.Sprint(value), "[]")
		}

		fmt.Printf("│ %-*s │ %-*s │\n", maxFieldLen, field, maxValLen, valStr)
	}

	fmt.Printf("└%s┴%s┘\n", strings.Repeat("─", maxFieldLen+2), strings.Repeat("─", maxValLen+2))
}

// Log prints v to stdout as pretty JSON (print_r-like).
// Falls back to a %#v dump if JSON marshaling fails.
func Log(v any) {
	// Try JSON pretty print
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		return
	}

	// Fallback: Go-syntax detailed dump
	fmt.Printf("%#v\n", v)
}
