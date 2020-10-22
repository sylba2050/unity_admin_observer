package main

import (
	"os"
	"bufio"
	"strconv"
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
	err = page.FindByName("commit").Click()
	if err != nil{
		panic(err)
	}
	time.Sleep(1000 * time.Millisecond)

	fp, err := os.OpenFile("/home/siruba_2050/unity_admin_observer/log.txt", os.O_WRONLY | os.O_APPEND, 0644)
        if err != nil {
	    panic(err)
        }
	defer fp.Close()
	url, err := page.URL()
        if err != nil {
	    panic(err)
        }
	writer := bufio.NewWriter(fp)
	_, err = writer.WriteString(url+"¥n")
        if err != nil {
	    panic(err)
        }
	writer.Flush()
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// page.FindByID("conversations_email_tfa_required_form_code").Fill(scanner.Text())
	// err = page.FindByName("commit").Click()
	// if err != nil{
	//	panic(err)
	// }
	time.Sleep(3000 * time.Millisecond)
}

func getSalesPageData(page *agouti.Page) ([]string, []int) {
	var packages []string
	var title []string
	var nowSales []int

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
		title = append(title, s.Text())
	})
	dom.Find("table#sales>tbody>tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(j int, ss *goquery.Selection) {
			if title[j] == "Package" {
				packages = append(packages, ss.Text())
			}
			if title[j] == "Qty" {
				atoi, err := strconv.Atoi(ss.Text())
				if err != nil {
					panic(err)
				}
				nowSales = append(nowSales, atoi)
			}
		})
	})
	return packages, nowSales
}

func buildUpdatedData(packages []string, nowSales []int) map[string]int {
	if len(packages) != len(nowSales) {
		panic("len(packages) != len(nowSales)")
	}

	cache := u.ReadCache()

	updated := make(map[string]int)
	for i := 0; i < len(packages); i++ {
		c, ok := cache[packages[i]]
		if !ok {
			continue
		}
		if nowSales[i] != c {
			updated[packages[i]] = nowSales[i]
		}
	}
	return updated
}

func main() {
	driver := getDriver()
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		panic(err)
	}

	login(page)
	packages, nowSales := getSalesPageData(page)
	updated := buildUpdatedData(packages, nowSales)
	u.SendSlackMessage(updated)

	u.WriteCache(packages, nowSales)
}
