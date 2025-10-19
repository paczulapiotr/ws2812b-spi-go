package ws2812b

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

// Color represents an RGB color
type Color struct {
	R, G, B uint8
}

// Strip represents a WS2812B LED strip
type Strip struct {
	numLEDs int
	spiConn spi.Conn
	port    spi.PortCloser
}

// NewStrip creates a new WS2812B LED strip controller
// bus: SPI bus number
// device: SPI device number
// speedMHz: SPI speed in MHz
// numLEDs: number of LEDs in the strip
func NewStrip(bus, device, speedMHz, numLEDs int) (*Strip, error) {
	// Initialize periph.io host
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize periph.io: %w", err)
	}

	// Open SPI port
	spiPort := fmt.Sprintf("/dev/spidev%d.%d", bus, device)
	port, err := spireg.Open(spiPort)
	if err != nil {
		return nil, fmt.Errorf("failed to open SPI port %s: %w", spiPort, err)
	}

	// Configure SPI connection
	conn, err := port.Connect(physic.MegaHertz*physic.Frequency(speedMHz), spi.Mode0, 8)
	if err != nil {
		port.Close()
		return nil, fmt.Errorf("failed to connect to SPI: %w", err)
	}

	strip := &Strip{
		numLEDs: numLEDs,
		spiConn: conn,
		port:    port,
	}

	// Send initial zero byte (like Python implementation)
	if err := strip.writeBytes([]byte{0x00}); err != nil {
		port.Close()
		return nil, fmt.Errorf("failed to initialize strip: %w", err)
	}

	return strip, nil
}

// Close closes the SPI connection
func (s *Strip) Close() error {
	if s.port != nil {
		return s.port.Close()
	}
	return nil
}

// writeBytes writes data to SPI
func (s *Strip) writeBytes(data []byte) error {
	write := make([]byte, len(data))
	copy(write, data)
	read := make([]byte, len(data))
	return s.spiConn.Tx(write, read)
}

// createLEDData converts RGB colors to SPI data
// WS2812B expects GRB order
func (s *Strip) createLEDData(colors []Color) []byte {
	data := make([]byte, 0, len(colors)*24)
	
	for _, color := range colors {
		// WS2812B expects GRB order
		colorBytes := []uint8{color.G, color.R, color.B}
		
		for _, colorByte := range colorBytes {
			// Convert each color byte to 8 SPI bytes
			for bit := 7; bit >= 0; bit-- {
				if (colorByte & (1 << uint(bit))) != 0 {
					data = append(data, 0xf8) // Bit 1
				} else {
					data = append(data, 0x80) // Bit 0
				}
			}
		}
	}
	
	return data
}

// SetLED sets a single LED to a specific color
func (s *Strip) SetLED(index int, color Color) error {
	if index < 0 || index >= s.numLEDs {
		return fmt.Errorf("LED index %d out of range [0, %d)", index, s.numLEDs)
	}
	
	colors := make([]Color, index+1)
	colors[index] = color
	
	data := s.createLEDData(colors)
	return s.writeBytes(data)
}

// SetAll sets all LEDs to the same color
func (s *Strip) SetAll(color Color) error {
	colors := make([]Color, s.numLEDs)
	for i := range colors {
		colors[i] = color
	}
	return s.SetColors(colors)
}

// SetColors sets all LEDs to the specified colors
func (s *Strip) SetColors(colors []Color) error {
	if len(colors) > s.numLEDs {
		return fmt.Errorf("too many colors: %d (max %d)", len(colors), s.numLEDs)
	}
	
	data := s.createLEDData(colors)
	return s.writeBytes(data)
}

// Clear turns off all LEDs
func (s *Strip) Clear() error {
	data := make([]byte, s.numLEDs*24)
	for i := range data {
		data[i] = 0x80
	}
	return s.writeBytes(data)
}

// Fill fills a range of LEDs with a color
func (s *Strip) Fill(start, end int, color Color) error {
	if start < 0 || start >= s.numLEDs {
		return fmt.Errorf("start index %d out of range [0, %d)", start, s.numLEDs)
	}
	if end < 0 || end > s.numLEDs {
		return fmt.Errorf("end index %d out of range [0, %d]", end, s.numLEDs)
	}
	if start >= end {
		return fmt.Errorf("start index %d must be less than end index %d", start, end)
	}
	
	colors := make([]Color, end)
	for i := start; i < end; i++ {
		colors[i] = color
	}
	
	data := s.createLEDData(colors)
	return s.writeBytes(data)
}

// Chase creates a chase effect with the given color
func (s *Strip) Chase(color Color, delay time.Duration) error {
	for i := 0; i < s.numLEDs; i++ {
		colors := make([]Color, s.numLEDs)
		colors[i] = color
		
		if err := s.SetColors(colors); err != nil {
			return err
		}
		
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	return nil
}

// Rainbow creates a rainbow effect across the strip
func (s *Strip) Rainbow(delay time.Duration, loops int) error {
	for loop := 0; loop < loops; loop++ {
		for hue := 0; hue < 360; hue += 5 {
			colors := make([]Color, s.numLEDs)
			for i := 0; i < s.numLEDs; i++ {
				ledHue := (hue + (i * 360 / s.numLEDs)) % 360
				colors[i] = hsvToRGB(float64(ledHue), 1.0, 1.0)
			}
			
			if err := s.SetColors(colors); err != nil {
				return err
			}
			
			if delay > 0 {
				time.Sleep(delay)
			}
		}
	}
	return nil
}

// hsvToRGB converts HSV color to RGB
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
