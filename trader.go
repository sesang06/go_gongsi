package main

import (
	"fmt"
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
		case <- time.After(duration): upbitTrader.sellAll(ticker)
	}
}

func (upbitTrader UpbitTrader) getBalances() {
	accounts, model, err := upbitTrader.upbit.GetAccounts()

	fmt.Print(accounts)
	fmt.Print(model)
	fmt.Print(err)

	fmt.Print(upbitTrader.upbit.GetAccounts())
}
//func main() {
//	u := NewUpbitTrader()
//	//u.getBalance()
//	//u.buy("KRW-ETH", 6000)
//	u.sellAll("KRW-ETH")
//}
