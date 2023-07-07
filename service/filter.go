package service

import (
	"errors"
	"strings"
	"sync"
)

type FilterService struct {
}

var (
	filterService     *FilterService
	filterServiceOnce sync.Once
)

func GetFilterService() *FilterService {
	filterServiceOnce.Do(func() {
		filterService = &FilterService{}
	})
	return filterService
}

/*
数据过滤规则

1. 买入席位总量 > 卖出席位总量
2. 买一 > 1.5 * 卖一
3. 买二 * 2 > 买一
4. 游资席位分析：买入机构席位大于3
5. 拉萨不要
*/
func (ins *FilterService) FilterLHBCode(rawList []RawDataLHB) ([]string, error) {

	if len(rawList) == 0 {
		return nil, errors.New("raw data is empty")
	}
	RecomendResult := []string{}
	for _, item := range rawList {

		xiweiList := item.XiWei
		//计算总量
		var buyTotal float64
		var sellTotal float64
		var buyOne float64
		var buyTwo float64
		var sellOne float64

		jigouCount := 0
		ifLaSha := false
		for j, xiwei := range xiweiList {
			if j < 5 {
				buyTotal = buyTotal + xiwei.Buy
				if strings.Contains(xiwei.Name, "机构") {
					jigouCount += 1
				}
			} else {
				sellTotal = sellTotal + xiwei.Sell
			}

			if strings.Contains(xiwei.Name, "拉萨") {
				ifLaSha = true
			}
			if j == 0 {
				buyOne = xiwei.Buy
			}
			if j == 1 {
				buyTwo = xiwei.Buy
			}
			if j == 5 {
				sellOne = xiwei.Sell
			}

		}

		//机构数量大于2
		if jigouCount > 2 {
			RecomendResult = append(RecomendResult, item.Code)
			continue
		}
		if ifLaSha {
			continue
		}

		//买方总量大于卖方总量
		if buyTotal < sellTotal {
			continue
		}

		if buyTwo*2 < buyOne {
			continue
		}
		if sellOne*1.5 > buyOne {
			continue
		}

		RecomendResult = append(RecomendResult, item.Code)
	}
	return RecomendResult, nil
}
