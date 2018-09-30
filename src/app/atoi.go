// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const fastSmalls = true // enable fast path for small integers
const nSmalls = 100

var smallsString = []byte{
	'0', '0', '0', '1', '0', '2', '0', '3', '0', '4', '0', '5', '0', '6', '0', '7', '0', '8', '0', '9',
	'1', '0', '1', '1', '1', '2', '1', '3', '1', '4', '1', '5', '1', '6', '1', '7', '1', '8', '1', '9',
	'2', '0', '2', '1', '2', '2', '2', '3', '2', '4', '2', '5', '2', '6', '2', '7', '2', '8', '2', '9',
	'3', '0', '3', '1', '3', '2', '3', '3', '3', '4', '3', '5', '3', '6', '3', '7', '3', '8', '3', '9',
	'4', '0', '4', '1', '4', '2', '4', '3', '4', '4', '4', '5', '4', '6', '4', '7', '4', '8', '4', '9',
	'5', '0', '5', '1', '5', '2', '5', '3', '5', '4', '5', '5', '5', '6', '5', '7', '5', '8', '5', '9',
	'6', '0', '6', '1', '6', '2', '6', '3', '6', '4', '6', '5', '6', '6', '6', '7', '6', '8', '6', '9',
	'7', '0', '7', '1', '7', '2', '7', '3', '7', '4', '7', '5', '7', '6', '7', '7', '7', '8', '7', '9',
	'8', '0', '8', '1', '8', '2', '8', '3', '8', '4', '8', '5', '8', '6', '8', '7', '8', '8', '8', '9',
	'9', '0', '9', '1', '9', '2', '9', '3', '9', '4', '9', '5', '9', '6', '9', '7', '9', '8', '9', '9',
}

const host32bit = ^uint(0)>>32 == 0

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) []byte { // WARN:
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

func FormatInt(i int64) []byte {
	if 0 <= i && i < nSmalls {
		return small(int(i))
	}
	return formatBits(nil, uint64(i), i < 0)
}

func formatBits(dst []byte, u uint64, neg bool) []byte {
	var a [64 + 1]byte
	i := len(a)

	if neg {
		u = -u
	}

	// convert bits
	// We use uint values where we can because those will
	// fit into a single register even on a 32bit machine.

	if host32bit {
		// convert the lower digits using 32bit operations
		for u >= 1e9 {
			// Avoid using r = a%b in addition to q = a/b
			// since 64bit division and modulo operations
			// are calculated by runtime functions on 32bit machines.
			q := u / 1e9
			us := uint(u - q*1e9) // u % 1e9 fits into a uint
			for j := 4; j > 0; j-- {
				is := us % 100 * 2
				us /= 100
				i -= 2
				a[i+1] = smallsString[is+1]
				a[i+0] = smallsString[is+0]
			}

			// us < 10, since it contains the last digit
			// from the initial 9-digit us.
			i--
			a[i] = smallsString[us*2+1]

			u = q
		}
		// u < 1e9
	}

	// u guaranteed to fit into a uint
	us := uint(u)
	for us >= 100 {
		is := us % 100 * 2
		us /= 100
		i -= 2
		a[i+1] = smallsString[is+1]
		a[i+0] = smallsString[is+0]
	}

	// us < 100
	is := us * 2
	i--
	a[i] = smallsString[is+1]
	if us >= 10 {
		i--
		a[i] = smallsString[is]
	}

	// add sign, if any
	if neg {
		i--
		a[i] = '-'
	}
	return a[i:]
}
