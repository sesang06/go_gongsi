package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	money          = 1000000
	sleep_duration = time.Second * 40
)

func CalDiff(original, new []NoticePost) []NoticePost {
	if len(new) == 0 {
		var emptyObject []NoticePost
		return emptyObject
	}
	newLen := len(new)

	var newLenIndex int = 0
	var diffObject []NoticePost

	if len(original) == 0 {
		var emptyObject []NoticePost
		return emptyObject
	}
	lastNoticeObject := original[0]

	for newLenIndex < newLen {
		currentPost := new[newLenIndex]
		if currentPost.ID != lastNoticeObject.ID {
			diffObject = append(diffObject, currentPost)
			newLenIndex += 1
		} else {
			break
		}
	}
	return diffObject

}


func CalNotiDiff(original, new []Notification) []Notification {
	if len(new) == 0 {
		var emptyObject []Notification
		return emptyObject
	}
	newLen := len(new)

	var newLenIndex int = 0
	var diffObject []Notification

	if len(original) == 0 {
		var emptyObject []Notification
		return emptyObject
	}
	lastNoticeObject := original[0]

	for newLenIndex < newLen {
		currentPost := new[newLenIndex]
		if currentPost.ID != lastNoticeObject.ID {
			diffObject = append(diffObject, currentPost)
			newLenIndex += 1
		} else {
			break
		}
	}
	return diffObject

}

type GonsiResult struct {
	err   error
	value []NoticePost
}

type NoticeResult struct {
	err error
	value []Notification
}
type UnknownError struct {
}

func (e *UnknownError) Error() string {
	return "POST IS NILL"
}


func DownloadNotificationFile(url string) <-chan NoticeResult {
	out := make(chan NoticeResult)
	go func() {
		defer close(out)
		func(url string) {

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				out <- NoticeResult{value: nil, err: err}
				return
			}

			//필요시 헤더 추가 가능
			req.Header.Add("Cache-Control", "no-cache, must-revalidate")

			// Client객체에서 Request 실행
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				out <- NoticeResult{value: nil, err: err}
				return
			}

			if res.StatusCode >= 400 {
				out <- NoticeResult{value: nil, err: &UnknownError{}}
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				out <- NoticeResult{value: nil, err: err}
				return
			}
			object := NotificationData{}

			jsonErr := json.Unmarshal(body, &object)

			if jsonErr != nil {
				out <- NoticeResult{value: nil, err: jsonErr}
				return
			}
			posts := object.Data.List
			if posts == nil {
				out <- NoticeResult{value: nil, err: &UnknownError{}}
				return
			}
			out <- NoticeResult{value: posts, err: nil}
		}(url)
	}()
	return out
}

func DownloadGongsiFile(url string) <-chan GonsiResult {
	out := make(chan GonsiResult)
	go func() {
		defer close(out)
		func(url string) {

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				out <- GonsiResult{value: nil, err: err}
				return
			}

			//필요시 헤더 추가 가능
			req.Header.Add("Cache-Control", "no-cache, must-revalidate")

			// Client객체에서 Request 실행
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				out <- GonsiResult{value: nil, err: err}
				return
			}

			if res.StatusCode >= 400 {
				out <- GonsiResult{value: nil, err: &UnknownError{}}
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				out <- GonsiResult{value: nil, err: err}
				return
			}
			object := NoticeObject{}

			jsonErr := json.Unmarshal(body, &object)

			if jsonErr != nil {
				out <- GonsiResult{value: nil, err: jsonErr}
				return
			}
			posts := object.Data.Posts
			if posts == nil {
				out <- GonsiResult{value: nil, err: &UnknownError{}}
				return
			}
			out <- GonsiResult{value: object.Data.Posts, err: nil}
		}(url)
	}()
	return out
}

func startTrading(upbitTrader *UpbitTrader, diff []NoticePost) {
	for _, post := range diff {
		ticker := "KRW-" + post.Assets
		if strings.Contains(post.Text, "기공개") {
			go upbitTrader.buyAndSell(ticker, 10000, sleep_duration)
		} else {
			go upbitTrader.buyAndSell(ticker, money, sleep_duration)
		}
	}
	go func() {
		for _, post := range diff {
			SendMessage("<!here> [GO] 공시감지 : " + post.Assets + post.Text)
		}
	}()
}

func startCrawling() {

	const gongsiUrl = "https://project-team.upbit.com/api/v1/disclosure?region=kr&per_page=10"
	var gonsiResultOut GonsiResult
	for true {
		gonsiResultOut := <-DownloadGongsiFile(gongsiUrl)
		if gonsiResultOut.err != nil {
			fmt.Println(gonsiResultOut.err)
		} else {
			break
		}
	}

	var currentGongsi = gonsiResultOut.value
	SendMessage("공시(Go) 시작")
	var cnt int = 0
	var has_error_occured = false
	upbitTrader := NewUpbitTrader()
	for true {
		time.Sleep(time.Millisecond * 2000)
		newPostResult := <-DownloadGongsiFile(gongsiUrl)
		if newPostResult.err != nil {
			if has_error_occured {
				continue
			}
			has_error_occured = true
			SendMessage("공시(Go) ERROR\n" + newPostResult.err.Error())
			continue
		}
		newPosts := newPostResult.value
		if len(newPosts) == 0 {
			fmt.Print("EMPTY")
			continue
		}
		diff := CalDiff(currentGongsi, newPosts)
		diffLen := len(diff)
		if diffLen > 0 && diffLen != 10 {
			startTrading(upbitTrader, diff)
		}

		if has_error_occured {
			has_error_occured = false
			SendMessage("<!here> Go 봇 정상화")
		}
		cnt += 1
		cnt %= 7200
		if cnt == 0 {
			var text string = ""
			for _, post := range newPosts {
				text += post.Text
				text += "\n"
			}
			SendMessage("공시(Go)\n" + text)
		}
		currentGongsi = newPosts
	}
}

func startCrawlingNotification() {
	const noticeUrl = "https://api-manager.upbit.com/api/v1/notices?page=1&per_page=10&thread_name=general"
	var notificationResultOut NoticeResult
	for true {
		notificationResultOut := <-DownloadNotificationFile(noticeUrl)
		if notificationResultOut.err != nil {
			fmt.Println(notificationResultOut.err)
		} else {
			break
		}
	}

	var currentGongsi = notificationResultOut.value
	SendMessage("공지(Go) 시작")
	var cnt int = 0
	var has_error_occured = false
	upbitTrader := NewUpbitTrader()
	for true {
		time.Sleep(time.Millisecond * 2000)
		newPostResult := <-DownloadNotificationFile(noticeUrl)
		if newPostResult.err != nil {
			if has_error_occured {
				continue
			}
			has_error_occured = true
			SendMessage("공지(Go) ERROR\n" + newPostResult.err.Error())
			continue
		}
		newPosts := newPostResult.value
		if len(newPosts) == 0 {
			fmt.Print("EMPTY")
			continue
		}
		diff := CalNotiDiff(currentGongsi, newPosts)
		diffLen := len(diff)
		if diffLen > 0 && diffLen != 10 {
			startTradeForNotification(upbitTrader, diff)
		}

		if has_error_occured {
			has_error_occured = false
			SendMessage("<!here> Go 봇 정상화")
		}
		cnt += 1
		cnt %= 7200
		if cnt == 0 {
			var text string = ""
			for _, post := range newPosts {
				text += post.Title
				text += "\n"
			}
			SendMessage("공시(Go)\n" + text)
		}
		currentGongsi = newPosts
	}
}

func main() {
	go func() {
		startCrawling()
	}()
	startCrawlingNotification()
}
