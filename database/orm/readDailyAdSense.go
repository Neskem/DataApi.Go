package orm

import (
	"DataApi.Go/database/models/AdSense"
	"DataApi.Go/lib/common"
	"fmt"
	"github.com/jinzhu/gorm"
	"sync"
)

type AdSenseReportDaily = AdSense.AdSenseReportDaily
type AdSenseRevenue = AdSense.AdSenseRevenue

func QueryAdSenseReport(db *gorm.DB, accountId string, startDate string, endDate string) []common.JSON {
	table := "adsense_report_daily"
	//db.Table(table).Where("account_id = ? AND date BETWEEN ? AND ?", accountId, startDate, endDate).Find(&adSenseReportDaily)
	rows, err := db.Table(table).Model(&AdSenseReportDaily{}).Where("account_id = ? AND date BETWEEN ? AND ?", accountId, startDate, endDate).Rows()
	var rowsList []common.JSON
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var adSenseReportDaily AdSenseReportDaily
		err := db.ScanRows(rows, &adSenseReportDaily)
		if err != nil {
			fmt.Println(err)
		}
		rowsList = append(rowsList, common.JSON{
			"account_id": adSenseReportDaily.AccountId,
			"ad_client_id": adSenseReportDaily.AdClientId,
			"ad_exchange_clicks": adSenseReportDaily.AdExchangeClicks,
			"ad_exchange_impression_rpm": adSenseReportDaily.AdExchangeImpressionRpm,
			"ad_exchange_impressions": adSenseReportDaily.AdExchangeImpressions,
			"clicks": adSenseReportDaily.Clicks,
			"customer_ad_exchange_estimated_revenue": adSenseReportDaily.CustomerAdExchangeEstimatedRevenue,
			"date": adSenseReportDaily.Date,
			"domain_name": adSenseReportDaily.DomainName,
			"earnings": adSenseReportDaily.Earnings,
			"impression_rpm": adSenseReportDaily.ImpressionRpm,
			"impressions": adSenseReportDaily.MatchedAdRequests,
			"network_code": adSenseReportDaily.NetworkCode,
			"page_rpm": adSenseReportDaily.PageRpm,
			"page_views": adSenseReportDaily.PageViews,
		})
	}
	fmt.Println(rowsList)
	return rowsList
}

func QueryAdSenseRevenue(db *gorm.DB, accountId string, startDate string, endDate string) []common.JSON {
	table := "adsense_report_daily"
	//db.Table(table).Where("account_id = ? AND date BETWEEN ? AND ?", accountId, startDate, endDate).Find(&adSenseReportDaily)
	rows, err := db.Table(table).Model(&AdSenseRevenue{}).Select("account_id, sum(customer_ad_exchange_estimated_revenue) as customer_ad_exchange_estimated_revenue").Where("account_id = ? AND date BETWEEN ? AND ?", accountId, startDate, endDate).Group("account_id").Rows()
	var rowsList []common.JSON
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var adSenseRevenue AdSenseRevenue
		err := db.ScanRows(rows, &adSenseRevenue)
		if err != nil {
			fmt.Println(err)
		}
		rowsList = append(rowsList, common.JSON{
			"account_id": adSenseRevenue.AccountId,
			"customer_ad_exchange_estimated_revenue": adSenseRevenue.CustomerAdExchangeEstimatedRevenue,
		})
	}
	fmt.Println(rowsList)
	return rowsList
}
func QueryAdSenseReportList(db *gorm.DB, accountId []string, startDate int, endDate int) []common.JSON{
	result := make(chan []common.JSON)
	wg := sync.WaitGroup{}
	wg.Add(len(accountId))
	for _, id := range accountId {
		go func(result chan<- []common.JSON) {
			defer wg.Done()
			reportList := QueryAdSenseReport(db, id, common.ConvertTime(startDate), common.ConvertTime(endDate))
			result <- reportList
		}(result)
	}

	go func(){
		wg.Wait()
		close(result)
	}()
	var response []common.JSON

	index := 0
	for n := range result {
		for _, s := range n {
			response = append(response, s)
		}
		index = index + 1
	}
	return response
}

func QueryAdSenseRevenueList(db *gorm.DB, accountId []string, startDate int, endDate int) []common.JSON{
	result := make(chan []common.JSON)
	wg := sync.WaitGroup{}
	wg.Add(len(accountId))
	for _, id := range accountId {
		go func(result chan<- []common.JSON) {
			defer wg.Done()
			reportList := QueryAdSenseRevenue(db, id, common.ConvertTime(startDate), common.ConvertTime(endDate))
			result <- reportList
		}(result)
	}

	go func(){
		wg.Wait()
		close(result)
	}()
	var response []common.JSON

	index := 0
	for n := range result {
		for _, s := range n {
			response = append(response, s)
		}
		index = index + 1
	}
	return response
}