// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const fastSmalls = true // enable fast path for small integers
const nSmalls = 100
const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"
const host32bit = ^uint(0)>>32 == 0
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) string { // WARN:
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

func FormatInt(i int64) []byte {
	if 0 <= i && i < nSmalls {
		return []byte(small(int(i)))
	}
	return formatBits(nil, uint64(i), i < 0, false)
}

func formatBits(dst []byte, u uint64, neg, append_ bool) []byte {

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

	if append_ {
		return append(dst, a[i:]...)
	}
	return a[i:]
}
