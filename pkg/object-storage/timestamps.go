package object_storage

import (
	"fmt"
	"log"
	"time"
)

type Bitstamp struct {
	Year   int
	Month  int
	Week   int
	Day    int
	Hour   int
	Minute int
	Second int
}

func NewBitstamp() *Bitstamp {
	t := time.Now()
	// ISOWeek returns the ISO 8601 year and week number in which t occurs.
	// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
	// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
	// of year n+1.
	year, week := t.ISOWeek()

	bs := Bitstamp{
		Year:   year - 2000,
		Month:  int(t.Month()),
		Week:   week,
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}
	return &bs
}

func CreateBitstamp(year int, month int, day int, hour int) *Bitstamp {
	return &Bitstamp{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
	}
}

func BitstampFromInt(bsi int) *Bitstamp {
	log.Println((2 ^ 8) - 1)
	year := (bsi & ((2 ^ 8) - 1) << 14) >> 14
	log.Println(year)
	month := bsi | (2^4-1)<<10
	day := bsi | (2^5-1)<<5
	hour := bsi | (2^5-1)<<5

	return &Bitstamp{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
	}
}

func (bs *Bitstamp) Int() int {
	return bs.Bits()
}

// As numeric, we get a single number representing the time.
// We only use year, week, day, and hour (for now)
// Year is in the range 0-258, so we will use 8 bits
// Month is 12, meaning 2^4
// Day is 31, meaning 2^5
// Hour is 24, meaning 2^5
// Total: 8 + 4 + 5 + 5 = 24 bits
// (If we add seconds, we need 2^6, still under a 32-bit number)
func (bs *Bitstamp) Bits() int {
	var bw int
	bsy := bs.Year << 14
	bsm := bs.Month << 10
	bsd := bs.Day << 5
	bsh := bs.Hour
	bw = bsy | bsm | bsd | bsh
	return bw
}

func (bs *Bitstamp) String() string {
	return fmt.Sprintf("%d-%02d-%02d-%02d",
		bs.Year, bs.Month, bs.Day, bs.Hour)
}

// 24-10-14-11
// 00011000 1010 01110 01011
//
//	24   10    14    11
func (bs *Bitstamp) BitString() string {
	return fmt.Sprintf("%08b%04b%05b%05b",
		bs.Year, bs.Month, bs.Day, bs.Hour)
}
