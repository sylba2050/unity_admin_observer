package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	page.FindByName("commit").Click()
	time.Sleep(5000 * time.Millisecond)
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
	dom.Find("table#sales>tbody>tr>td").Each(func(i int, s *goquery.Selection) {
		if title[i] == "Package" {
			packages = append(packages, s.Text())
		}
		if title[i] == "Qty" {
			atoi, err := strconv.Atoi(s.Text())
			if err != nil {
				panic(err)
			}
			nowSales = append(nowSales, atoi)
		}
	})
	return packages, nowSales
}

func sendSlackMessage(packages []string, nowSales []int) {
	text := buildSlackNotificationText(packages, nowSales)
	message := u.Message{Channel: u.Config.SlackChannelName, Text: text}
	rawSendJSON, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", u.Config.SlackURL, bytes.NewBuffer(rawSendJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.Config.SlackToken))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func buildSlackNotificationText(packages []string, nowSales []int) string {
	if len(packages) != len(nowSales) {
		panic(errors.New("len(packages) != len(nowSales)"))
	}
	var text string
	for _, user := range u.Config.SlackUserIDs {
		text += fmt.Sprintf("<@%s> ", user)
	}
	for i := 0; i < len(packages); i++ {
		text += "\n"
		text += fmt.Sprintf("%sが購入されました。現在の累計販売数は%d個です。", packages[i], nowSales[i])
	}

	return text
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
	sendSlackMessage(packages, nowSales)
}
