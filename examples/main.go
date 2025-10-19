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
	strip, err := ws2812b.NewStrip(1, 0, 8, 8)
	if err != nil {
		log.Fatalf("Failed to create LED strip: %v", err)
	}
	defer strip.Close()

	// Example 1: Turn on individual LEDs with different colors
	// fmt.Println("Example 1: Individual LED control")
	// strip.TurnOnLED(0, ws2812b.ColorRed)
	// time.Sleep(500 * time.Millisecond)

	// strip.TurnOnLED(1, ws2812b.ColorGreen)
	// time.Sleep(500 * time.Millisecond)

	// strip.TurnOnLED(2, ws2812b.ColorBlue)
	// time.Sleep(500 * time.Millisecond)

	// // Custom color
	// strip.TurnOnLED(3, ws2812b.Color{R: 220, G: 20, B: 60}) // Crimson
	// time.Sleep(500 * time.Millisecond)

	// time.Sleep(2 * time.Second)

	// // Example 2: Turn off specific LEDs
	// fmt.Println("Example 2: Turning off individual LEDs")
	// strip.TurnOffLED(0)
	// time.Sleep(500 * time.Millisecond)
	// strip.TurnOffLED(1)
	// time.Sleep(500 * time.Millisecond)

	// // Example 3: Set all LEDs at once
	// fmt.Println("Example 3: Set all LEDs")
	// strip.SetAll(ws2812b.ColorYellow)
	// strip.Show()
	// time.Sleep(2 * time.Second)

	// // Example 4: Clear all LEDs
	// fmt.Println("Example 4: Clear all LEDs")
	// strip.Clear()
	// strip.Show()
	// time.Sleep(1 * time.Second)

	// // Example 5: Fill a range
	// fmt.Println("Example 5: Fill range of LEDs")
	// strip.Fill(2, 5, ws2812b.ColorPurple)
	// strip.Show()
	// time.Sleep(2 * time.Second)

	// // Example 6: Chase effect
	// fmt.Println("Example 6: Chase effect")
	// strip.Clear()
	// strip.Show()
	// time.Sleep(500 * time.Millisecond)
	// strip.Chase(ws2812b.ColorCyan, 100*time.Millisecond, 3)

	// Example 7: Rainbow effect
	fmt.Println("Example 7: Rainbow effect")
	strip.Rainbow(10 * time.Millisecond)

	// // Example 8: Manual control with batch update
	// fmt.Println("Example 8: Manual batch control")
	// strip.Clear()
	// strip.SetLED(0, ws2812b.ColorRed)
	// strip.SetLED(2, ws2812b.ColorGreen)
	// strip.SetLED(4, ws2812b.ColorBlue)
	// strip.SetLED(6, ws2812b.ColorWhite)
	// strip.Show() // Update all at once
	// time.Sleep(2 * time.Second)

	// // Example 9: Light LEDs in order (like Python example)
	// fmt.Println("Example 9: Light in order")
	// strip.Clear()
	// strip.Show()
	// time.Sleep(500 * time.Millisecond)

	// for i := 0; i < 8; i++ {
	// 	strip.SetLED(i, ws2812b.Color{R: 255, G: uint8((i * 32)), B: uint8((i * 16))})
	// 	strip.Show()
	// 	time.Sleep(500 * time.Millisecond)
	// }

	time.Sleep(2 * time.Second)

	// Turn off all LEDs
	fmt.Println("Turning off all LEDs")
	strip.Clear()
	strip.Show()
}
