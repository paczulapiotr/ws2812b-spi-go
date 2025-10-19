#!/usr/bin/env python
# -*- encoding: utf-8 -*-

import spidev
from time import sleep

class WS2812B:

	def __init__(self, bus, device, speed, num_leds):
		self.num_leds = num_leds
		self.spi = spidev.SpiDev()
		self.spi.open(bus, device)
		self.spi.max_speed_hz = int(speed * 1000000)
		self.spi.mode = 0b0  # Mode 0 for Rock 5B
		self.spi.xfer([0x00])

	def create_led_data(self, colors):
		"""
		Create SPI data for LEDs
		colors: list of (R, G, B) tuples
		Returns: list of bytes for SPI transmission
		"""
		data = []
		for r, g, b in colors:
			# WS2812B expects GRB order
			for color_byte in [g, r, b]:
				# Convert each color byte to 8 SPI bytes
				for bit in range(7, -1, -1):
					if (color_byte & (1 << bit)) != 0:
						data.append(0xf8)  # Bit 1
					else:
						data.append(0x80)  # Bit 0
		return data

	def light_in_order(self, delay=0.1, loops=3):
		"""
		Light up LEDs one by one in red, green, blue pattern
		delay: time between each LED (seconds)
		loops: number of times to repeat the pattern
		"""
		print(f"Lighting {self.num_leds} LEDs in order")
		
		# Define color pattern: Red, Green, Blue repeating
		colors_pattern = [(255, 0, 0), (0, 255, 0), (0, 0, 255)]
		
		for loop in range(loops):
			print(f"Loop {loop + 1}/{loops}")
			
			# Light up LEDs one by one
			for i in range(self.num_leds):
				# Create list of colors for LEDs 0 to i
				colors = [colors_pattern[j % 3] for j in range(i + 1)]
				led_data = self.create_led_data(colors)
				self.spi.xfer(led_data)
				sleep(delay)
			
			# Turn off all LEDs
			print("Turning off all LEDs")
			off_data = [0x80] * (self.num_leds * 24)
			self.spi.xfer(off_data)
			
			if loop < loops - 1:
				sleep(0.5)  # Pause between loops

	def close(self):
		self.spi.close()

if __name__ == "__main__":
	led = WS2812B(bus=0, device=0, speed=8, num_leds=148)
	try:
		led.light_in_order(delay=0, loops=3)  # 10ms delay, 3 loops
	finally:
		led.close()

