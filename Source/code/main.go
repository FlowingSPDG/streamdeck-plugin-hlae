package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FlowingSPDG/streamdeck"

	mirvpgl "github.com/FlowingSPDG/HLAE-Server-GO"
)

var hlaeserver *mirvpgl.HLAEServer

const (
	// AppName Streamdeck plugin app name
	AppName = "dev.flowingspdg.hlae.sdPlugin"

	// Action Name
	Action = "dev.flowingspdg.hlae.command"
)

func sendCommand(cmd string) error {
	if hlaeserver == nil {
		return fmt.Errorf("hlaeserver is not initialized")
	}
	if !launched || !connected {
		return fmt.Errorf("hlaeserver is not connected")
	}
	return hlaeserver.BroadcastRCON(cmd)
}

func main() {
	// Setup log file
	logfile, err := os.OpenFile("./streamdeck-hlae-plugin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	go func() {
		var err error
		hlaeserver, err = mirvpgl.New(":65535", "/std") // mirv_pgl url "ws://localhost:65535/std"
		if err != nil {
			panic(err)
		}

		if err := hlaeserver.Start(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	log.Println("Starting...")
	if err := run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	setup(client)

	return client.Run()
}

func setup(client *streamdeck.Client) {
	contexts := make(map[string]struct{})

	client.RegisterNoActionHandler(streamdeck.ApplicationDidLaunch, ApplicationDidLaunchHandler)
	client.RegisterNoActionHandler(streamdeck.ApplicationDidTerminate, ApplicationDidTerminateHandler)

	action := client.Action(Action)

	action.RegisterHandler(streamdeck.WillAppear, WillAppearHandler)
	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		contexts[event.Context] = struct{}{}
		return nil
	})

	action.RegisterHandler(streamdeck.WillDisappear, WillDisappearHandler)
	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		delete(contexts, event.Context)
		return nil
	})
	action.RegisterHandler(streamdeck.KeyDown, KeyDownHandler)

	action.RegisterHandler(streamdeck.DidReceiveSettings, DidReceiveSettingsHandler)
	action.RegisterHandler(streamdeck.SendToPlugin, SendToPluginHandler)
}
