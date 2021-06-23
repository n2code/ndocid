# ndocid
> *human-writeable & human-typeable identifiers*
```console
$ ndocid -n #Thu, 10 Oct 2019 01:41:40 +0200
96822L9IPD
```
This command-line tool generates an alphanumeric identifier for documents, records or anything, really, from a given point in time, number or a short bitstream (technically speaking: any 64 bit unsigned integer).

The input is encoded in an identifier which has many useful practical properties.
It is...
* **unique & reversible** (because it is an encoding and _not_ a hash)
* **robust**:
  * only numbers and upper-case latin characters
  * alphabet of 32 characters without `S/G/1/0` to avoid confusion with `5/6/I/O` in written form
  * no leading zeros to prevent annoyance when working with spreadsheets, input forms etc.
* **spread evenly**, i.e. when grouping generated identifiers according to their first 1/2/3/... characters you end up with buckets of roughly the same size
* **input-friendly**:
  * contains a checksum to detect typos
  * early-match:tm: implementation possible for auto-completed and database-backed input forms
    * ID starts with five digits so input with numpad is possibly sufficient (especially if the search space is small enough)
    * checksum for partial identifier
* as **short** as possible while retaining all features above
  * 0  to 12 bit => 6 character ID (5 digits + 1 letter/digit)
  * up to 17 bit => 7 character ID (5 digits + 2 letters/digits)
  * up to 22 bit => 8 character ID (5 digits + 3 letters/digits)
  * ...

## Examples
```console
$ ndocid -i 123456789 #plain integer input
674685WFX
```
```console
$ ndocid -d 20191231223000 #new year's eve
52598NV7RD
```
```console
$ ndocid -b "11010011 01000111" #16 bit
89273IF
```

## Usage
**`ndocid`** `[-v]` `[MODE] INPUT`
```console
$ ndocid -h
Usage of ./ndocid:
  -b "00010110 11011011"
    	BITSTRING-MODE: Generate ID from string of bits, e.g. "00010110 11011011".
    	  Spaces, tabs, underscores and leading zeros are being dropped.
    	  The maximum length is 64 bits.
    	  Bad input will result in an exit code greater than 0.
  -d 20060102150405
    	DATE-MODE: Generate ID from given date and time.
    	  For example 20060102150405 which represents "Mon Jan 2 15:04:05 2006".
    	  Evaluated in the machine's time zone.
    	  Exit code greater than 0 if the input is not according to format.
  -i 42
    	INTEGER-MODE: Generate ID from number, e.g. 42.
    	  Accepts any positive decimal number that can fit in an unsigned 64 bit integer.
    	  Exit code greater than 0 if input exceeds range.
  -n	NOW-MODE: Generate ID from current date and time of this machine.
  -r 72639D77LD
    	REVERSING/CHECK-MODE: Validates given ID, e.g. 72639D77LD.
    	  Exit code 0: Valid full ID
    	  Exit code 1: Invalid ID
    	  Exit code 4: Plausible partial ID (beginning), needs further digits
    	  The first line returned is OK / ERROR / PARTIAL for exit codes 0 / 1 / 4.
  -v	Verbose option: Generate more human-readable output.
    	  Explains algorithm in MODEs that generate IDs.
    	  Provides possible source representations when reversing is successful.
```
