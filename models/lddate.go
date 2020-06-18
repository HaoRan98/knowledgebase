package models

import (
	"log"
	"strconv"
	"time"
)

//基准时间
type LdDate struct {
	LdCurrentDate string //当前日期
	LdCurrentYear string //当前年度
	LdHzrq        string //汇总日期,为当前日期前一天
	LdNext        string //当前日期的下一天
	LdHzrqSntq    string //汇总日期(上年同期),为上年同期的当前日期前一天的日期

	LdMonthFirst         string //本月第一天
	LdMonthEnd           string //本月最后一天
	LdNextMonthFirst     string //下月第一天
	LdNextMonthEnd       string //下月最后一天
	LdLastMonthFirst     string //上月第一天
	LdLastMonthEnd       string //上月最后一天
	LdMonthFirstSntq     string //上年同期本月第一天
	LdMonthEndSntq       string //上年同期本月最后一天
	LdLastMonthFirstSntq string //上年同期上月第一天
	LdLastMonthEndSntq   string //上年同期上月最后一天
	LdLastTwoMonthFirst  string //上上月第一天
	LdLastTwoMonthEnd    string //上上月最后一天

	LdLastMonth            string //上月月份
	LdLastMonthYear        string //上月年度
	LdQrqBn                string //起日期,为本年1月1日
	LdZrqBn                string //止日期,为本年12月31日
	LdQrqLastMonthYear     string //上月年度起日期,为上月年度1月1日
	LdZrqLastMonthYear     string //上月年度止日期,为上月年度12月31日
	LdQrqLastMonthYearSntq string //上月年度上年同期起日期,为上月年度1月1日
	LdZrqLastMonthYearSntq string //上月年度上年同期止日期,为上月年度12月31日
}

//传入参数：20XX-XX-XX
func NewLdDate(date string) *LdDate {
	var ldToday string
	if date == "" { //日期参数,默认为null，取当前系统时间
		ldToday = time.Now().Format("2006-01-02")
	} else {
		ldToday = date
	}
	ldHzrq := GetYesterDay(ldToday)
	ldNextDay := GetTomorrow(ldToday)
	ldTime, err := time.Parse("2006-01-02", ldToday)
	if err != nil {
		log.Println("Parse err==>", err)
		return nil
	}

	ldLastMonthToday := ldTime.AddDate(0, -1, 0).Format("2006-01-02")
	ldNextMonthToday := ldTime.AddDate(0, 1, 0).Format("2006-01-02")
	ldLastTwoMonthToday := ldTime.AddDate(0, -2, 0).Format("2006-01-02")
	//上年同期的今天
	ldTodaySntq := ldTime.AddDate(-1, 0, 0).Format("2006-01-02")
	//上年同期的上月今天
	ldLastMonthTodaySntq := ldTime.AddDate(-1, -1, 0).Format("2006-01-02")

	return &LdDate{
		LdCurrentDate: ldToday,
		LdCurrentYear: ldToday[:4],
		LdMonthFirst:  GetFirstDayOfMonth(ldToday),
		LdMonthEnd:    GetLastDayOfMonth(ldToday),
		LdQrqBn:       GetFirstDayOfYear(ldToday),
		LdZrqBn:       GetLastDayOfYear(ldToday),
		LdHzrq:        ldHzrq,
		LdNext:        ldNextDay,
		LdHzrqSntq:    GetDateOfLastYear(ldHzrq),

		LdLastMonth:            ldLastMonthToday[5:7],
		LdLastMonthYear:        ldLastMonthToday[:4],
		LdQrqLastMonthYear:     GetFirstDayOfYear(ldLastMonthToday),
		LdZrqLastMonthYear:     GetLastDayOfYear(ldLastMonthToday),
		LdQrqLastMonthYearSntq: GetFirstDayOfYear(ldLastMonthTodaySntq),
		LdZrqLastMonthYearSntq: GetLastDayOfYear(ldLastMonthTodaySntq),

		LdNextMonthFirst:     GetFirstDayOfMonth(ldNextMonthToday),
		LdNextMonthEnd:       GetLastDayOfMonth(ldNextMonthToday),
		LdLastMonthFirst:     GetFirstDayOfMonth(ldLastMonthToday),
		LdLastMonthEnd:       GetLastDayOfMonth(ldLastMonthToday),
		LdMonthFirstSntq:     GetFirstDayOfMonth(ldTodaySntq),
		LdMonthEndSntq:       GetLastDayOfMonth(ldTodaySntq),
		LdLastMonthFirstSntq: GetFirstDayOfMonth(ldLastMonthTodaySntq),
		LdLastMonthEndSntq:   GetLastDayOfMonth(ldLastMonthTodaySntq),
		LdLastTwoMonthFirst:  GetFirstDayOfMonth(ldLastTwoMonthToday),
		LdLastTwoMonthEnd:    GetLastDayOfMonth(ldLastTwoMonthToday),
	}
}

//获取昨天的日期
func GetYesterDay(today string) string {
	nTime, _ := time.Parse("2006-01-02", today)
	yesterday := nTime.AddDate(0, 0, -1).Format("2006-01-02")
	return yesterday
}

//获取明天的日期
func GetTomorrow(today string) string {
	nTime, _ := time.Parse("2006-01-02", today)
	tomorrow := nTime.AddDate(0, 0, 1).Format("2006-01-02")
	return tomorrow
}

//获取月份第一天
func GetFirstDayOfMonth(today string) string {
	nTime, _ := time.Parse("2006-01-02", today)
	currentYear, currentMonth, _ := nTime.Date()
	currentLocation := nTime.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	return firstOfMonth.Format("2006-01-02")
}

//获取月份最后一天
func GetLastDayOfMonth(today string) string {
	nTime, _ := time.Parse("2006-01-02", today)
	currentYear, currentMonth, _ := nTime.Date()
	currentLocation := nTime.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Format("2006-01-02")
	return lastOfMonth
}

//获取年度第一天
func GetFirstDayOfYear(today string) string {
	return today[:4] + "-01-01"
}

//获取年度最后一天
func GetLastDayOfYear(today string) string {
	return today[:4] + "-12-31"
}

//获取上年同期
func GetDateOfLastYear(today string) string {
	year, _ := strconv.Atoi(today[:4])
	lastYear := year - 1
	return strconv.Itoa(lastYear) + today[4:]
}
