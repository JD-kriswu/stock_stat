package service

import (
	"errors"
	"fmt"
	"net/http"
	"stock_stat/utils"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ClientService struct {
}

var (
	clientService     *ClientService
	clientServiceOnce sync.Once
)

func GetClientService() *ClientService {
	clientServiceOnce.Do(func() {
		clientService = &ClientService{}
	})
	return clientService
}

func (ins *ClientService) GetLHBFromTHSAndParse(date string) ([]LHBCodeItem, error) {

	client := &http.Client{}
	//now := time.Now()
	//url := "https://data.10jqka.com.cn/ifmarket/lhbggxq/report/" + now.Format("2006-01-02") + "/"
	url := "https://data.10jqka.com.cn/ifmarket/lhbggxq/report/" + "2023-06-15" + "/"
	//fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,immaaaaaaage/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Cookie", "Hm_lvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; v=A_HvWkPOMePy8p33zWJidJfMBnaO3mVQD1IJZNMG7bjX-h_oGy51IJ-iGTBg")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	//req.Header.Add("Referrer Policy", "strict-origin-when-cross-origin")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("GetLHBFromTHS failed ")
		return nil, err

	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Get Doc from Body failed")
		return nil, err
	}
	result := []LHBCodeItem{}

	//获取列表
	doc.Find(".page-table .twrap .m-table tr").Each(func(i int, s *goquery.Selection) {
		//每一行 ，一个代码对象
		var codeItem LHBCodeItem
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if i != 0 {
				content, _ := utils.GbkToUtf8([]byte(s.Text()))
				switch i {
				case 1:
					codeItem.Code = string(content)
					break

				case 2:
					codeItem.Name = string(content)
					break
				case 3:
					codeItem.CurrentPrice, _ = strconv.ParseFloat(string(content), 64)
					break
				case 4:
					codeItem.Range = string(content)
					break
				case 5:
					codeItem.TotalAmount = string(content)
					break
				case 6:
					codeItem.BalanceBuy = string(content)
					break
				}

			}
		})
		if codeItem.Code != "" {
			result = append(result, codeItem)
		}

	})
	return result, nil

}

// 获取每一个code对应的历史榜单
func (ins *ClientService) GetLHBListByCode(code string) ([]HistoryLHBItem, error) {

	if code == "" {
		fmt.Println("err code :", code)
		return nil, errors.New("err code ")
	}

	url := "http://data.10jqka.com.cn/market/lhbgg/code/" + code
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,immaaaaaaage/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Cookie", "Hm_lvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; v=A_HvWkPOMePy8p33zWJidJfMBnaO3mVQD1IJZNMG7bjX-h_oGy51IJ-iGTBg")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	//req.Header.Add("Referrer Policy", "strict-origin-when-cross-origin")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("GetLHBListByCode failed ")
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Get Doc from Body failed")
		return nil, err
	}
	result := []HistoryLHBItem{}
	//获取历史列表
	doc.Find("#ggsj tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			//获取最后一个列数据，拿到三个属性数据
			if i == 7 {
				attrCode, _ := s.Find("a").Attr("code")
				attrDate, _ := s.Find("a").Attr("date")
				attrRid, _ := s.Find("a").Attr("rid")

				historyItem := HistoryLHBItem{
					Code: attrCode,
					Date: attrDate,
					Rid:  attrRid,
				}
				result = append(result, historyItem)

			}
		})
	})

	return result, nil
}

// 获取最近一个交易日的lhb席位信息
// http://data.10jqka.com.cn/ifmarket/getnewlh/code/300288/date/2023-05-17/rid/44/
func (ins *ClientService) GetLHBListByCodeAndRid(code, date, rid string) ([]LHBXWItem, error) {
	if code == "" || date == "" || rid == "" {
		fmt.Println("err code date rid", code, date, rid)
		return nil, errors.New("err params")
	}

	url := "http://data.10jqka.com.cn/ifmarket/getnewlh/code/" + code + "/date/" + date + "/rid/" + rid + "/"
	//fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cookie", "Hm_lvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1684246227; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1684246227; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1684246227; v=A_HvWkPOMePy8p33zWJidJfMBnaO3mVQD1IJZNMG7bjX-h_oGy51IJ-iGTBg")
	//req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	//GetLHBListByCodeAndRidreq.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Upgrade-Insecure-Requests", "1")
	//req.Header.Add("Referrer Policy", "strict-origin-when-cross-origin")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("GetLHBListByCodeAndRid failed ")
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Get Doc from Body failed")
		return nil, err
	}
	result := []LHBXWItem{}
	//解析出来席位
	doc.Find(".m_table tr").Each(func(i int, s *goquery.Selection) {

		var xwItem LHBXWItem
		s.Find("td").Each(func(i int, s *goquery.Selection) {

			if i != 0 {
				var content []byte
				if i == 1 {
					content, _ = utils.GbkToUtf8([]byte(s.Find("a").Text()))
				} else {
					content, _ = utils.GbkToUtf8([]byte(s.Text()))
				}

				switch i {
				case 1:
					xwItem.Name = string(content)
					break
				case 2:
					xwItem.Buy, _ = strconv.ParseFloat(string(content), 64)
					break
				case 3:
					xwItem.BuyPercent = string(content)
					break
				case 4:
					xwItem.Sell, _ = strconv.ParseFloat(string(content), 64)
					break
				case 5:
					xwItem.SellPercent = string(content)
					break
				case 6:
					xwItem.Balance, _ = strconv.ParseFloat(string(content), 64)
					break
				}

			}
		})
		if xwItem.Name != "" {
			result = append(result, xwItem)
		}

	})
	return result, nil
}
