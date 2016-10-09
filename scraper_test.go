package main

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	testDate     time.Time
	testTime     time.Time
	testResult   *Notification
	testDoc      *goquery.Document
	testDocEmpty *goquery.Document
)

func TestMain(m *testing.M) {
	l, e := time.LoadLocation("Asia/Tokyo")
	if e != nil {
		panic(e)
	}

	testDate, e = time.ParseInLocation("2006年1月2日15時04分", "2011年11月11日00時00分", l)
	if e != nil {
		panic(e)
	}

	testTime, e = time.ParseInLocation("2006年1月2日15時04分", "2011年11月11日11時11分", l)
	if e != nil {
		panic(e)
	}

	testResult = &Notification{
		Location: "location",
		Classes: []Class{
			{1, "class A", "A A", "reason A"},
			{2, "class B", "B B", "reason B"},
			{2, "class C", "C C C", "reason C"},
			{2, "class D", "D D", "reason D"},
			{3, "class E", "E E", "reason E"},
			{3, "class F", "F F F", "reason F"},
			{3, "class G", "G G", "reason G"},
			{3, "class H", "H H H", "reason H"},
			{4, "class I", "I I", "reason I"},
		},
		Date:      testDate,
		UpdatedAt: testTime,
	}

	testDoc = loadDoc("./testdata/page.html")
	testDocEmpty = loadDoc("./testdata/page_empty.html")

	m.Run()
}

func loadDoc(path string) *goquery.Document {
	f, e := os.Open(path)
	if e != nil {
		panic(e)
	}
	defer f.Close()

	d, e := goquery.NewDocumentFromReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	if e != nil {
		panic(e)
	}

	return d
}

func TestParseClasses(t *testing.T) {
	c, e := parseClasses(testDoc)
	if e != nil {
		panic(e)
	}

	if !reflect.DeepEqual(testResult.Classes, *c) {
		t.Errorf("parseClasses() mismatch:\ngot  %+v\nwant %+v", c, testResult.Classes)
	}
}

func TestParseClassesEmpty(t *testing.T) {
	c, e := parseClasses(testDocEmpty)
	if e != nil {
		panic(e)
	}

	if len(*c) != 0 {
		t.Errorf("parseClasses() mismatch:\ngot  %+v\nwant %+v", c, &[]Class{})
	}
}

func TestParseLastModified(t *testing.T) {
	tt, e := parseLastModified(testDoc)
	if e != nil {
		panic(e)
	}

	if !reflect.DeepEqual(testResult.UpdatedAt, tt) {
		t.Errorf("parseLastModified() mismatch:\ngot  %+v\nwant %+v", tt, testResult.UpdatedAt)
	}
}

func TestParseHeaderDate(t *testing.T) {
	tt, e := parseHeaderDate(testDoc)
	if e != nil {
		panic(e)
	}

	if !reflect.DeepEqual(testResult.Date, tt) {
		t.Errorf("parseHeaderDate() mismatch:\ngot  %+v\nwant %+v", tt, testResult.Date)
	}
}

func TestParseHeaderLocation(t *testing.T) {
	l, e := parseHeaderLocation(testDoc)
	if e != nil {
		panic(e)
	}

	if testResult.Location != l {
		t.Errorf("parseHeaderLocation() mismatch:\ngot  %+v\nwant %+v", l, testResult.Location)
	}
}

func TestParse(t *testing.T) {
	f, e := os.Open("./testdata/page.html")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	n, e := Parse(f)
	if e != nil {
		panic(e)
	}

	if !reflect.DeepEqual(testResult, n) {
		t.Errorf("Parse() mismatch:\ngot  %+v\nwant %+v", n, testResult)
	}
}
