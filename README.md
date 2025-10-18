# WS2812B LED Strip Control Library for Go

A Go library for controlling WS2812B (NeoPixel) LED strips via SPI on Linux systems (like Raspberry Pi, Radxa boards, etc.).

## Features

- Control individual LEDs with specific RGB colors
- Turn LEDs on/off individually or all at once
- Set color ranges on the strip
- Built-in color constants (Red, Green, Blue, Yellow, Cyan, Magenta, White, Purple, Pink)
- Pre-built animation effects:
  - Rainbow effect
  - Chase effect
- Efficient batch updates
- Easy-to-use API

## Hardware Requirements

- Linux-based SBC (Single Board Computer) like Raspberry Pi or Radxa board
- WS2812B LED strip
- SPI interface enabled on your device

## Installation

First, make sure you have Go installed and SPI enabled on your device.

Install the library:

```bash
cd go-ws2812b
go mod download
```

## Project Structure

```
go-ws2812b/
├── ws2812b.go      # Main library file
├── go.mod          # Go module definition
├── examples/
│   └── main.go     # Example usage
└── README.md       # This file
```

## Usage

### Basic Example

```go
package main

import (
    "log"
    "time"
    "github.com/radxa/sample_code/go-ws2812b"
)

func main() {
    // Create a new LED strip with 8 LEDs
    // Parameters: SPI bus, device, speed (MHz), number of LEDs
    strip, err := ws2812b.NewStrip(1, 0, 8, 8)
    if err != nil {
        log.Fatal(err)
    }
    defer strip.Close()

    // Turn on first LED with red color
    strip.TurnOnLED(0, ws2812b.ColorRed)
    
    // Turn on second LED with custom color
    strip.TurnOnLED(1, ws2812b.Color{R: 220, G: 20, B: 60})
    
    time.Sleep(2 * time.Second)
    
    // Turn off first LED
    strip.TurnOffLED(0)
    
    // Set all LEDs to blue
    strip.SetAll(ws2812b.ColorBlue)
    strip.Show()
    
    time.Sleep(2 * time.Second)
    
    // Clear all LEDs
    strip.Clear()
    strip.Show()
}
```

## API Reference

### Creating a Strip

```go
strip, err := ws2812b.NewStrip(bus, device, speedMHz, numLEDs)
```

- `bus`: SPI bus number (typically 0 or 1)
- `device`: SPI device number (usually 0)
- `speedMHz`: SPI speed in MHz (typically 8)
- `numLEDs`: Total number of LEDs in your strip

### Individual LED Control

```go
// Turn on a specific LED with a color
strip.TurnOnLED(index, color)

// Turn off a specific LED
strip.TurnOffLED(index)

// Set LED color (doesn't update strip immediately)
strip.SetLED(index, color)

// Get current color of an LED
color, err := strip.GetLED(index)
```

### Bulk Operations

```go
// Set all LEDs to the same color
strip.SetAll(color)

// Clear all LEDs (turn them off)
strip.Clear()

// Fill a range of LEDs
strip.Fill(startIndex, endIndex, color)

// Update the physical strip with current values
strip.Show()
```

### Animation Effects

```go
// Rainbow effect (cycles through all hues)
strip.Rainbow(delay time.Duration)

// Chase effect (single LED moving across strip)
strip.Chase(color, delay time.Duration, iterations int)
```

### Color Constants

Pre-defined colors available:
- `ColorOff` (black/off)
- `ColorRed`
- `ColorGreen`
- `ColorBlue`
- `ColorYellow`
- `ColorCyan`
- `ColorMagenta`
- `ColorWhite`
- `ColorPurple`
- `ColorPink`

### Custom Colors

```go
customColor := ws2812b.Color{R: 255, G: 128, B: 0} // Orange
strip.TurnOnLED(0, customColor)
```

## Advanced Usage

### Batch Updates for Better Performance

When updating multiple LEDs, use `SetLED()` followed by a single `Show()` call:

```go
// Efficient - one SPI transaction
strip.SetLED(0, ws2812b.ColorRed)
strip.SetLED(1, ws2812b.ColorGreen)
strip.SetLED(2, ws2812b.ColorBlue)
strip.Show()

// Less efficient - three SPI transactions
strip.TurnOnLED(0, ws2812b.ColorRed)
strip.TurnOnLED(1, ws2812b.ColorGreen)
strip.TurnOnLED(2, ws2812b.ColorBlue)
```

## Running the Example

```bash
cd go-ws2812b/examples
go run main.go
```

## Technical Details

### WS2812B Protocol

The WS2812B uses a single-wire protocol that can be emulated using SPI:
- Each LED requires 24 bits of data (8 bits per color channel)
- Color order is GRB (Green, Red, Blue) - the library handles this automatically
- Each bit is encoded as either `0x80` (representing 0) or `0xf8` (representing 1)
- SPI Mode 1 is used
- Typical speed: 8 MHz

### SPI Configuration

Make sure SPI is enabled on your device. On Raspberry Pi/Radxa:

```bash
# Check if SPI device exists
ls -l /dev/spidev*

# The device should appear as /dev/spidev1.0 or similar
```

## Troubleshooting

### Permission Denied Error

If you get a permission error when opening the SPI device, you may need to run with sudo or add your user to the appropriate group:

```bash
sudo usermod -a -G spi $USER
# Log out and back in for changes to take effect
```

Or run with sudo:
```bash
sudo go run main.go
```

### LEDs Not Lighting Up

1. Check your wiring (Data In, VCC, GND)
2. Verify SPI is enabled on your device
3. Ensure you're using the correct SPI bus and device numbers
4. Check that your power supply is adequate for your LED strip
5. Try a lower SPI speed (e.g., 4 MHz instead of 8 MHz)

### Colors Appear Wrong

- WS2812B uses GRB color order (not RGB) - the library handles this automatically
- Check your power supply voltage (WS2812B typically requires 5V)

### Module Not Found

Make sure you're in the correct directory and have run:
```bash
go mod download
```

## Comparison with Python Version

This Go library provides the same functionality as the Python reference implementation:

**Python:**
```python
led = WS2812B(1, 0, 8)
led.light_in_order()
led.all_led_off()
```

**Go:**
```go
strip, _ := ws2812b.NewStrip(1, 0, 8, 8)
defer strip.Close()
// ... control individual LEDs ...
strip.Clear()
strip.Show()
```

## License

This library is provided as-is for educational and commercial use.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
