package ndocid

import (
	"fmt"
	"math/bits"
	"os"
	"strings"
	"time"
	"unicode"
)

// DateFormat referencing Mon Jan 2 15:04:05 -0700 MST 2006
const DateFormat = "20060102150405"

const customBase32Alphabet = string("23456789ABCDEFHIJKLMNOPQRTUVWXYZ")

var Verbose bool
var verboseLineOut = func(format string, msg ...interface{}) {
	if Verbose {
		fmt.Fprintf(os.Stdout, format+"\n", msg...)
	}
}

var baseLocation = time.Local

// encodes an integer in range [0,32)
func customBase32Encode(i int) rune {
	return rune(customBase32Alphabet[i])
}

func customBase32Decode(r rune) (i int, ok bool) {
	var decoding = map[rune]int{
		'2': 0,
		'3': 1,
		'4': 2,
		'5': 3, 'S': 3,
		'6': 4, 'G': 4,
		'7': 5,
		'8': 6,
		'9': 7,
		'A': 8,
		'B': 9,
		'C': 10,
		'D': 11,
		'E': 12,
		'F': 13,
		'H': 14,
		'I': 15, '1': 15,
		'J': 16,
		'K': 17,
		'L': 18,
		'M': 19,
		'N': 20,
		'O': 21, '0': 21,
		'P': 22,
		'Q': 23,
		'R': 24,
		'T': 25,
		'U': 26,
		'V': 27,
		'W': 28,
		'X': 29,
		'Y': 30,
		'Z': 31,
	}
	i, ok = decoding[unicode.ToUpper(r)]
	return
}

func EncodeUint64(i uint64) string {
	verboseLineOut("Received numeric input: %d", i)
	return encode(i)
}

func EncodeBitstring(s string) (result string, err error) {
	verboseLineOut("Received bitstring input: %s", s)
	var number uint64
	if s == "" {
		err = fmt.Errorf("Empty bitstring input")
		return
	}
	shiftLeft := func() {
		if bits.LeadingZeros64(number) == 0 {
			err = fmt.Errorf("Bitstring input exceeds 64 bits")
		}
		number <<= 1
	}
	for _, char := range s {
		switch char {
		case ' ':
		case '\t':
		case '_':
		case '0':
			shiftLeft()
		case '1':
			shiftLeft()
			number++
		default:
			err = fmt.Errorf("Bad character in bitstring input: %c (%U)", char, char)
			return
		}
	}
	result = encode(number)
	return
}

func EncodeDatetime(s string) (result string, err error) {
	if len(s) != len(DateFormat) {
		err = fmt.Errorf("Input date does not match required %d-character-format (see -h)", len(DateFormat))
		return
	}
	t, err := time.ParseInLocation(DateFormat, s, baseLocation)
	if err != nil {
		err = fmt.Errorf("Bad date format: %s", err)
		return
	}
	unixSeconds := t.Unix()
	verboseLineOut("Received date input: %s (unix time in seconds: %d)", t.Format(time.RFC1123Z), unixSeconds)
	result = EncodeUint64(uint64(unixSeconds))
	return
}

