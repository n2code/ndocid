package ndocid

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

var coverVerboseButHideIt = func() {
	Verbose = true
	verboseLineOut = func(format string, msg ...interface{}) {}
}

func TestCustomBase32(t *testing.T) {
	for x := 0; x < 32; x++ {
		e := customBase32Encode(x)
		d, ok := customBase32Decode(e)
		if !ok {
			t.Fatalf("%c not decodable", e)
		}
		if d != x {
			t.Fatalf("decoding of %c which is %d does not match encoding of %d", e, d, x)
		}
	}
}

func TestEncodeUint64(t *testing.T) {
	assertEncoded := func(input uint64, exp string) {
		act := EncodeUint64(input)
		if act != exp {
			t.Errorf(`%d encoded, got "%s" but expected "%s"`, input, act, exp)
		}
	}

	coverVerboseButHideIt()

	//the following encodings have been calculated manually
	assertEncoded(1552572000, "72639D77LD")                   //Pi-Day 2019, 3pm in Germany
	assertEncoded(1567856598, "68495LTTOD")                   //07.09.2019 13:43:18 in Germany, first ndocid source file created
	assertEncoded(4133980800, "52247CRMTY")                   //01.01.2101 00:00:00 UTC, new millennium
	assertEncoded(1234567890, "3445352D8B")                   //Ascending numbers
	assertEncoded(0, "22222X")                                //Minimum
	assertEncoded(0xFFFF_FFFF_FFFF_FFFF, "499997ZZZZZZZZZZ5") //Maximum
}

func TestEncodeDatetime(t *testing.T) {
	// tests assume Germany's time zone
	var err error
	baseLocation, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Skip(`cannot test datetime conversion: `, err)
	}

	coverVerboseButHideIt()

	assertEncoded := func(input string, exp string) {
		act, err := EncodeDatetime(input)
		if err != nil {
			t.Errorf(`unexpected error on encoding "%s": %s`, input, err)
		}
		if act != exp {
			t.Errorf(`"%s" encoded, got "%s" but expected "%s"`, input, act, exp)
		}
	}
	assertEncodingFailure := func(input string) {
		_, err := EncodeDatetime(input)
		if err == nil {
			t.Errorf(`no error on attempt of encoding "%s"`, input)
		}
	}

	Verbose = true //to cover verbose output

	//the following encodings have been calculated manually
	assertEncoded("20190314150000", "72639D77LD") //Pi-Day 2019, 3pm in Germany
	assertEncoded("20190907134318", "68495LTTOD") //07.09.2019 13:43:18 in Germany, first ndocid source file created
	assertEncoded("21010101010000", "52247CRMTY") //01.01.2101 00:00:00 UTC, new millennium
	assertEncoded("20190915215634", "94875JBZOD")

	assertEncodingFailure("42")
	assertEncodingFailure("20001231123000.0000")
	assertEncodingFailure("20000230000000")

	baseLocation = time.Local
}

func TestEncodeBitstring(t *testing.T) {
	assertEncoded := func(input string, exp string) {
		act, err := EncodeBitstring(input)
		if err != nil {
			t.Errorf(`unexpected error on encoding "%s": %s`, input, err)
		}
		if act != exp {
			t.Errorf(`"%s" encoded, got "%s" but expected "%s"`, input, act, exp)
		}
	}
	assertEncodingFailure := func(input string) {
		_, err := EncodeBitstring(input)
		if err == nil {
			t.Errorf(`no error on attempt of encoding "%s"`, input)
		}
	}

	coverVerboseButHideIt()

	assertEncoded("0", EncodeUint64(0))
	assertEncoded("000", EncodeUint64(0))
	assertEncoded("00 \t 00 01 \t 10 01", EncodeUint64(25))
	assertEncoded("1010_1111", EncodeUint64(175))
	assertEncoded("01011101 01110011 10010111 11010110", "68495LTTOD")
	assertEncoded(fmt.Sprintf("00%s", strings.Repeat("1", 64)), EncodeUint64(0xFFFF_FFFF_FFFF_FFFF))

	assertEncodingFailure("")
	assertEncodingFailure("0x01")
	assertEncodingFailure(fmt.Sprintf("%s0", strings.Repeat("1", 64)))
	assertEncodingFailure(fmt.Sprintf("%s1", strings.Repeat("1", 64)))
}

