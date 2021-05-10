package main

import "strings"


func parse_ticker(text string) []string {
	ticker_candidate := []string{"NEO", "MTL", "LTC", "XRP", "ETC", "OMG", "SNT", "WAVES", "XEM", "QTUM", "LSK", "STEEM", "XLM", "ARDR", "KMD", "ARK", "STORJ", "GRS", "REP", "EMC2", "ADA", "SBD", "POWR", "BTG", "ICX", "EOS", "TRX", "SC", "IGNIS", "ONT", "ZIL", "POLY", "ZRX", "SRN", "LOOM", "BCH", "ADX", "BAT", "IOST", "DMT", "RFR", "CVC", "IQ", "IOTA", "MFT", "ONG", "GAS", "UPP", "ELF", "KNC", "BSV", "THETA", "EDR", "QKC", "BTT", "MOC", "ENJ", "TFUEL", "MANA", "ANKR", "NPXS", "AERGO", "ATOM", "TT", "CRE", "SOLVE", "MBL", "TSHP", "WAXP", "HBAR", "MED", "MLK", "STPT", "ORBS", "VET", "CHZ", "PXL", "STMX", "DKA", "HIVE", "KAVA", "AHT", "SPND", "LINK", "XTZ", "BORA", "JST", "CRO", "TON", "SXP", "LAMB", "HUNT", "MARO", "PLA", "DOT", "SRM", "MVL", "PCI", "STRAX", "AQT", "BCHA", "GLM", "QTCON", "SSX", "META", "OBSR", "FCT2", "LBC", "CBK", "SAND", "HUM", "DOGE"}
	var containing_tickers []string
	for _, ticker := range ticker_candidate {
		if strings.Contains(text, ticker) {
			containing_tickers = append(containing_tickers, ticker)
		}
	}
	return containing_tickers
}

func startTradeForNotification(upbitTrader *UpbitTrader, diff []Notification) {
	for _, post := range diff {
		tickers := parse_ticker(post.Title)
		for _, ticker := range tickers {
			SendMessage("<!here> [Go] 공지감지" + ticker)
			tickerWithKRW := "KRW-" + ticker
			upbitTrader := NewUpbitTrader()
			go upbitTrader.buyAndSell(tickerWithKRW, 10000, sleep_duration)
		}
	}
	go func() {
		for _, post := range diff {
			SendMessage("<!here> [GO] 공지감지 : " + post.Title)
		}
	}()
}