func encode(x uint64) (r string) {
	var acc strings.Builder

	//Algorithm
	f1 := int(x & 0b000000000111 >> 0)
	f2 := int(x & 0b000000111000 >> 3)
	f3 := int(x & 0b000111000000 >> 6)
	f4 := int(x & 0b111000000000 >> 9)
	p2 := bits.OnesCount64(x&0b000000111111) & 0b1
	p3 := bits.OnesCount64(x&0b000111111111) & 0b1
	p4 := bits.OnesCount64(x&0b111111111111) & 0b1
	fc := p2<<2 + p3<<1 + p4
	fp := [...]int{fc, f1, f2, f3, f4}
	fs := 3*fc + 1*f1 + 3*f2 + 1*f3 + 3*f4

	vp := make([]int, 0, 11)
	vs := 0
	for vr := x >> 12; vr > 0; vr >>= 5 {
		n := int(vr & 0b11111)
		vp = append(vp, n)
		vs += n * (1 + len(vp)%2*2)
	}

	mc := 29 - (fs+vs)%29

	for _, i := range fp {
		acc.WriteRune(customBase32Encode(i))
	}
	acc.WriteRune(customBase32Encode(mc))
	for _, i := range vp {
		acc.WriteRune(customBase32Encode(i))
	}
	r = acc.String()

	if Verbose {
		//Verbose output
		var th, tb string
		for i := bits.LeadingZeros64(x) / 8; i < 8; i++ {
			b := x >> ((7 - i) * 8) & 0xFF
			th += fmt.Sprintf("%02X ", b)
			tb += fmt.Sprintf("%08b ", b)
		}
		verboseLineOut("Encoding %d ( %s/ %s):", x, th, tb)
		verboseLineOut("  Calculating leading [F]ixed [P]art FP:")
		verboseLineOut("    F1 := LSB  3 to  1: %03b (%1d)", f1, f1)
		verboseLineOut("    F2 := LSB  6 to  4: %03b (%1d)", f2, f2)
		verboseLineOut("    F3 := LSB  9 to  7: %03b (%1d)", f3, f3)
		verboseLineOut("    F4 := LSB 12 to 10: %03b (%1d)", f4, f4)
		verboseLineOut("    P2 := Even [P]arity bit for  6 LSB := %1b", p2)
		verboseLineOut("    P3 := Even [P]arity bit for  9 LSB := %1b", p3)
		verboseLineOut("    P4 := Even [P]arity bit for 12 LSB := %1b", p4)
		verboseLineOut("    FC := P2 * 4 + P3 * 2 + P4 : %03b (%1d) as input [C]heck digit", fc, fc)
		verboseLineOut("    < FP := [FC F1 F2 F3 F4]: %v", fp)
		verboseLineOut("  Calculating trailing [V]ariable [P]art VP:")
		verboseLineOut("    %d bits remaining: %b", bits.Len64(x>>12), x>>12)
		for i, v := range vp {
			verboseLineOut("    V%d := next 5 LSB: %05b (%2d)", i+1, v, v)
		}
		verboseLineOut("    < VP := [V1 V2 ...]: %v", vp)
		verboseLineOut("  Calculating [M]aster [C]heck digit MC:")
		verboseLineOut("    FS := [F]ixed    part weighted [S]um: 3*FC + 1*F1 + 3*F2 + 1*F3 + 3*F4: %d", fs)
		verboseLineOut("    VS := [V]ariable part weighted [S]um: 3*V1 + 1*V2 + 3*V3 + 1*V4 + ... : %d", vs)
		verboseLineOut("    < MC := 29 - ( ( FS + VS ) modulo 29 ): %d", mc)
		verboseLineOut("  Concatenating FP & MC & VP:")
		verboseLineOut("    < %v & [%d] & %v", fp, mc, vp)
		verboseLineOut("  Encoding using custom Base32 mapping...")
		verboseLineOut("    Alphabet used: %s", customBase32Alphabet)
		verboseLineOut("< Result: %s", r)
	}

	return r
}

func Decode(x string) (r uint64, err error, complete bool) {
	defer func() {
		if err != nil {
			complete = false
		}
		if !complete {
			r = 0
		}
	}()

	id := make([]uint64, 0, 17)

	pos := 0
	check := 0

	for _, char := range x {
		pos++
		d, mapped := customBase32Decode(char)
		if !mapped {
			err = fmt.Errorf("Bad character in position %d: %c (%U)", pos, char, char)
			return
		}
		if pos <= 5 && d > 7 {
			err = fmt.Errorf("Non-[2,9]-numeric character in position %d: %c (%U)", pos, char, char)
			return
		}
		id = append(id, uint64(d))
		check += (1 + pos%2*2) * d
	}

	if pos >= 2 {
		r |= id[1] << 0
	}
	if pos >= 3 {
		r |= id[2] << 3
		if (bits.OnesCount64(r)+int((id[0]&0b100)>>2))%2 == 1 {
			err = fmt.Errorf("ID invalid starting at position 3")
			return
		}
	}
	if pos >= 4 {
		r |= id[3] << 6
		if (bits.OnesCount64(r)+int((id[0]&0b010)>>1))%2 == 1 {
			err = fmt.Errorf("ID invalid starting at position 4")
			return
		}
	}
	if pos >= 5 {
		r |= id[4] << 9
		if (bits.OnesCount64(r)+int((id[0]&0b001)>>0))%2 == 1 {
			err = fmt.Errorf("ID invalid starting at position 5")
			return
		}
	}

	for i := 6; i < len(id); i++ {
		r |= uint64(id[i] << (12 + (i-6)*5))
	}

	if pos < 6 {
		return
	}

	if check%29 != 0 {
		err = fmt.Errorf("ID invalid starting after position 5")
		return
	}

	complete = true
	return
}
