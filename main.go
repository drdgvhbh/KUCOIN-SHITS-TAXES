package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

	startDate, err := time.Parse("2006-01-02", os.Getenv("START_DATE"))
	if err != nil {
		log.Fatalln(err, "start date is invalid")
	}
	endDate, err := time.Parse("2006-01-02", os.Getenv("END_DATE"))
	if err != nil {
		log.Fatalln(err, "end date is invalid")
	}

	recordFile, err := os.Create("orders.csv")
	if err != nil {
		log.Fatalln(err, "failed to create orders.csv")
	}

	writer := csv.NewWriter(recordFile)
	err = writer.Write([]string{
		"id", "symbol", "op_type", "type", "side", "price",
		"size", "funds", "deal_funds", "deal_size", "fee",
		"fee_currency", "stp", "stop", "stop_triggered", "stop_price", "time_in_force", "post_only",
		"hidden", "iceberg", "visible_size", "cancel_after", "channel", "client_oid", "remark",
		"tags", "is_active", "cancel_exist", "created_at", "trade_type", "status", "fail_msg",
	})
	if err != nil {
		log.Fatalln(err, "failed to write to orders.csv")
	}

	currentStartDate := startDate
	currentEndDate := startDate.AddDate(0, 0, 7)

	paginationParams := kucoin.PaginationParam{CurrentPage: 1, PageSize: 500}

	for currentStartDate.Before(endDate) {
		params := make(map[string]string)
		params["startAt"] = fmt.Sprintf("%d", currentStartDate.UnixNano()/1e6)
		params["endAt"] = fmt.Sprintf("%d", currentEndDate.UnixNano()/1e6)
		params["tradeType"] = "TRADE"

		rsp, err := s.Orders(params, &paginationParams)
		if err != nil {
			log.Fatalln(err, "failed to retrieved orders")
		}
		os := kucoin.OrdersModel{}
		pa, err := rsp.ReadPaginationData(&os)
		if err != nil {
			log.Fatalln(err, "failed to read pagination data")
		}
		log.Printf("Total num: %d, current page: %d, total page: %d, page size: %d", pa.TotalNum, pa.CurrentPage, pa.TotalPage, pa.PageSize)
		log.Printf("Start Date: %s, End Date: %s", currentStartDate.Format("2006-01-02"), currentEndDate.Format("2006-01-02"))
		for _, o := range os {
			err = writer.Write([]string{
				o.Id, o.Symbol, o.OpType, o.Type, o.Side, o.Price,
				o.Size, o.Funds, o.DealFunds, o.DealSize, o.Fee,
				o.FeeCurrency, o.Stp, o.Stop, strconv.FormatBool(o.StopTriggered), o.StopPrice,
				o.TimeInForce, strconv.FormatBool(o.PostOnly), strconv.FormatBool(o.Hidden),
				strconv.FormatBool(o.IceBerg), o.VisibleSize, strconv.FormatUint(o.CancelAfter, 10),
				o.Channel, o.ClientOid, o.Remark, o.Tags, strconv.FormatBool(o.IsActive),
				strconv.FormatBool(o.CancelExist), strconv.FormatInt(o.CreatedAt, 10),
				o.TradeType, o.Status, o.FailMsg,
			})
		}

		if pa.CurrentPage < pa.TotalPage {
			paginationParams.CurrentPage++
		} else {
			paginationParams = kucoin.PaginationParam{CurrentPage: 1, PageSize: 500}
			currentStartDate = currentStartDate.AddDate(0, 0, 7)
			currentEndDate = currentEndDate.AddDate(0, 0, 7)
		}
	}

	writer.Flush()
	err = recordFile.Close()
	if err != nil {
		log.Fatalln(err, "failed to close orders.csv")
	}
}
