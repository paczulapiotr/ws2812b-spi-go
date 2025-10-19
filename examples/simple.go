package main

import (
	"log"
	"time"

	"github.com/paczulapiotr/ws2812b-spi-go"
)

func main() {
	// Initialize LED strip (SPI bus 1, device 0, 8MHz, 8 LEDs)
	strip, err := ws2812b.NewStrip(1, 0, 8, 8)
	if err != nil {
		log.Fatalf("Failed to initialize LED strip: %v", err)
	}
	defer strip.Close()

	// Turn on LED 0 - Red
	strip.TurnOnLED(0, ws2812b.ColorRed)
	time.Sleep(1 * time.Second)

	// Turn on LED 1 - Green
	strip.TurnOnLED(1, ws2812b.ColorGreen)
	time.Sleep(1 * time.Second)

	// Turn on LED 2 - Blue
	strip.TurnOnLED(2, ws2812b.ColorBlue)
	time.Sleep(1 * time.Second)

	// Turn off LED 0
	strip.TurnOffLED(0)
	time.Sleep(1 * time.Second)

	// Set all LEDs to yellow
	strip.SetAll(ws2812b.ColorYellow)
	strip.Show()
	time.Sleep(2 * time.Second)

	// Turn off all LEDs
	strip.Clear()
	strip.Show()
}
