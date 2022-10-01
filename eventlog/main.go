package main

import (
	"fmt"

	"golang.org/x/sys/windows/svc/eventlog"
)

func main() {

	elog, err := eventlog.Open("Winlogon")
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	fmt.Println(elog.Handle)
}
