package service

type HistoryLHBItem struct {
	Code string
	Date string
	Rid  string
}

// 席位买卖信息
type LHBXWItem struct {
	Name        string
	Buy         float64
	BuyPercent  string
	Sell        float64
	SellPercent string
	Balance     float64
}

// 龙虎榜列表
type LHBCodeItem struct {
	Code         string
	Name         string
	CurrentPrice float64
	Range        string
	TotalAmount  string
	BalanceBuy   string
}

// 抓取出来的数据集合
type RawDataLHB struct {
	Code         string
	Name         string
	CurrentPrice float64
	Range        string
	XiWei        []LHBXWItem
}
