package pagination

import (
	"fmt"
	"net/url"
	"strconv"
)

type Pagination struct {
	Offset        int
	Limit         int
	ActualPage    int
	TotalPages    int
	TotalElements int
	LinkSelf      string
	LinkFirst     string
	LinkPrev      string
	LinkNext      string
	LinkLast      string
	Request       url.URL
}

func CreatePagination(reqUrl url.URL) *Pagination {
	// get values and convert string to int
	intPage, _ := strconv.Atoi(reqUrl.Query().Get("page"))
	intPerPage, _ := strconv.Atoi(reqUrl.Query().Get("per_page"))

	// check limit
	limit := 25
	if intPerPage < 1 {
		limit = 25
	} else if intPerPage > 100 {
		limit = 100
	}
	// check page
	if intPage < 1 {
		intPage = 1
	}

	// create pag element
	pag := Pagination{
		Offset:     calcOffset(intPage, limit),
		Limit:      limit,
		ActualPage: intPage,
		Request:    reqUrl,
	}

	return &pag
}

func (pag *Pagination) SetTotalElements(rows int) error {
	// total elements and pages
	pag.TotalElements = rows
	err := pag.setTotalPages(rows)
	if err != nil {
		return err
	}

	// check page again
	if pag.ActualPage > pag.TotalPages {
		pag.ActualPage = pag.TotalPages
		pag.Offset = calcOffset(pag.ActualPage, pag.Limit)
	}

	// set links
	err = pag.setLinkSelf(pag.Request)
	if err != nil {
		return err
	}
	err = pag.setLinkFirst()
	if err != nil {
		return err
	}
	err = pag.setLinkPrev()
	if err != nil {
		return err
	}
	err = pag.setLinkNext()
	if err != nil {
		return err
	}
	err = pag.setLinkLast()
	if err != nil {
		return err
	}
	return nil
}

func (pag *Pagination) GetLinks() string {
	links := ""
	if pag.ActualPage == 1 && pag.TotalPages == 1 {
		links = fmt.Sprintf(`<%s>;rel="self"`, pag.LinkSelf)
	} else if pag.ActualPage == 1 && pag.TotalPages > 1 {
		links = fmt.Sprintf(`<%s>;rel="self",<%s>;rel="next",<%s>;rel="last"`, pag.LinkSelf, pag.LinkNext, pag.LinkLast)
	} else if pag.ActualPage == pag.TotalPages && pag.TotalPages > 1 {
		links = fmt.Sprintf(`<%s>;rel="self",<%s>;rel="first",<%s>;rel="prev"`, pag.LinkSelf, pag.LinkFirst, pag.LinkPrev)
	} else {
		links = fmt.Sprintf(`<%s>;rel="self",<%s>;rel="first",<%s>;rel="prev",<%s>;rel="next",<%s>;rel="last"`, pag.LinkSelf, pag.LinkFirst, pag.LinkPrev, pag.LinkNext, pag.LinkLast)
	}
	return links
}

// private methods

func calcOffset(page, limit int) int {
	offset := 0
	if page < 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	return offset
}

func (pag *Pagination) setLinkSelf(reqUrl url.URL) error {
	// copy and change url attributes
	values := reqUrl.Query()
	values["page"] = []string{fmt.Sprintf("%v", pag.ActualPage)}
	values["per_page"] = []string{fmt.Sprintf("%v", pag.Limit)}

	// create url
	newUrl, err := url.Parse(reqUrl.RequestURI())
	if err != nil {
		return err
	}

	// add values
	newUrl.RawQuery = values.Encode()

	// save url
	pag.LinkSelf = newUrl.RequestURI()

	return nil
}

func (pag *Pagination) setLinkFirst() error {
	// create url
	newUrl, err := url.Parse(pag.LinkSelf)
	if err != nil {
		return err
	}

	// copy and change url attributes
	values := newUrl.Query()
	values["page"] = []string{fmt.Sprintf("%v", 1)}
	values["per_page"] = []string{fmt.Sprintf("%v", pag.Limit)}

	// add values
	newUrl.RawQuery = values.Encode()

	// save url
	pag.LinkFirst = newUrl.RequestURI()

	return nil
}

func (pag *Pagination) setLinkPrev() error {
	// create url
	newUrl, err := url.Parse(pag.LinkSelf)
	if err != nil {
		return err
	}

	prevPage := 1
	if pag.ActualPage <= 1 {
		prevPage = 1
	} else if pag.ActualPage > pag.TotalPages {
		prevPage = pag.TotalPages
	} else {
		prevPage = pag.ActualPage - 1
	}

	// copy and change url attributes
	values := newUrl.Query()
	values["page"] = []string{fmt.Sprintf("%v", prevPage)}
	values["per_page"] = []string{fmt.Sprintf("%v", pag.Limit)}

	// add values
	newUrl.RawQuery = values.Encode()

	// save url
	pag.LinkPrev = newUrl.RequestURI()

	return nil
}

func (pag *Pagination) setLinkNext() error {
	// create url
	newUrl, err := url.Parse(pag.LinkSelf)
	if err != nil {
		return err
	}

	nextPage := 1
	if pag.ActualPage >= pag.TotalPages {
		nextPage = pag.TotalPages
	} else if pag.ActualPage < 1 {
		nextPage = 1
	} else {
		nextPage = pag.ActualPage + 1
	}

	// copy and change url attributes
	values := newUrl.Query()
	values["page"] = []string{fmt.Sprintf("%v", nextPage)}
	values["per_page"] = []string{fmt.Sprintf("%v", pag.Limit)}

	// add values
	newUrl.RawQuery = values.Encode()

	// save url
	pag.LinkNext = newUrl.RequestURI()

	return nil
}

func (pag *Pagination) setLinkLast() error {
	// create url
	newUrl, err := url.Parse(pag.LinkSelf)
	if err != nil {
		return err
	}

	// copy and change url attributes
	values := newUrl.Query()
	values["page"] = []string{fmt.Sprintf("%v", pag.TotalPages)}
	values["per_page"] = []string{fmt.Sprintf("%v", pag.Limit)}

	// add values
	newUrl.RawQuery = values.Encode()

	// save url
	pag.LinkLast = newUrl.RequestURI()

	return nil
}

func (pag *Pagination) setTotalPages(rows int) error {
	// get total pages
	totalPages := rows / pag.Limit
	if rows%pag.Limit > 0 {
		totalPages = totalPages + 1
	}
	pag.TotalPages = totalPages

	return nil
}
