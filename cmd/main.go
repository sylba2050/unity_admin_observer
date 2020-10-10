package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	u "github.com/sylba2050/unity_admin_observer"
)

func getDriver() *agouti.WebDriver {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
				"no-sandbox",
			}),
	)
	err := driver.Start()
	if err != nil {
		panic(err)
	}

	return driver
}

func login(page *agouti.Page) {
	err := page.Navigate(u.Config.LoginURL)
	if err != nil {
		panic(err)
	}
	time.Sleep(2000 * time.Millisecond)
	page.FindByID("conversations_create_session_form_email").Fill(u.Config.Mail)
	page.FindByID("conversations_create_session_form_password").Fill(u.Config.Password)
	page.FindByName("commit").Click()
	time.Sleep(5000 * time.Millisecond)
}

func getSalesPageData(page *agouti.Page) {
	err := page.Navigate(u.Config.SalesURL)
	if err != nil {
		panic(err)
	}
	time.Sleep(3000 * time.Millisecond)

	html, err := page.HTML()
	if err != nil {
		panic(err)
	}
	pageReader := strings.NewReader(html)
	dom, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		panic(err)
	}
	dom.Find("table#sales>thead>tr>td").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
	dom.Find("table#sales>tbody>tr>td").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}

func main() {
	driver := getDriver()
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		panic(err)
	}

	login(page)
	getSalesPageData(page)
}
