package main

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	timeLocation         = "Asia/Tokyo"
	selectorHeader       = "body > table > tbody > tr > td > div:nth-child(4) > table:nth-child(1) > tbody > tr > th"
	selectorClasses      = "body > table > tbody > tr > td > div:nth-child(4) > table:nth-child(3) > tbody > .style1"
	selectorLastModified = "body > table > tbody > tr > td > div:nth-child(7) > center > table > tbody > tr:nth-child(1) > th"
	regexpTrimGarbage    = "( \n| {2,})"
)

var (
	reTrimGarbage *regexp.Regexp
)

func init() {
	l, err := time.LoadLocation(timeLocation)
	if err != nil {
		l = time.UTC
	}
	time.Local = l

	reTrimGarbage = regexp.MustCompile(regexpTrimGarbage)
	reTrimGarbage.Longest()
}

func parseClasses(root *goquery.Document) (*[]Class, error) {
	var classes []Class

	root.Find(selectorClasses).Each(func(i int, s *goquery.Selection) {
		var cells []string

		s.Children().Each(func(i int, s *goquery.Selection) {
			cells = append(cells, s.Text())
		})
		if len(cells) != 4 {
			return
		}

		period, _ := strconv.Atoi(strings.Replace(cells[0], "講時", "", 1))
		if period == 0 {
			period = classes[len(classes)-1].Period
		}

		instructor := reTrimGarbage.ReplaceAllString(cells[2], " ")
		reason := reTrimGarbage.ReplaceAllString(cells[3], "")

		class := Class{
			Period:     period,
			Name:       cells[1],
			Instructor: instructor,
			Reason:     reason,
		}

		classes = append(classes, class)
	})

	return &classes, nil
}

func parseLastModified(root *goquery.Document) (time.Time, error) {
	m := root.Find(selectorLastModified).Text()
	m = strings.NewReplacer("更新日時：", "", " ", "", "\n", "").Replace(m)

	t, err := time.ParseInLocation("2006年1月2日15時04分", m, time.Local)
	if err != nil {
		return t, errors.Wrap(err, "can't parse updated date")
	}

	return t, nil
}

func parseHeaderDate(root *goquery.Document) (time.Time, error) {
	h := root.Find(selectorHeader).Text()

	h = strings.NewReplacer(" ", "", " ", "", "\n", "").Replace(h)
	h = h[strings.Index(h, "]")+1 : strings.Index(h, "(")]

	t, err := time.ParseInLocation("2006年1月2日", h, time.Local)
	if err != nil {
		return t, errors.Wrap(err, "can't parse date in header")
	}

	return t, nil
}

func parseHeaderLocation(root *goquery.Document) (string, error) {
	h := root.Find(selectorHeader).Text()
	h = h[strings.Index(h, "[")+1 : strings.Index(h, "]")]

	return h, nil
}

func Parse(f io.Reader) (*Notification, error) {
	n := &Notification{}

	/*
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	*/

	root, err := goquery.NewDocumentFromReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return n, err
	}

	location, err := parseHeaderLocation(root)
	if err != nil {
		return n, err
	}

	classes, err := parseClasses(root)
	if err != nil {
		return n, err
	}

	date, err := parseHeaderDate(root)
	if err != nil {
		return n, err
	}

	updatedAt, err := parseLastModified(root)
	if err != nil {
		return n, err
	}

	n = &Notification{location, *classes, date, updatedAt}

	return n, nil
}
