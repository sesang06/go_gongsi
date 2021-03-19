package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	money          = 10000
	sleep_duration = time.Second * 30
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

type Result struct {
	err   error
	value []NoticePost
}

type UnknownError struct {
}

func (e *UnknownError) Error() string {
	return "POST IS NILL"
}

func DownloadFile(url string) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		func(url string) {

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				out <- Result{value: nil, err: err}
				return
			}

			//필요시 헤더 추가 가능
			req.Header.Add("Cache-Control", "no-cache, must-revalidate")

			// Client객체에서 Request 실행
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				out <- Result{value: nil, err: err}
				return
			}

			if res.StatusCode >= 400 {
				out <- Result{value: nil, err: &UnknownError{}}
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				out <- Result{value: nil, err: err}
				return
			}
			object := NoticeObject{}

			jsonErr := json.Unmarshal(body, &object)

			if jsonErr != nil {
				out <- Result{value: nil, err: jsonErr}
				return
			}
			posts := object.Data.Posts
			if posts == nil {
				out <- Result{value: nil, err: &UnknownError{}}
				return
			}
			out <- Result{value: object.Data.Posts, err: nil}
		}(url)
	}()
	return out
}

func startTrading(upbitTrader *UpbitTrader, diff []NoticePost) {
	for _, post := range diff {
		ticker := "KRW-" + post.Assets
		go upbitTrader.buyAndSell(ticker, money, sleep_duration)
	}
	go func() {
		for _, post := range diff {
			SendMessage("<!here> [GO] 공시감지 : " + post.Assets + post.Text)
		}
	}()
}

func startCrawling() {

	const url = "https://project-team.upbit.com/api/v1/disclosure?region=kr&per_page=10"
	var resultOut Result
	for true {
		resultOut := <-DownloadFile(url)
		if resultOut.err != nil {
			fmt.Println(resultOut.err)
		} else {
			break
		}
	}
	var current = resultOut.value
	SendMessage("공시(Go) 시작")
	var cnt int = 0
	var has_error_occured = false
	upbitTrader := NewUpbitTrader()
	for true {
		time.Sleep(time.Millisecond * 2000)
		newPostResult := <-DownloadFile(url)
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
		diff := CalDiff(current, newPosts)
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
		current = newPosts
	}
}

func main() {
	startCrawling()
}
