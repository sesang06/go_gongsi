package main
import (
	"fmt"
	"sync"
	"testing"
)

func TestParse(t *testing.T) {
	text := "[디지털 자산] 펀디엑스(NPXS) 토큰 스왑 및 심볼 변경에 따른 입"
	tickers := parse_ticker(text)
	fmt.Println(tickers)
}

func TestStartNotification(t *testing.T) {
	notice := NoticePost{}
	notice.Assets = "AHT"
	noticeData := []NoticePost{notice}
	upbitTrader := NewUpbitTrader()
	startTrading(upbitTrader, noticeData)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}