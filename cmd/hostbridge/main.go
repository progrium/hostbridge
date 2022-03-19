package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/progrium/qtalk-go/mux"
	"tractor.dev/hostbridge/bridge"
	"tractor.dev/hostbridge/bridge/core"
	"tractor.dev/hostbridge/cmd/hostbridge/demo"
)

const Version = "0.1.0"

func init() {
	runtime.LockOSThread()
}

func main() {
	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if flag.Arg(0) == "demo-build" {
		demo.Build()
		return
	}

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "hostbridge %s\n", Version)
	}

	sess, err := mux.DialIO(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	srv := bridge.NewServer()
	go srv.Respond(sess, context.Background())
	go func() {
		sess.Wait()
		core.Quit()
	}()
	core.Run(nil)
}
