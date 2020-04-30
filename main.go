package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Kucoin/kucoin-go-sdk"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	kucoin.DebugMode = true
	kucoin.SetLoggerDirectory(exPath)
	s := kucoin.NewApiService(
		// kucoin.ApiBaseURIOption("https://api.kucoin.com"),
		kucoin.ApiKeyOption(os.Getenv("API_KEY")),
		kucoin.ApiSecretOption(os.Getenv("SECRET")),
		kucoin.ApiPassPhraseOption(os.Getenv("PASSWORD")),
	)

	params := make(map[string]string)
	params["startAt"] = "1559606400000"
	params["endAt"] = "1560211200000"
	params["tradeType"] = "TRADE"
	rsp, err := s.Orders(params, &kucoin.PaginationParam{CurrentPage: 1, PageSize: 10})
	if err != nil {
		// Handle error
		return
	}

	os := kucoin.OrdersModel{}
	pa, err := rsp.ReadPaginationData(&os)
	if err != nil {
		// Handle error
		return
	}
	log.Printf("Total num: %d, total page: %d", pa.TotalNum, pa.TotalPage)
	for _, o := range os {
		log.Printf("Order: %s, %s, %s", o.Id, o.Type, o.Price)
	}
}