func TestDecodeSuccess(t *testing.T) {
	assertDecoded := func(input string, exp uint64) {
		act, err, complete := Decode(input)
		if err != nil {
			t.Errorf(`unexpected error on decoding "%s": %s`, input, err)
		}
		if act != exp {
			t.Errorf(`"%s" decoded, got %d but expected %d`, input, act, exp)
		}
		if !complete {
			t.Errorf(`"%s" unexpectedly incomplete`, input)
		}
	}

	//proper values
	assertDecoded("72639D77LD", 1552572000)                   //Pi-Day 2019, 3pm in Germany
	assertDecoded("68495LTTOD", 1567856598)                   //07.09.2019 13:43:18 in Germany, first ndocid source file created
	assertDecoded("52247CRMTY", 4133980800)                   //01.01.2101 00:00:00 UTC, new millennium
	assertDecoded("3445352D8B", 1234567890)                   //Ascending numbers
	assertDecoded("22222X", 0)                                //Minimum
	assertDecoded("499997ZZZZZZZZZZ5", 0xFFFF_FFFF_FFFF_FFFF) //Maximum

}

func TestDecodePartialSuccess(t *testing.T) {
	assertIncompleteButValid := func(input string) {
		act, err, complete := Decode(input)
		if err != nil {
			t.Errorf(`unexpected error on decoding partial "%s": %s`, input, err)
		}
		if complete {
			t.Errorf(`"%s" decoding unexpectedly complete`, input)
		}
		if act != 0 {
			t.Errorf(`partial "%s" decoding unexpectedly returned %d`, input, act)
		}
	}

	//partial but successful
	assertIncompleteButValid("")
	assertIncompleteButValid("6")
	assertIncompleteButValid("68")
	assertIncompleteButValid("684")
	assertIncompleteButValid("6849")
	assertIncompleteButValid("68495")
	assertIncompleteButValid("5")
	assertIncompleteButValid("52")
	assertIncompleteButValid("522")
	assertIncompleteButValid("5224")
	assertIncompleteButValid("52247")
}

func TestDecodeFailure(t *testing.T) {
	assertDecodingFailure := func(input string) {
		_, err, complete := Decode(input)
		if err == nil {
			t.Errorf(`no error on attempt of decoding "%s"`, input)
		}
		if complete {
			t.Errorf(`"%s" unexpectedly complete`, input)
		}
	}

	//bad chars
	assertDecodingFailure("2222!X")
	assertDecodingFailure("22220X")
	assertDecodingFailure("22221X")
	assertDecodingFailure("2222TX")

	//partial failures
	assertDecodingFailure("22332")
	assertDecodingFailure("22332")
	assertDecodingFailure("22233")
	assertDecodingFailure("22223")
	assertDecodingFailure("92222")
	assertDecodingFailure("92332")
	assertDecodingFailure("92323")

	//bad checksums

	//single digit change: 94875JBZOD
	assertDecodingFailure("94875KBZOD")
	assertDecodingFailure("94875J5ZOD")
	assertDecodingFailure("94875JBSOD")
	assertDecodingFailure("94875JBZQD")
	assertDecodingFailure("94875JBZO4")

	//two digit swap:      94875JBZOD
	assertDecodingFailure("49875JBZOD")
	assertDecodingFailure("98475JBZOD")
	assertDecodingFailure("94785JBZOD")
	assertDecodingFailure("94857JBZOD")
	assertDecodingFailure("94875BJZOD")
	assertDecodingFailure("94875JZBOD")
	assertDecodingFailure("94875JBOZD")
	assertDecodingFailure("94875JBZDO")
}

func TestAlphabetConsistency(t *testing.T) {
	if len(customBase32Alphabet) != 32 || utf8.RuneCountInString(customBase32Alphabet) != 32 {
		t.Fatal("alphabet broken, does not contain 32 single-byte characters")
	}
}
