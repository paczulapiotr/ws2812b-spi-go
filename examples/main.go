package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paczulapiotr/ws2812b-spi-go"
)

func main() {
	// Create a new LED strip with:
	// - SPI bus 1
	// - SPI device 0
	// - 8 MHz speed
	// - 8 LEDs
	strip, err := ws2812b.NewStrip(0, 0, 8, 8)
	if err != nil {
		log.Fatalf("Failed to create LED strip: %v", err)
	}
	defer strip.Close()

	// Example 1: Turn on individual LEDs with different colors
	fmt.Println("Example 1: Individual LED control")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 2: Turn off specific LEDs
	fmt.Println("Example 2: Turning off individual LEDs")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 3: Set all LEDs at once
	fmt.Println("Example 3: Set all LEDs")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 4: Clear all LEDs
	fmt.Println("Example 4: Clear all LEDs")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 5: Fill a range
	fmt.Println("Example 5: Fill range of LEDs")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 6: Chase effect
	fmt.Println("Example 6: Chase effect")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 7: Rainbow effect
	fmt.Println("Example 7: Rainbow effect")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 8: Manual control with batch update
	fmt.Println("Example 8: Manual batch control")
	/// ... implementation
	time.Sleep(1 * time.Second)

	// Example 9: Light LEDs in order (like Python example)
	fmt.Println("Example 9: Light in order")
	/// ... implementation
	time.Sleep(1 * time.Second)
}
