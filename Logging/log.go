package Logging

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var flog *os.File = nil
var tFormat = "2006-01-02 15:04:05"
var version = "Bilivedown 0.1 alpha"

func PrintVersion() {
	fmt.Printf("** %s **\n", version)
}

func GetTime() string {
	return time.Now().Format(tFormat)
}

func Init(logtoFile bool) (err error) {
	if logtoFile {
		if _, err = os.Stat("Log"); err != nil {
			if err = os.Mkdir("Log", os.ModePerm); err != nil {
				return
			}
		}
		if flog, err = os.Create(path.Join("Log", GetTime()+".log")); err != nil {
			return
		}
		fmt.Fprintf(flog, "%s I: Init module. %s on %s %s.\n", GetTime(), version, runtime.GOARCH, runtime.GOOS)
	}
	return nil
}

func Close() error {
	if flog != nil {
		if err := flog.Close(); err != nil {
			return err
		}
		flog = nil
	}
	return nil
}

func Error(description string, err error) {
	fmt.Fprintf(os.Stderr, "\n\033[31mError: %s\n - %s\n\033[0m", description, err)
	if flog != nil {
		fmt.Fprintf(flog, "%s E: %s\n%20sE: - %s\n", GetTime(), description, " ", err)
	}
}

func ErrorC(description string, err string) {
	fmt.Fprintf(os.Stderr, "\n\033[31mError: %s\n - %s\n\033[0m", description, err)
	if flog != nil {
		fmt.Fprintf(flog, "%s E: %s\n%20sE: - %s\n", GetTime(), description, " ", err)
	}
}

func ErrorS(description string) {
	fmt.Fprintf(os.Stderr, "\n\033[31mError: %s\033[0m\n", description)
	if flog != nil {
		fmt.Fprintf(flog, "%s E: %s\n", GetTime(), description)
	}
}

func Info(description string) {
	fmt.Fprintln(os.Stdout, description)
	if flog != nil {
		fmt.Fprintf(flog, "%s I: %s\n", GetTime(), description)
	}
}

func Warning(description string) {
	fmt.Fprintf(os.Stdout, "\033[33m*-* Warning: %s\033[0m\n", description)
	if flog != nil {
		fmt.Fprintf(flog, "%s W: %s\n", GetTime(), description)
	}
}

func Note(description string) {
	if flog == nil {
		Warning("Use note function without init logging moudle. Ignore note things.")
		return
	}
	fmt.Fprintf(flog, "%s N: %s\n", GetTime(), description)
}

func RePrint(description string) {
	fmt.Fprintf(os.Stdout, "\r\033[32m%s\033[0m", description)
	if flog != nil {
		fmt.Fprintf(flog, "%s I: %s\n", GetTime(), description)
	}
}

func Dump(varname string, data string) {
	if flog == nil {
		Warning("Use dump function without init logging moudle. Ignore dump things.")
		return
	}
	if strings.IndexByte(data, '\n') == -1 {
		fmt.Fprintf(flog, "%s D: %s = %s\n", GetTime(), varname, data)
	} else {
		fmt.Fprintf(flog, "%s D: %s(list) -->\n%s\n<-- dump list end\n", GetTime(), varname, data)
	}
}

func ReadArg(varname string, data string) {
	if flog == nil {
		Warning("Use readarg function without init logging moudle. Ignore arg things.")
		return
	}
	fmt.Fprintf(flog, "%s A: %s = %s\n", GetTime(), varname, data)
}
