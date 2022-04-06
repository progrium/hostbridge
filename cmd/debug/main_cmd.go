//go:build cmd

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"tractor.dev/hostbridge/bridge/misc"
	"tractor.dev/hostbridge/client"
)

func main() {
	os.Setenv("BRIDGECMD", "./hostbridge")
	c, err := client.Spawn()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.OnEvent = func(e client.Event) {
		log.Println(e)
	}

	ctx := context.Background()

	if err := c.App.Run(ctx, client.AppOptions{}); err != nil {
		panic(err)
	}

	options := client.WindowOptions{
		Title: "Demo window",
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		Transparent: true,
		Frameless:   false,
		Resizable:   true,
		Visible:     true,
		//Position: window.Position{X: 10, Y: 10},
		Size:   client.Size{Width: 360, Height: 240},
		Center: true,
		HTML: `
			<!doctype html>
			<html>
				<body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
				<script>
					window.onload = function() {
						document.body.innerHTML = '<div style="padding: 30px">Transparency Test!<br><br>${navigator.userAgent}</div>';
					};
				</script>
			</html>
		`,
	}

	w1, err := c.Window.New(ctx, options)
	if err != nil {
		panic(err)
	}

	fmt.Println("[main] window", w1)

	if w1 == nil {
		return
	}

	if err := w1.SetTitle(ctx, "Hello, Sailor!"); err != nil {
		panic(err)
	}

	trayTemplate := []client.MenuItem{
		{
			Title: "Click on this here thing",
		},
		{
			Title: "Secret stuff",
			SubMenu: []client.MenuItem{
				{
					ID:    1,
					Title: "I'm nested!!",
				},
				{
					ID:       101,
					Disabled: true,
					Title:    "Can't touch this",
				},
			},
		},
		{
			Title:       "Quit",
			Accelerator: "Command+T",
		},
	}

	iconData, err := misc.Assets.ReadFile("icon.png")
	if err != nil {
		fmt.Println("Error reading icon file:", err)
	}

	if err := c.App.NewIndicator(ctx, iconData, trayTemplate); err != nil {
		log.Fatal(err)
	}

	select {}
}
