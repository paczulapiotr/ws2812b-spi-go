package ws2812b

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/spi"
)

// Color represents an RGB color value
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Common color constants
var (
	ColorOff     = Color{0, 0, 0}
	ColorRed     = Color{255, 0, 0}
	ColorGreen   = Color{0, 255, 0}
	ColorBlue    = Color{0, 0, 255}
	ColorYellow  = Color{255, 255, 0}
	ColorCyan    = Color{0, 255, 255}
	ColorMagenta = Color{255, 0, 255}
	ColorWhite   = Color{255, 255, 255}
	ColorPurple  = Color{128, 0, 128}
	ColorPink    = Color{255, 192, 203}
)

// Strip represents a WS2812B LED strip
type Strip struct {
	device   *spi.Device
	numLEDs  int
	leds     []Color
	spiSpeed int
}

// NewStrip creates a new WS2812B LED strip controller
// bus: SPI bus number (e.g., 0 or 1)
// device: SPI device number (usually 0)
// speedMHz: SPI speed in MHz (recommended: 8)
// numLEDs: number of LEDs in the strip
func NewStrip(bus, device, speedMHz, numLEDs int) (*Strip, error) {
	dev, err := spi.Open(&spi.Devfs{
		Dev:      fmt.Sprintf("/dev/spidev%d.%d", bus, device),
		Mode:     spi.Mode0, // Mode 0 for Radxa Rock 5B
		MaxSpeed: int64(speedMHz * 1000000),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open SPI device: %w", err)
	}

	// Send initialization byte
	dev.Tx([]byte{0x00}, nil)

	strip := &Strip{
		device:   dev,
		numLEDs:  numLEDs,
		leds:     make([]Color, numLEDs),
		spiSpeed: speedMHz * 1000000,
	}

	// Initialize all LEDs to off
	strip.Clear()
	strip.Show()

	return strip, nil
}

// Close closes the SPI connection
func (s *Strip) Close() error {
	if s.device != nil {
		return s.device.Close()
	}
	return nil
}

// SetLED sets the color of a specific LED (0-indexed)
func (s *Strip) SetLED(index int, color Color) error {
	if index < 0 || index >= s.numLEDs {
		return fmt.Errorf("LED index %d out of range (0-%d)", index, s.numLEDs-1)
	}
	s.leds[index] = color
	return nil
}

// GetLED returns the current color of a specific LED
func (s *Strip) GetLED(index int) (Color, error) {
	if index < 0 || index >= s.numLEDs {
		return ColorOff, fmt.Errorf("LED index %d out of range (0-%d)", index, s.numLEDs-1)
	}
	return s.leds[index], nil
}

// TurnOnLED turns on a specific LED with the given color
func (s *Strip) TurnOnLED(index int, color Color) error {
	if err := s.SetLED(index, color); err != nil {
		return err
	}
	return s.Show()
}

// TurnOffLED turns off a specific LED
func (s *Strip) TurnOffLED(index int) error {
	return s.TurnOnLED(index, ColorOff)
}

// SetAll sets all LEDs to the same color
func (s *Strip) SetAll(color Color) {
	for i := range s.leds {
		s.leds[i] = color
	}
}

// Clear turns off all LEDs (sets them to black)
func (s *Strip) Clear() {
	s.SetAll(ColorOff)
}

// Show updates the LED strip with the current color values
func (s *Strip) Show() error {
	data := s.encode()
	if err := s.device.Tx(data, nil); err != nil {
		return fmt.Errorf("failed to send data to LED strip: %w", err)
	}
	return nil
}

// encode converts the LED color values to SPI data
// This uses the proven encoding from the working Python implementation
// WS2812B uses GRB color order (Green, Red, Blue)
// 0x80 = bit 0, 0xF8 = bit 1
func (s *Strip) encode() []byte {
	// Each LED needs 24 bytes (8 bytes per color channel * 3 channels)
	data := make([]byte, s.numLEDs*24)
	
	offset := 0
	for _, led := range s.leds {
		// WS2812B expects GRB order
		colors := []uint8{led.G, led.R, led.B}
		
		for _, colorByte := range colors {
			// Convert each color byte to 8 SPI bytes
			// Process each bit from MSB to LSB
			for bit := 7; bit >= 0; bit-- {
				if (colorByte & (1 << bit)) != 0 {
					// Bit is 1
					data[offset] = 0xF8
				} else {
					// Bit is 0
					data[offset] = 0x80
				}
				offset++
			}
		}
	}

	return data
}

// Rainbow creates a rainbow effect across the strip
func (s *Strip) Rainbow(delay time.Duration) error {
	for hue := 0; hue < 360; hue++ {
		for i := 0; i < s.numLEDs; i++ {
			ledHue := (hue + (i * 360 / s.numLEDs)) % 360
			s.leds[i] = hsvToRGB(float64(ledHue), 1.0, 1.0)
		}
		if err := s.Show(); err != nil {
			return err
		}
		time.Sleep(delay)
	}
	return nil
}

// Fill sets a range of LEDs to a specific color
func (s *Strip) Fill(start, end int, color Color) error {
	if start < 0 || start >= s.numLEDs {
		return fmt.Errorf("start index %d out of range (0-%d)", start, s.numLEDs-1)
	}
	if end < 0 || end >= s.numLEDs {
		return fmt.Errorf("end index %d out of range (0-%d)", end, s.numLEDs-1)
	}
	if start > end {
		return fmt.Errorf("start index %d is greater than end index %d", start, end)
	}

	for i := start; i <= end; i++ {
		s.leds[i] = color
	}
	return nil
}

// Chase creates a chase effect with the given color
func (s *Strip) Chase(color Color, delay time.Duration, iterations int) error {
	for iter := 0; iter < iterations; iter++ {
		for i := 0; i < s.numLEDs; i++ {
			s.Clear()
			s.leds[i] = color
			if err := s.Show(); err != nil {
				return err
			}
			time.Sleep(delay)
		}
	}
	return nil
}

// hsvToRGB converts HSV color space to RGB
// h: 0-360, s: 0-1, v: 0-1
func hsvToRGB(h, s, v float64) Color {
	h = h / 60.0
	i := int(h)
	f := h - float64(i)
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return Color{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
	}
}
