package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gammazero/workerpool"
	"github.com/gocolly/colly/v2"
)

var userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36"
var queue int

func main() {
	fmt.Println(`
	 Utku Sen's
	\ \     / /
	 \ \   / /
	  \ \_/ /
	   \   /
	socialhunter
	   / _ \
	  / / \ \
	 / /   \ \
	/_/     \_\
	utkusen.com
`)
	urlFile := flag.String("f", "", "Path of the URL file")
	numWorker := flag.Int("w", 5, "Number of worker.")
	flag.Parse()
	if *urlFile == "" {
		fmt.Println("Please specify all arguments!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	file, err := ioutil.ReadFile(*urlFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	urls := strings.Split(string(file), "\n")
	queue = len(urls)
	fmt.Println("Total URLs:", queue)
	wp := workerpool.New(*numWorker)

	for _, url := range urls {
		url := url
		wp.Submit(func() {
			fmt.Println("Checking:", url)
			action(url)
		})

	}
	wp.StopWait()

	color.Cyan("Scan Completed")
}

func action(url string) {
	sl := visitor(url, 10)
	checkTakeover(removeDuplicateStr(sl))
	color.Magenta("Finished Checking: " + url)
	queue--
	fmt.Println("Remaining URLs:", queue)
}

func stringInSlice(a string, list *[]string) bool {
	for _, b := range *list {
		if b == a {
			return true
		}
	}
	return false
}

func checkTakeover(socialLinks []string) {
	var alreadyChecked []string
	for _, value := range socialLinks {
		foundLink := strings.Split(value, "|")[0]
		socialLink := strings.Split(value, "|")[1]
		if stringInSlice(socialLink, &alreadyChecked) {
			continue
		}
		alreadyChecked = append(alreadyChecked, socialLink)
		if len(socialLink) > 60 || strings.Contains(socialLink, "intent/tweet") || strings.Contains(socialLink, "twitter.com/share") || strings.Contains(socialLink, "twitter.com/privacy") || strings.Contains(socialLink, "facebook.com/home") || strings.Contains(socialLink, "instagram.com/p/") {
			continue
		}
		u, err := url.Parse(socialLink)
		if err != nil {
			continue
		}
		domain := u.Host
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if strings.Contains(domain, "facebook.com") {
			if strings.Count(socialLink, ".") > 1 {
				socialLink = "https://" + strings.Split(socialLink, ".")[1] + "." + strings.Split(socialLink, ".")[2]
			}
			socialLink = strings.Replace(socialLink, "www.", "", -1)
			tempLink := strings.Replace(socialLink, "facebook.com", "tr-tr.facebook.com", -1)
			resp, err := http.Get(tempLink)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			if strings.Contains(string(body), "Sayfa Bulunamadı") {
				color.Green("Possible Takeover: " + socialLink + " at " + foundLink)

			}

		}
		if strings.Contains(domain, "tiktok.com") {
			if strings.Count(strings.Replace(socialLink, "www.", "", -1), ".") > 1 {
				continue
			}
			client := &http.Client{Transport: tr}

			req, err := http.NewRequest("GET", socialLink, nil)
			if err != nil {
				continue
			}

			req.Header.Set("User-Agent", userAgent)

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 404 {
				color.Green("Possible Takeover: " + socialLink + " at " + foundLink)
			}
		}
		if strings.Contains(domain, "instagram.com") {

			if strings.Count(strings.Replace(socialLink, "www.", "", -1), ".") > 1 {
				continue
			}
			if !strings.Contains(socialLink, "instagram.com/") {
				continue
			}
			tempLink := "https://www.picuki.com/profile/" + strings.Split(socialLink, "instagram.com/")[1]
			client := &http.Client{Transport: tr}
			req, err := http.NewRequest("GET", tempLink, nil)
			if err != nil {
				continue
			}

			req.Header.Set("User-Agent", userAgent)

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 404 {
				color.Green("Possible Takeover: " + socialLink + " at " + foundLink)
			}
		}
		if strings.Contains(domain, "twitter.com") {
			if strings.Count(strings.Replace(socialLink, "www.", "", -1), ".") > 1 {
				continue
			}
			u, err := url.Parse(socialLink)
			userName := u.Path
			tempLink := "https://nitter.net" + userName
			client := &http.Client{}
			req, err := http.NewRequest("GET", tempLink, nil)
			if err != nil {
				continue
			}

			req.Header.Set("User-Agent", userAgent)

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 404 {
				color.Green("Possible Takeover: " + socialLink + " at " + foundLink)
			}
		}
	}
	return
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}

	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func visitor(visitURL string, maxDepth int) []string {
	socialDomains := []string{"twitter.com", "instagram.com", "facebook.com", "twitch.tv", "tiktok.com"}
	var socialLinks []string
	var visitedLinks []string
	denyList := []string{".js", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".mp4", ".webm", ".mp3", ".csv", ".ogg", ".wav", ".flac", ".aac", ".wma", ".wmv", ".avi", ".mpg", ".mpeg", ".mov", ".mkv", ".zip", ".rar", ".7z", ".tar", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".rtf", ".odt", ".ods", ".odp", ".odg", ".odf", ".odb", ".odc", ".odm", ".avi", ".mpg", ".mpeg", ".mov", ".mkv", ".zip", ".rar", ".7z", ".tar", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".rtf", ".odt", ".ods", ".odp", ".odg", ".odf", ".odb", ".odc", ".odm", ".mp4", ".webm", ".mp3", ".ogg", ".wav", ".flac", ".aac", ".wma", ".wmv", ".avi", ".mpg", ".mpeg", ".mov", ".mkv", ".zip", ".rar", ".7z", ".tar", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".rtf", ".odt", ".ods", ".odp", ".odg", ".odf", ".odb", ".odc", ".odm", ".mp4", ".webm", ".mp3", ".ogg", ".wav", ".flac", ".aac", ".wma", ".wmv", ".avi", ".mpg", ".mpeg", ".mov", ".mkv", ".zip", ".rar", ".7z", ".tar", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".rtf", ".odt"}

	c := colly.NewCollector()
	c.UserAgent = userAgent
	c.SetRequestTimeout(5 * time.Second)
	c.MaxDepth = maxDepth
	c.AllowURLRevisit = false //there is a bug in colly that prevents this from working. We have to check it manually
	u, err := url.Parse(visitURL)
	if err != nil {
		panic(err)
	}
	domain := u.Host
	path := u.Path
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		u2, err := url.Parse(link)
		if err != nil {
			panic(err)
		}
		linkDomain := u2.Host
		for _, domain := range socialDomains {
			if strings.Contains(linkDomain, domain) {
				socialLinks = append(socialLinks, e.Request.URL.String()+"|"+link)
			}
		}
		if strings.Contains(linkDomain, domain) {
			visitFlag := true
			for _, extension := range denyList {
				if strings.Contains(strings.ToLower(link), extension) {
					visitFlag = false
				}
			}
			for _, value := range visitedLinks {
				if strings.ToLower(link) == value {
					visitFlag = false
				}
			}

			if !strings.HasPrefix(u2.Path, path) {
				visitFlag = false
			}

			if visitFlag == true {
				visitedLinks = append(visitedLinks, link)
				e.Request.Visit(link)

			}
		}

	})

	c.Visit(visitURL)
	return socialLinks
}
