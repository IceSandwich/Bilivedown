package Logging

import (
	"errors"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	nerr := errors.New("some error")
	fmt.Printf("\nLogToConsole-*-\n")
	err := Init(false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("LogError:")
	Error("some descr", nerr)
	fmt.Println("LogRePrint:")
	RePrint("first descr")
	RePrint("second desc")
	fmt.Printf("\nLogToFile-*-\n")
	err = Init(true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("LogError:")
	Error("file descr", nerr)
	fmt.Println("LogRePrint:")
	RePrint("file first descr")
	RePrint("file second desc")
	err = Close()
	if err != nil {
		t.Error(err)
	}
}
