package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/FlowingSPDG/streamdeck"
)

var (
	connected = true // TODO...
	launched  = false
)

func init() {
}

// WillAppearHandler willAppear handler.
func WillAppearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	p := streamdeck.WillAppearPayload{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		return err
	}
	log.Println("WillAppearHandler:", p)

	s := PropertyInspector{}
	if err := json.Unmarshal(p.Settings, &s); err != nil {
		return err
	}

	settings.Save(event.Context, &s)

	log.Printf("settings for context%s context:%#v\n", event.Context, s)
	return nil
}

// WillDisappearHandler willDisappear handler
func WillDisappearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	log.Println("WillDisappearHandler")
	settings.Save(event.Context, &PropertyInspector{})
	log.Println("Refreshing settings for this context:", event.Context)
	s, _ := settings.Load(event.Context)
	return client.SetSettings(ctx, s)
}

// KeyDownHandler keyDown handler
func KeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	log.Println("KeyDownHandler")
	s, err := settings.Load(event.Context)
	if err != nil {
		return fmt.Errorf("couldn't find settings for context %v", event.Context)
	}
	log.Println("settings for this context:", s)

	if !connected || !launched {
		return client.ShowAlert(ctx)
	}

	if err := sendCommand(s.Command); err != nil {
		log.Println("Failed to execute command:", err)
		client.ShowAlert(ctx)
		return err
	}

	return client.ShowOk(ctx)
}

// ApplicationDidLaunchHandler applicationDidLaunch handler
func ApplicationDidLaunchHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	p := streamdeck.ApplicationDidLaunchPayload{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		log.Println("ERR:", err)
		return err
	}
	log.Println("ApplicationDidLaunchHandler:", p)
	if p.Application == "csgo.exe" {
		launched = true
	}
	return nil
}

// ApplicationDidTerminateHandler applicationDidTerminate handler
func ApplicationDidTerminateHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	p := streamdeck.ApplicationDidTerminatePayload{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		log.Println("ERR:", err)
		return err
	}
	log.Println("ApplicationDidTerminateHandler:", p)
	if p.Application == "csgo.exe" {
		launched = false
	}

	return nil
}

// DidReceiveSettingsHandler didReceiveSettings Handler
func DidReceiveSettingsHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	p := streamdeck.DidReceiveSettingsPayload{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		log.Println("ERR:", err)
		return err
	}
	log.Println("DidReceiveSettingsHandler:", p)

	s := &PropertyInspector{}
	if err := json.Unmarshal(p.Settings, s); err != nil {
		log.Println("ERR:", err)
		return err
	}
	settings.Save(event.Context, s)

	return nil
}

// SendToPluginHandler SendToPlugin Handler
func SendToPluginHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	s := PropertyInspector{}
	if err := json.Unmarshal(event.Payload, &s); err != nil {
		log.Println("ERR:", err)
		return err
	}
	log.Println("SendToPluginHandler:", s)

	settings.Save(event.Context, &s)
	return client.SetSettings(ctx, s)
}
