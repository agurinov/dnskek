package log

import (
	"fmt"
	golog "log"
	"os"
	"strconv"
	"strings"
)

const modChar = "\x1b"

type Modification int
type Modifications []Modification

func (ms Modifications) String() string {
	toString := make([]string, len(ms))

	for i, v := range ms {
		toString[i] = strconv.Itoa(int(v))
	}

	return fmt.Sprintf("%s[%sm", modChar, strings.Join(toString, ";"))
}

const (
	// Common attributes [0, 1, 2, 4, 5, 7]
	Default Modification = iota
	Bold
	SemiBright
	_
	UnderLine
	Blink
	_
	Reversion
)

const (
	// Char colors [30-37]
	BlackChar Modification = iota + 30
	RedChar
	GreenChar
	YellowChar
	BlueChar
	PurpleChar
	AquamarineChar
	GrayChar
)

const (
	// Background colors [40-47]
	BlackBackground Modification = iota + 40
	RedBackground
	GreenBackground
	YellowBackground
	BlueBackground
	PurpleBackground
	AquamarineBackground
	GrayBackground
)

// colors and modes
// https://habrahabr.ru/post/119436/
// https://misc.flogisoft.com/bash/tip_colors_and_formatting
// TODO play with this https://stackoverflow.com/questions/28432398/difference-between-some-operators-golang
func Wrap(s string, mods ...Modification) string {
	// no modifications
	if mods == nil {
		return s
	}
	// calculate wrappers (opening and closing)
	openingWrapper := Modifications(mods).String()
	closingWrapper := Modifications([]Modification{Default}).String()
	// concat
	return openingWrapper + s + closingWrapper
}

// just shortcuts
var (
	debugPrefix   = Wrap("[DEBUG]", Bold, GrayChar) + "\t"
	errorPrefix   = Wrap("[ERROR]", Bold, RedChar, Blink) + "\t"
	infoPrefix    = Wrap("[INFO]", Bold, GreenChar) + "\t"
	warningPrefix = Wrap("[WARN]", Bold, YellowChar, Blink) + "\t"
)

type Logger struct {
	out   *golog.Logger
	err   *golog.Logger
	debug bool
}

func NewLogger() *Logger {
	return &Logger{
		out:   golog.New(os.Stdout, "", golog.LstdFlags),
		err:   golog.New(os.Stderr, "", golog.LstdFlags),
		debug: false,
	}
}

func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}

func (l *Logger) Debug(args ...interface{}) {
	if l.debug {
		l.err.SetPrefix(debugPrefix)
		l.err.Print(args...)
	}
}

func (l *Logger) Debugf(fmtString string, args ...interface{}) {
	if l.debug {
		l.err.SetPrefix(debugPrefix)
		l.err.Printf(fmtString+"\n", args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	l.err.SetPrefix(errorPrefix)
	l.err.Print(args...)
}

func (l *Logger) Errorf(fmtString string, args ...interface{}) {
	l.err.SetPrefix(errorPrefix)
	l.err.Printf(fmtString+"\n", args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.out.SetPrefix(infoPrefix)
	l.out.Print(args...)
}

func (l *Logger) Infof(fmtString string, args ...interface{}) {
	l.out.SetPrefix(infoPrefix)
	l.out.Printf(fmtString+"\n", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.err.SetPrefix(warningPrefix)
	l.err.Print(args...)
}

func (l *Logger) Warnf(fmtString string, args ...interface{}) {
	l.err.SetPrefix(warningPrefix)
	l.err.Printf(fmtString+"\n", args...)
}
