package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paczulapiotr/ws2812b-spi-go"
)

func main() {
	// Initialize LED strip (SPI bus 1, device 0, 8MHz, 144 LEDs)
	strip, err := ws2812b.NewStrip(1, 0, 8, 144)
	if err != nil {
		log.Fatalf("Failed to initialize LED strip: %v", err)
	}
	defer strip.Close()

	fmt.Println("Testing first LED only...")

	// Turn on LED 0 - Red
	fmt.Println("Red")
	strip.TurnOnLED(0, ws2812b.ColorRed)
	time.Sleep(2 * time.Second)

	// Turn on LED 0 - Green
	fmt.Println("Green")
	strip.TurnOnLED(0, ws2812b.ColorGreen)
	time.Sleep(2 * time.Second)

	// Turn on LED 0 - Blue
	fmt.Println("Blue")
	strip.TurnOnLED(0, ws2812b.ColorBlue)
	time.Sleep(2 * time.Second)

	// Turn off LED 0
	fmt.Println("Off")
	strip.TurnOffLED(0)
	time.Sleep(1 * time.Second)

	fmt.Println("Done!")
}
