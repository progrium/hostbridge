package bridge

import (
	"context"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/window"
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/talk"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	go func() {
		m.Run()
		app.Quit()
	}()
	app.Run(nil)
}

func TestBridge(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv := NewServer()
	go srv.Serve(l)

	client, err := talk.Dial("tcp", l.Addr().String(), codec.JSONCodec{})
	if err != nil {
		t.Fatal(err)
	}

	var opts interface{}
	opts = window.Options{
		HTML: `
			<!doctype html>
			<html>
				<body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
				<script>
					window.onload = function() {
						document.body.innerHTML = '<div style="padding: 30px">TEST</div>';
					};
				</script>
			</html>
		`,
	}
	_, err = client.Call(context.Background(), "window.Create", fn.Args{opts}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// uncomment to see visually
	time.Sleep(1 * time.Second)
}
