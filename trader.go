package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type UpbitTrader struct {
	upbit *Upbit
}

func NewUpbitTrader() *UpbitTrader {
	accessKey := g_accesskey
	secretKey := g_secretkey
	u := NewUpbit(accessKey, secretKey)
	return &UpbitTrader{
		u,
	}
}

func (upbitTrader UpbitTrader) buy(ticker string, money int) {
	order, _, error := upbitTrader.upbit.PurchaseOrder(ticker, "", strconv.Itoa(money), "price", "")
	fmt.Print(order)
	fmt.Print(error)
}

func (upbitTrader UpbitTrader) getBalance(ticker string) string {
	slice := strings.Split(ticker, "-")
	currency := slice[1]
	accounts, _, err := upbitTrader.upbit.GetAccounts()

	if err != nil {
		return "0"
	}
	for _, account := range accounts {
		if account.Currency == currency {
			return account.Balance
		}
	}
	return "0"
}

func (upbitTrader UpbitTrader) getBalanceAndAvgBuyPrice(ticker string) (string, string) {
	slice := strings.Split(ticker, "-")
	currency := slice[1]
	accounts, _, err := upbitTrader.upbit.GetAccounts()

	if err != nil {
		return "0", "0"
	}
	for _, account := range accounts {
		if account.Currency == currency {
			return account.Balance, account.AvgBuyPrice
		}
	}
	return "0", "0"
}

func (upbitTrader UpbitTrader) sellAll(ticker string) {
	balance := upbitTrader.getBalance(ticker)
	_, _, err := upbitTrader.upbit.SellOrder(ticker, balance, "", "market", "")
	if err != nil {
		fmt.Println(err)
	}
}

func (upbitTrader UpbitTrader) buyAndSell(ticker string, money int, duration time.Duration) {
	upbitTrader.buy(ticker, money)
	select {
	case <-time.After(time.Second):
		upbitTrader.getLimitSell(ticker)
	}
	select {
	case <-time.After(duration - time.Second):
		upbitTrader.cancelSellOrders(ticker)
	}
	select {
	case <-time.After(time.Second):
		upbitTrader.sellAll(ticker)
	}
}

func (upbitTrader UpbitTrader) getLimitSell(ticker string) {
	balance, avgPricestr := upbitTrader.getBalanceAndAvgBuyPrice(ticker)
	avgPriceFloat, _ := strconv.ParseFloat(avgPricestr, 64)
	price := fmt.Sprintf("%f", getTickerSize(avgPriceFloat*1.10))
	fmt.Println(price)
	fmt.Println(balance)
	_, _, err := upbitTrader.upbit.SellOrder(ticker, balance, price, "limit", "")
	if err != nil {
		fmt.Println(err)
	}
}

func (upbitTrader UpbitTrader) cancelSellOrders(ticker string) {
	orders, _, _ := upbitTrader.upbit.GetOrders(ticker, "wait", []string{}, []string{}, "", "")
	for _, order := range orders {
		if order.Side == "ask" {
			upbitTrader.upbit.CancelOrder(order.UUID)
		}
	}
}

func getTickerSize(price float64) float64 {
	if price >= 2000000 {
		return math.Round(price/1000) * 1000
	}

	if price >= 1000000 {
		return math.Round(price/500) * 500
	}

	if price >= 500000 {
		return math.Round(price/100) * 100
	}

	if price >= 100000 {
		return math.Round(price/50) * 50
	}
	if price >= 10000 {
		return math.Round(price/10) * 10
	}
	if price >= 1000 {
		return math.Round(price/5) * 5
	}
	if price >= 100 {
		return math.Round(price/1) * 1
	}
	if price >= 10 {
		return math.Round(price/0.1) * 0.1
	}
	return math.Round(price/0.01) * 0.01
}

//func (upbitTrader UpbitTrader) getSellableMoney(ticker string) {
//
//}

func (upbitTrader UpbitTrader) getBalances() {
	accounts, model, err := upbitTrader.upbit.GetAccounts()

	fmt.Print(accounts)
	fmt.Print(model)
	fmt.Print(err)

	fmt.Print(upbitTrader.upbit.GetAccounts())
}

//func main() {
//	u := NewUpbitTrader()
//	//u.getLimitSell("KRW-BTC")
//	//u.cancelSellOrders("KRW-BTC")
//	u.buyAndSell("KRW-NPXS", 10000, time.Second * 50)
//	//select {
//	//case <- time.After(time.Second): fmt.Println("HELLO WOLRD")
//	//}
//	//select {
//	//case <- time.After(time.Second): fmt.Println("HELLO WOLRD")
//	//}
//
//}
