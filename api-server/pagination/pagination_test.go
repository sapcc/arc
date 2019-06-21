// +build integration

package pagination_test

import (
	. "github.com/sapcc/arc/api-server/pagination"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pagination", func() {

	Describe("create pagination", func() {

		var ()

		JustBeforeEach(func() {
		})

		It("should set default values", func() {
			// create url
			newUrl, err := url.Parse("http://www.someurl.com")
			Expect(err).NotTo(HaveOccurred())
			// create pagination
			pagination := CreatePagination(*newUrl)
			Expect(pagination.ActualPage).To(Equal(1))
			Expect(pagination.Offset).To(Equal(0))
			Expect(pagination.Limit).To(Equal(25))
			eq := reflect.DeepEqual(pagination.Request, *newUrl)
			Expect(eq).To(Equal(true))
		})

		It("should reset wrong values", func() {
			// create url
			newUrl, err := url.Parse("http://www.someurl.com?page=-1&per_page=2000")
			Expect(err).NotTo(HaveOccurred())
			// create pagination
			pagination := CreatePagination(*newUrl)
			Expect(pagination.ActualPage).To(Equal(1))
			Expect(pagination.Offset).To(Equal(0))
			Expect(pagination.Limit).To(Equal(100))
			eq := reflect.DeepEqual(pagination.Request, *newUrl)
			Expect(eq).To(Equal(true))
		})

		It("should set values", func() {
			// create url
			newUrl, err := url.Parse("http://www.someurl.com?page=2&per_page=50")
			Expect(err).NotTo(HaveOccurred())
			// create pagination
			pagination := CreatePagination(*newUrl)
			Expect(pagination.ActualPage).To(Equal(2))
			Expect(pagination.Offset).To(Equal(50))
			Expect(pagination.Limit).To(Equal(50))
			eq := reflect.DeepEqual(pagination.Request, *newUrl)
			Expect(eq).To(Equal(true))
		})

	})

	Describe("setting total of Elements", func() {

		It("should set the total elements and total pages attribute", func() {
			// create url
			newUrl, err := url.Parse("http://www.someurl.com?page=2&per_page=25")
			Expect(err).NotTo(HaveOccurred())
			// create pagination
			pagination := CreatePagination(*newUrl)
			pagination.SetTotalElements(130)
			Expect(pagination.TotalElements).To(Equal(130))
			Expect(pagination.TotalPages).To(Equal(6))
		})

		It("should correct actual page and offset", func() {
			// create url
			newUrl, err := url.Parse("http://www.someurl.com?page=8&per_page=25")
			Expect(err).NotTo(HaveOccurred())
			// create pagination
			pagination := CreatePagination(*newUrl)
			Expect(pagination.ActualPage).To(Equal(8))
			Expect(pagination.Offset).To(Equal(175))
			// set total elements and correct actual page and offset
			pagination.SetTotalElements(130)
			Expect(pagination.ActualPage).To(Equal(6))
			Expect(pagination.Offset).To(Equal(125))
		})

		Describe("links", func() {

			It("should set self link when there is just one page", func() {
				// create url
				newUrl, err := url.Parse("http://www.someurl.com/relative_path?page=3&per_page=25")
				Expect(err).NotTo(HaveOccurred())
				// create pagination
				pagination := CreatePagination(*newUrl)
				// set total elements and correct actual page and offset
				pagination.SetTotalElements(1)
				Expect(pagination.GetLinks()).To(Equal(fmt.Sprintf(`<%s>;rel="self"`, "/relative_path?page=1&per_page=25")))
			})

			It("should set self, next and last link when first page and there is more than one page", func() {
				// create url
				newUrl, err := url.Parse("http://www.someurl.com/relative_path?page=1&per_page=25")
				Expect(err).NotTo(HaveOccurred())
				// create pagination
				pagination := CreatePagination(*newUrl)
				// set total elements and correct actual page and offset
				pagination.SetTotalElements(30)
				Expect(pagination.GetLinks()).To(Equal(fmt.Sprintf(`<%s>;rel="self",<%s>;rel="next",<%s>;rel="last"`, "/relative_path?page=1&per_page=25", "/relative_path?page=2&per_page=25", "/relative_path?page=2&per_page=25")))
			})

			It("should set self, first, prev, next and last link when between first and last page and there is more than two pages", func() {
				// create url
				newUrl, err := url.Parse("http://www.someurl.com/relative_path?page=2&per_page=25")
				Expect(err).NotTo(HaveOccurred())
				// create pagination
				pagination := CreatePagination(*newUrl)
				// set total elements and correct actual page and offset
				pagination.SetTotalElements(100)
				Expect(pagination.GetLinks()).To(Equal(fmt.Sprintf(`<%s>;rel="self",<%s>;rel="first",<%s>;rel="prev",<%s>;rel="next",<%s>;rel="last"`, "/relative_path?page=2&per_page=25", "/relative_path?page=1&per_page=25", "/relative_path?page=1&per_page=25", "/relative_path?page=3&per_page=25", "/relative_path?page=4&per_page=25")))
			})

			It("should set self, first and prev when last page and there is more than one page", func() {
				// create url
				newUrl, err := url.Parse("http://www.someurl.com/relative_path?page=4&per_page=25")
				Expect(err).NotTo(HaveOccurred())
				// create pagination
				pagination := CreatePagination(*newUrl)
				// set total elements and correct actual page and offset
				pagination.SetTotalElements(100)
				Expect(pagination.GetLinks()).To(Equal(fmt.Sprintf(`<%s>;rel="self",<%s>;rel="first",<%s>;rel="prev"`, "/relative_path?page=4&per_page=25", "/relative_path?page=1&per_page=25", "/relative_path?page=3&per_page=25")))
			})

		})

		Describe("adding headers", func() {

			It("should set pages, elements, perpage and links headers", func() {
				pages := ""
				elements := ""
				per_page := ""
				link := ""
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					pag := CreatePagination(*r.URL)
					pag.SetHeaders(w)
					pages = w.Header().Get("Pagination-Pages")
					elements = w.Header().Get("Pagination-Elements")
					per_page = w.Header().Get("Pagination-Per-Page")
					link = w.Header().Get("Link")
				}))
				defer server.Close()
				// send request
				http.Get(server.URL)

				Expect(pages).To(Equal("0"))
				Expect(elements).To(Equal("0"))
				Expect(per_page).To(Equal("25")) // default value
				Expect(link).To(Equal(`<>;rel="self",<>;rel="first",<>;rel="prev",<>;rel="next",<>;rel="last"`))
			})

		})

	})

})
