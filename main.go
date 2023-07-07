package main

import (
	"fmt"
	"stock_stat/service"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()

	app.Name = "DragonTiger"

	app.Usage = "get dragon tiger data"

	app.Version = "V0.01"

	now := time.Now()
	date := now.Format("2006-01-02")
	//1.获取当天的龙虎榜列表
	lhbList, err := service.GetClientService().GetLHBFromTHSAndParse(date)
	if err != nil {
		fmt.Println("GetLHBFromTHSAndParse failed ", err)
		return
	}

	if len(lhbList) == 0 {
		fmt.Println("GetLHBFromTHSAndParse return empty")
		return
	}

	//测试，打印一下龙虎榜数据
	fmt.Println("-----------------------------龙虎榜上榜个股---------------------------------")

	for _, item := range lhbList {
		fmt.Printf(item.Code)
		fmt.Printf("    ")
		fmt.Printf(item.Name)
		fmt.Printf("    ")
		fmt.Printf(strconv.FormatFloat(item.CurrentPrice, 'f', 2, 64))
		fmt.Printf("    ")
		fmt.Printf(item.Range)
		fmt.Println()
	}
	rawDataList := []service.RawDataLHB{}
	//fmt.Println("-----------------------------个股龙虎榜席位买入情况---------------------------")
	//根据龙虎榜列表，获取每个code的历史榜单
	for _, item := range lhbList {

		var rawData service.RawDataLHB

		histotyItemList, err := service.GetClientService().GetLHBListByCode(item.Code)
		if err != nil {
			fmt.Println("GetLHBListByCode failed ,code:", item.Code)
			continue
		}
		if len(histotyItemList) == 0 {
			continue
		}
		//获取当天的信息就可以
		hitoryItem := histotyItemList[0]
		xwList, err := service.GetClientService().GetLHBListByCodeAndRid(hitoryItem.Code, hitoryItem.Date, hitoryItem.Rid)
		if err != nil {
			fmt.Println("GetLHBListByCodeAndRid failed ", err)
			return
		}
		if len(xwList) == 0 {
			continue
		}
		rawData.Code = item.Code
		rawData.Name = item.Name
		rawData.CurrentPrice = item.CurrentPrice
		rawData.Range = item.Range
		rawData.XiWei = xwList
		rawDataList = append(rawDataList, rawData)
		/*
		   //测试打印一下
		   fmt.Println("-------------------------------")
		   fmt.Println(item.Code, ":", item.Name)
		   for _, xwItem := range xwList {
		    fmt.Printf(xwItem.Name)
		    fmt.Printf("    ")
		    fmt.Printf(strconv.FormatFloat(xwItem.Buy, 'f', 2, 64))
		    fmt.Printf("    ")
		    fmt.Printf(xwItem.BuyPercent)
		    fmt.Printf("    ")
		    fmt.Printf(strconv.FormatFloat(xwItem.Sell, 'f', 2, 64))
		    fmt.Printf("    ")
		    fmt.Printf(xwItem.SellPercent)
		    fmt.Printf("    ")
		    fmt.Printf(strconv.FormatFloat(xwItem.Balance, 'f', 2, 64))
		    fmt.Println()
		   }*/
	}

	//做规则过滤
	result, err := service.GetFilterService().FilterLHBCode(rawDataList)
	if err != nil {
		fmt.Println("FilterLHBCode", err)
		return
	}

	fmt.Println(result)

	return
}
