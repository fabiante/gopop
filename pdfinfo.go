package gopop

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type PageSize struct {
	Width  float64
	Height float64
	Unit   string
	Note   string
}

// PDFInfoPage contains information about a single page.
//
//	Page    1 size:  595.32 x 841.92 pts (A4)
//	Page    1 rot:   0
//	Page    2 size:  595.32 x 841.92 pts (A4)
//	Page    2 rot:   0
type PDFInfoPage map[string]string

var regexPDFInfoPageSize = regexp.MustCompile("^(.*) x (.*) (pts) ?(.*)$")

// Size parses the value of SizeRaw and returns a PageSize.
//
// Panics if parsing fails
func (page PDFInfoPage) Size() PageSize {
	raw := page.SizeRaw()
	find := regexPDFInfoPageSize.FindStringSubmatch(raw)
	if find == nil {
		panic(fmt.Errorf("parsing page size %q failed", raw))
	}
	w := find[1]
	h := find[2]
	unit := find[3]
	note := ""
	if len(find) > 4 {
		note = find[4]
	}

	wf, err := strconv.ParseFloat(w, 64)
	if err != nil {
		panic(fmt.Errorf("parsing %q to float failed", w))
	}

	hf, err := strconv.ParseFloat(h, 64)
	if err != nil {
		panic(fmt.Errorf("parsing %q to float failed", h))
	}

	return PageSize{
		Width:  wf,
		Height: hf,
		Unit:   unit,
		Note:   note,
	}
}

func (page PDFInfoPage) SizeRaw() string {
	return page["size"]
}

func (page PDFInfoPage) RotRaw() string {
	return page["rot"]
}

type PDFInfo struct {
	// Properties holds arbitrary key/value pairs that are not further interpreted.
	// An example:
	//
	//     Producer:        PDFKit.NET 23.1.101.39834 DMV9
	//     CreationDate:    Wed Apr 19 12:58:37 2023 CEST
	//     ModDate:         Wed Apr 19 12:58:37 2023 CEST
	//     Custom Metadata: no
	//     Metadata Stream: yes
	//     Tagged:          no
	//     UserProperties:  no
	//     Suspects:        no
	//     Form:            AcroForm
	//     JavaScript:      no
	//     Pages:           2
	//     Encrypted:       no
	//     File size:       230296 bytes
	//     Optimized:       no
	//     PDF version:     1.7
	Properties map[string]string

	Pages map[string]PDFInfoPage
}

func NewPDFInfo() *PDFInfo {
	return &PDFInfo{
		Properties: map[string]string{},
		Pages:      map[string]PDFInfoPage{},
	}
}

func (info *PDFInfo) Producer() string {
	return info.Properties["Producer"]
}

// editPage is a helper method to edit a page. This method will create a page if
// it does not exist.
func (info *PDFInfo) editPage(page string, edit func(page PDFInfoPage)) {
	p := info.Pages[page]
	if p == nil {
		p = PDFInfoPage{}
		info.Pages[page] = p
	}
	edit(p)
}

// Page returns the given page or nil if it does not exist.
func (info *PDFInfo) Page(page int) PDFInfoPage {
	return info.Pages[strconv.Itoa(page)]
}

// PDFInfoParser is a parser for the output of the pdfinfo command.
//
// This parser is not fully tested and may break if pdfinfo outputs something that we don't expect.
type PDFInfoParser struct {
	reader io.Reader
	info   *PDFInfo
}

var (
	regexParseLine    = regexp.MustCompile("^(.*?):[\t ]*(.*)$")
	regexParsePageKey = regexp.MustCompile("^Page[\\t ]*?(\\d+)[\\t ](.*)$")
)

func NewPDFInfoParser(reader io.Reader) *PDFInfoParser {
	return &PDFInfoParser{reader: reader, info: NewPDFInfo()}
}

func (parser *PDFInfoParser) Parse() (*PDFInfo, error) {
	scanner := bufio.NewScanner(parser.reader)

	for scanner.Scan() {
		line := scanner.Text()
		if err := parser.parseLine(line); err != nil {
			return nil, err
		}
	}

	return parser.info, nil
}

func (parser *PDFInfoParser) parseLine(line string) error {
	match := regexParseLine.FindStringSubmatch(line)
	if match == nil {
		return nil
	}

	key := match[1]
	value := match[2]

	switch {
	case strings.HasPrefix(key, "Page "):
		if err := parser.parsePageKey(key, value); err != nil {
			return err
		}
		break
	default:
		parser.info.Properties[key] = value
		break
	}

	return nil
}

func (parser *PDFInfoParser) parsePageKey(key string, value string) error {
	keyMatch := regexParsePageKey.FindStringSubmatch(key)
	if keyMatch == nil {
		return fmt.Errorf("invalid page key %q", key)
	}

	pageNumber := keyMatch[1]
	pageKey := keyMatch[2]

	parser.info.editPage(pageNumber, func(page PDFInfoPage) {
		page[pageKey] = value
	})

	return nil
}
