package bit

import ()

func parse(b uint8) [8]bool {
	var (
		a [8]bool
		c uint8
	)
	c = b
	for i := 7; i >= 0; {
		if c%2 == 0 {
			a[i] = false
		} else {
			a[i] = true
		}
		c = c >> 1
		i--
	}
	return a
}
func transfer(b [8]bool) uint8 {
	var re uint8
	for i := 0; i < 8; {
		re *= 2
		if b[i] {
			re += 1
		}
		i++
	}
	return re
}
