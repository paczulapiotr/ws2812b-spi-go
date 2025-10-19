package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paczulapiotr/ws2812b-spi-go"
)

func main() {
	// Create a new LED strip with:
	// - SPI bus 0
	// - SPI device 0
	// - 8 MHz speed
	// - 148 LEDs
	strip, err := ws2812b.NewStrip(0, 0, 8, 148)
	if err != nil {
		log.Fatalf("Failed to create LED strip: %v", err)
	}
	defer strip.Close()

	// Example 1: Set all LEDs at once
	fmt.Println("Example 1: Set all LEDs")
	if err := strip.SetAll(ws2812b.Color{R: 255, G: 0, B: 255}); err != nil {
		log.Printf("Error setting all LEDs: %v", err)
	}
	time.Sleep(1 * time.Second)
	strip.Clear()

	// Example 2: Clear all LEDs
	fmt.Println("Example 2: Clear all LEDs")
	if err := strip.Clear(); err != nil {
		log.Printf("Error clearing LEDs: %v", err)
	}
	time.Sleep(1 * time.Second)
	strip.Clear()

	// Example 3: Fill a range
	fmt.Println("Example 3: Fill range of LEDs")
	if err := strip.Fill(0, 20, ws2812b.Color{R: 255, G: 0, B: 0}); err != nil {
		log.Printf("Error filling range: %v", err)
	}
	time.Sleep(1 * time.Second)
	strip.Clear()

	// Example 4: Chase effect
	fmt.Println("Example 4: Chase effect")
	if err := strip.Chase(ws2812b.Color{R: 0, G: 0, B: 255}, 100*time.Millisecond); err != nil {
		log.Printf("Error running chase: %v", err)
	}
	time.Sleep(1 * time.Second)
	strip.Clear()

	// Example 5: Rainbow effect
	fmt.Println("Example 5: Rainbow effect")
	if err := strip.Rainbow(50*time.Millisecond, 2); err != nil {
		log.Printf("Error running rainbow: %v", err)
	}
	time.Sleep(1 * time.Second)
	strip.Clear()

	// Example 6: Manual control with batch update
	fmt.Println("Example 6: Manual batch control")
	colors := []ws2812b.Color{
		{R: 255, G: 0, B: 0},     // Red
		{R: 255, G: 127, B: 0},   // Orange
		{R: 255, G: 255, B: 0},   // Yellow
		{R: 0, G: 255, B: 0},     // Green
		{R: 0, G: 0, B: 255},     // Blue
		{R: 75, G: 0, B: 130},    // Indigo
		{R: 148, G: 0, B: 211},   // Violet
		{R: 255, G: 255, B: 255}, // White
	}
	if err := strip.SetColors(colors); err != nil {
		log.Printf("Error setting colors: %v", err)
	}
	time.Sleep(1 * time.Second)
}
