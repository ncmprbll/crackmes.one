package main

import "fmt"

func main() {
	var4E := []byte{0x9e, 0x81, 0x9b, 0xf3, 0xd0, 0xd6, 0x03, 0xb4, 0xe9, 0x27, 0x00, 0x66, 0x2b, 0x05}
	var30 := []byte{0x0e, 0x03, 0x0b, 0x08, 0x01, 0x02, 0x0d, 0x04, 0x08, 0x0a, 0x06, 0x0c, 0x05, 0x09, 0x07, 0x0f}
	var52 := []byte{0x5a, 0x43, 0xc2, 0xf8}

	for k, v := range var4E {
		var4E[k] = v ^ 0xA5
	}

	output := var4E

	for i := 0; i < 0xd; i += 2 {
		var7B, var7A := output[i], output[i+1]

		for var6C := 5; var6C >= 0; var6C-- {
			var76 := var7A
			var7A = var7B

			var77 := var30[var7A&0xF] | var30[var7A/16]*16
			c := var52[var6C&0x3]

			var77 = var77 + c
			var77 = var77 ^ ((var77 / 32) | (var77 * 8))

			var7B = var76 ^ var77
		}

		output[i], output[i+1] = var7B, var7A
	}

	fmt.Println(string(output))
}
