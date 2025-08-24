// pkg/logger/pretty_logger.go
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
)

// ANSI Color Codes (Bright versions)
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[91m"
	ColorGreen  = "\033[92m"
	ColorYellow = "\033[93m"
	ColorBlue   = "\033[94m"
	ColorCyan   = "\033[96m"
	ColorPurple = "\033[95m"
)

type prettyLogger struct{}

func NewPrettyLogger() Logger {
	return &prettyLogger{}
}

// ✨ อัปเกรด getFileInfo ให้ฉลาดขึ้น
// เราจะหา Caller ที่อยู่นอก package 'logger' ของเรา
func getFileInfo() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// ถ้าเจอไฟล์ที่ไม่ได้อยู่ใน package logger ก็คือไฟล์ที่เรียกเราจริงๆ!
		if !strings.Contains(file, "pkg/logger") {
			parts := strings.Split(file, "/")
			return fmt.Sprintf("%s:%d", parts[len(parts)-1], line)
		}
	}
	return "???:0"
}

// ✨ อัปเกรดให้เข้าใจ key-value pairs แบบ slog
func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(" |") // ใช้ | คั่นระหว่าง message หลักกับ args

	for i := 0; i < len(args); i += 2 {
		key := args[i]
		var value any = "(MISSING)"
		if i+1 < len(args) {
			value = args[i+1]
		}
		builder.WriteString(fmt.Sprintf(" %s=%v", key, value))
	}
	return builder.String()
}

func (l *prettyLogger) Debug(msg string, args ...any) {
	location := getFileInfo()
	formattedArgs := formatArgs(args...)
	log.Printf("%s🐛 DEBUG %s: %s%s%s", ColorBlue, location, msg, formattedArgs, ColorReset)
}

func (l *prettyLogger) Info(msg string, args ...any) {
	location := getFileInfo()
	formattedArgs := formatArgs(args...)
	log.Printf("%sℹ️  INFO  %s: %s%s%s", ColorCyan, location, msg, formattedArgs, ColorReset)
}

func (l *prettyLogger) Success(msg string, args ...any) {
	location := getFileInfo()
	formattedArgs := formatArgs(args...)
	log.Printf("%s✅ SUCCESS %s: %s%s%s", ColorGreen, location, msg, formattedArgs, ColorReset)
}

func (l *prettyLogger) Warn(msg string, args ...any) {
	location := getFileInfo()
	formattedArgs := formatArgs(args...)
	log.Printf("%s⚠️  WARN  %s: %s%s%s", ColorYellow, location, msg, formattedArgs, ColorReset)
}

func (l *prettyLogger) Error(msg string, err error, args ...any) {
	location := getFileInfo()
	// สำหรับ Error เราจะเพิ่ม field 'err' เข้าไปใน args ด้วย
	allArgs := append(args, "err", err)
	formattedArgs := formatArgs(allArgs...)
	log.Printf("%s❌ ERROR %s: %s%s%s", ColorRed, location, msg, formattedArgs, ColorReset)
}

func (l *prettyLogger) Print(msg string) {
	location := getFileInfo()
	// เลือกสีตาม Level
	color := ColorPurple
	log.Printf("%s🔍 Print  %s: %s\n%s%s", color, location, msg, ColorReset)
}

func (l *prettyLogger) Dump(data interface{}) {
	location := getFileInfo()
	// แปลง object เป็น JSON สวยๆ
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		l.Error("Failed to dump data", err)
		return
	}
	color := ColorPurple

	log.Printf("%s🔍 DUMP  %s: %s\n%s%s", color, location, string(jsonBytes), ColorReset)
}

func (l *prettyLogger) Dumpf(level string, msg string, data interface{}) {
	location := getFileInfo()

	// แปลง object เป็น JSON สวยๆ
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		l.Error("Failed to dump data", err)
		return
	}
	// เลือกสีตาม Level
	color := ColorPurple
	switch strings.ToUpper(level) {
	case LevelDebug:
		color = ColorBlue
	case LevelInfo:
		color = ColorCyan
	case LevelWarn:
		color = ColorYellow
	case LevelError:
		color = ColorRed
	case LevelSuccess:
		color = ColorGreen
	}
	log.Printf("%s🔍 DUMP_F  %s: %s\n%s%s", color, location, msg, string(jsonBytes), ColorReset)
}
