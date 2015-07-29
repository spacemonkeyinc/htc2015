// Copyright (C) 2015 Space Monkey, Inc.

package general

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/spacemonkeygo/flagfile"
	"github.com/spacemonkeygo/spacelog/setup"
	"gopkg.in/spacemonkeygo/monitor.v1"
)

var (
	debugAddr = flag.String("debug_addr", "localhost:0",
		"address to listen on for debugging")
)

func Run(main func() error) {
	flagfile.Load()
	setup.MustSetup(os.Args[0])
	monitor.RegisterEnvironment()
	go http.ListenAndServe(*debugAddr, monitor.DefaultStore)

	err := main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
