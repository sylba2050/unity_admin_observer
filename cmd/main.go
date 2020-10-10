package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	u "github.com/Sylba2050/unity_admin_observer"
	"github.com/sclevine/agouti"
)

func main() {
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
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		panic(err)
	}
	err = page.Navigate(u.Config.LoginURL)
	if err != nil {
		panic(err)
	}
	time.Sleep(2000 * time.Millisecond)
	page.FindByID("conversations_create_session_form_email").Fill(u.Config.Mail)
	page.FindByID("conversations_create_session_form_password").Fill(u.Config.Password)
	page.FindByName("commit").Click()
	time.Sleep(5000 * time.Millisecond)

	err = page.Navigate(u.Config.SalesURL)
	if err != nil {
		panic(err)
	}
	time.Sleep(1000 * time.Millisecond)

	html, err := page.HTML()
	if err != nil {
		panic(err)
	}
	pageReader := strings.NewReader(html)
	dom, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		panic(err)
	}
	dom.Find("table#sales>tbody>tr>td").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}
