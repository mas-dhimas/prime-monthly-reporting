package network

import (
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/lib"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/models"
)

func (s *Service) GenerateDeviceAvailabilityReporting(from, to int) (string, error) {
	data := []models.DeviceAvailabilityReport{}

	year := time.Now().Year()
	layout := "02/01/2006 MST"

	res, err := s.repo.GetDeviceInventory()
	if err != nil {
		logrus.Error("unable to get device inventory : ", err)
		return "", err
	}

	logrus.Info("success get device inventory")

	// // loop every inventory devices
	for i, v := range res.Data.Nodes {
		if i != 0 {
			continue
		}
		// 	// loop each months
		for i := from; i <= to; i++ {
			month := strconv.Itoa(i)
			if i < 10 {
				month = fmt.Sprintf("0%s", month)
			}
			value := fmt.Sprintf("01/%s/%s WIB", month, strconv.Itoa(year))

			start, err := time.Parse(layout, value)
			if err != nil {
				return "", err
			}

			startUnix := start.Unix()
			endUnix := lib.EndOfMonth(start).Unix()

			availability, err := s.repo.GetDeviceAvailibilityReporting(v.ID.Hex(), int(startUnix), int(endUnix))
			if err != nil {
				return "", err
			}

			availability.Data.Month = start.Month().String()

			data = append(data, *availability)
		}
	}

	xlsx := excelize.NewFile()
	sheet1Name := "Sheet One"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	// ip mgmt
	xlsx.MergeCell(sheet1Name, "A1", "A3")
	xlsx.SetCellValue(sheet1Name, "A1", "IP Address")

	// month
	xlsx.MergeCell(sheet1Name, "B1", "B3")
	xlsx.SetCellValue(sheet1Name, "B1", "Month")

	// availability table name
	xlsx.MergeCell(sheet1Name, "C1", "J1")
	// xlsx.MergeCell(sheet1Name, "C1", "F1")
	xlsx.SetCellValue(sheet1Name, "C1", "Availability")

	// ICMP PING table name
	xlsx.MergeCell(sheet1Name, "C2", "F2")
	xlsx.SetCellValue(sheet1Name, "C2", "ICMP / Ping")

	// ICMP Availability utilization
	xlsx.SetCellValue(sheet1Name, "C3", "Percentage (%)")
	xlsx.SetCellValue(sheet1Name, "D3", "Uptime")
	xlsx.SetCellValue(sheet1Name, "E3", "Downtime")
	xlsx.SetCellValue(sheet1Name, "F3", "Undetected")

	// SNMP Uptime table name
	xlsx.MergeCell(sheet1Name, "G2", "J2")
	xlsx.SetCellValue(sheet1Name, "G2", "SNMP Uptime")

	// SNMP Uptime utilization
	xlsx.SetCellValue(sheet1Name, "G3", "Percentage (%)")
	xlsx.SetCellValue(sheet1Name, "H3", "Uptime")
	xlsx.SetCellValue(sheet1Name, "I3", "Downtime")
	xlsx.SetCellValue(sheet1Name, "J3", "Undetected")

	row := 0
	lastPoint := ""
	for i, v := range data {
		startpoint := 4
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", row+startpoint), v.Data.IP)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", row+startpoint), v.Data.Month)
		// ICMP / PING
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", row+startpoint), v.Data.IcmpPing.Ratio)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", row+startpoint), v.Data.IcmpPing.Uptime)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", row+startpoint), v.Data.IcmpPing.Downtime)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", row+startpoint), v.Data.IcmpPing.UnknownTime)

		// SNMP Uptime
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", row+startpoint), v.Data.SnmpUptime.Ratio)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", row+startpoint), v.Data.SnmpUptime.Uptime)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", row+startpoint), v.Data.SnmpUptime.Downtime)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", row+startpoint), v.Data.SnmpUptime.UnknownTime)
		if i == len(data)-1 {
			lastPoint = fmt.Sprintf("J%d", row+startpoint)
		}
		row++
	}

	style, err := xlsx.NewStyle(`
	{
		"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false},
		"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]
	}`)
	if err != nil {
		logrus.Error("error set broder style : ", err)
	}
	xlsx.SetCellStyle(sheet1Name, "A1", lastPoint, style)

	filename := fmt.Sprintf("%s-%s.xlsx", "device_availability", time.Now().Format(time.RFC822))
	err = xlsx.SaveAs(fmt.Sprintf("./docs/%s", filename))
	if err != nil {
		fmt.Println(err)
	}

	return filename, nil
}
