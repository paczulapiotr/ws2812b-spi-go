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
// speedMHz: SPI speed in MHz (ignored, uses 6.4 MHz for WS2812B timing)
// numLEDs: number of LEDs in the strip
func NewStrip(bus, device, speedMHz, numLEDs int) (*Strip, error) {
	// WS2812B requires 6.4 MHz SPI speed for proper timing
	// This gives us 156.25ns per bit which works for WS2812B protocol
	dev, err := spi.Open(&spi.Devfs{
		Dev:      fmt.Sprintf("/dev/spidev%d.%d", bus, device),
		Mode:     spi.Mode0,
		MaxSpeed: 6400000, // 6.4 MHz
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open SPI device: %w", err)
	}

	strip := &Strip{
		device:   dev,
		numLEDs:  numLEDs,
		leds:     make([]Color, numLEDs),
		spiSpeed: 6400000,
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
	// Add a small delay for the reset period (>50µs)
	time.Sleep(100 * time.Microsecond)
	return nil
}

// encode converts the LED color values to SPI data
// WS2812B protocol requires specific timing:
// - Bit 1: 0.8µs high, 0.45µs low (total 1.25µs)
// - Bit 0: 0.4µs high, 0.85µs low (total 1.25µs)
// At 6.4 MHz SPI, each bit is 156.25ns
// So we encode each WS2812B bit as 8 SPI bits (~1.25µs)
// WS2812B uses GRB color order (Green, Red, Blue)
func (s *Strip) encode() []byte {
	// Reset signal: at least 50µs of low (we'll use ~100µs)
	// At 6.4MHz, that's about 640 bits, so 80 bytes
	resetBytes := 80
	
	// Each LED color bit becomes 8 SPI bits
	// 24 color bits per LED * 8 SPI bits = 192 SPI bits = 24 bytes per LED
	bytesPerLED := 24
	totalBytes := resetBytes + (s.numLEDs * bytesPerLED) + resetBytes
	
	data := make([]byte, totalBytes)
	
	// Start after reset bytes
	offset := resetBytes

	for _, led := range s.leds {
		// WS2812B expects GRB order
		colors := []uint8{led.G, led.R, led.B}
		
		for _, colorByte := range colors {
			// Process each bit from MSB to LSB
			for bit := 7; bit >= 0; bit-- {
				if (colorByte & (1 << bit)) != 0 {
					// Bit is 1: 110000 pattern (5 highs, 3 lows) = 0b11111000 = 0xF8
					data[offset] = 0xF8
				} else {
					// Bit is 0: 100000 pattern (2 highs, 6 lows) = 0b11000000 = 0xC0
					data[offset] = 0xC0
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
