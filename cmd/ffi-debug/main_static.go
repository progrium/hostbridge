package main

import "fmt"

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is necessary to fix the linker errors about `dlsym` that would otherwise appear.

/*
#cgo LDFLAGS: ./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit
#include "../../lib/hostbridge.h"
*/
import "C"

// 
// Lib
//

var user_main_loop func(event_type int)

//export go_main_loop
func go_main_loop(i C.int) {
	event_type := int(i)

	if (user_main_loop != nil) {
		user_main_loop(event_type)
	}
}

func Run(user_callback func(event_type int)) {
    user_main_loop = user_callback
    C.run(C.closure(C.go_main_loop))
}

type Window struct {
	Id    int
	Title string
}

func WindowCreate() Window {
	result := Window{}
	return result
}

func (it *Window) SetTitle(Title string) {
	it.Title = Title
}

//
// User Code
//

func main_loop(event_type int) {
	if (event_type > 0) {
		fmt.Println("%d", event_type);
	}
}

func main() {
	fmt.Println("[go] main");

	w := Window{}
	w.SetTitle("Hello, Sailor!")

	C.window_create(C.int(1280), C.int(720), C.CString("Hey"));

	fmt.Printf("%s\n", w.Title)

	fmt.Println("[go] run");
	Run(main_loop)

	fmt.Println("[go] this will never fire");
}
