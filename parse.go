// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import "strings"

const phoneE164Length = 8

// ParseE164 validates HotelBeds-styled phone number and converts into international E164 phone number.
func ParseE163(raw string) string {
	var e164Number string
	delimPos := strings.IndexAny(raw, ".-, ")
	if delimPos > 0 {
		char := raw[delimPos]
		if char != 0 { // EOF
			e164Number = strings.ReplaceAll(raw, string(char), "")
		} else {
			e164Number = raw
		}
	} else {
		e164Number = raw
	}
	if len(e164Number) < phoneE164Length {
		return ""
	}

	var formattedNumber string
	if strings.HasPrefix(e164Number, "+00") {
		formattedNumber = e164Number[3:] // international
	} else if strings.HasPrefix(e164Number, "00") {
		formattedNumber = e164Number[2:] // international VoIP
	} else {
		formattedNumber = e164Number
	}

	// Check if have sign in result, otherwise append it
	if formattedNumber[0] != '+' {
		formattedNumber = "+" + formattedNumber
	}

	return formattedNumber
}
