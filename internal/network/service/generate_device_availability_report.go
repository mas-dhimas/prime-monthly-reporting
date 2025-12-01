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
	for _, v := range res.Data.Nodes {
		// 	// loop each months
		for i := from; i <= to; i++ {
			if i < int(v.CreatedAt.Month()) && int(v.CreatedAt.Year()) == year {
				continue
			}

			month := strconv.Itoa(i)
			if i < 10 {
				month = fmt.Sprintf("0%s", month)
			}
			value := fmt.Sprintf("01/%s/%s WIB", month, strconv.Itoa(year))

			start, err := time.Parse(layout, value)
			if err != nil {
				return "", err
			}

			startUnix := start.Add(-7 * time.Hour).Unix()
			endUnix := lib.EndOfMonth(start).Add(-7 * time.Hour).Add(-1 * time.Second).Unix()

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

	// node name
	xlsx.MergeCell(sheet1Name, "B1", "B3")
	xlsx.SetCellValue(sheet1Name, "B1", "Node Name")

	// month
	xlsx.MergeCell(sheet1Name, "C1", "C3")
	xlsx.SetCellValue(sheet1Name, "C1", "Month")

	// availability table name
	xlsx.MergeCell(sheet1Name, "D1", "K1")
	// xlsx.MergeCell(sheet1Name, "C1", "F1")
	xlsx.SetCellValue(sheet1Name, "D1", "Availability")

	// ICMP PING table name
	xlsx.MergeCell(sheet1Name, "D2", "G2")
	xlsx.SetCellValue(sheet1Name, "D2", "ICMP / Ping")

	// ICMP Availability utilization
	xlsx.SetCellValue(sheet1Name, "D3", "Percentage (%)")
	xlsx.SetCellValue(sheet1Name, "E3", "Uptime")
	xlsx.SetCellValue(sheet1Name, "F3", "Downtime")
	xlsx.SetCellValue(sheet1Name, "G3", "Undetected")

	// SNMP Uptime table name
	xlsx.MergeCell(sheet1Name, "H2", "K2")
	xlsx.SetCellValue(sheet1Name, "H2", "SNMP Uptime")

	// SNMP Uptime utilization
	xlsx.SetCellValue(sheet1Name, "H3", "Percentage (%)")
	xlsx.SetCellValue(sheet1Name, "I3", "Uptime")
	xlsx.SetCellValue(sheet1Name, "J3", "Downtime")
	xlsx.SetCellValue(sheet1Name, "K3", "Undetected")

	row := 0
	lastPoint := ""
	for i, v := range data {
		startpoint := 4
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", row+startpoint), v.Data.IP)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", row+startpoint), v.Data.NodeName)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", row+startpoint), v.Data.Month)
		// ICMP / PING
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", row+startpoint), fmt.Sprintf("%.2f", v.Data.IcmpPing.Ratio))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", row+startpoint), lib.FormatSeconds(v.Data.IcmpPing.Uptime))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", row+startpoint), lib.FormatSeconds(v.Data.IcmpPing.Downtime))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", row+startpoint), lib.FormatSeconds(v.Data.IcmpPing.UnknownTime))

		// SNMP Uptime
		ratio := fmt.Sprintf("%.2f", v.Data.SnmpUptime.Ratio)
		if v.Data.SnmpUptime.Ratio < 0 {
			ratio = "N/A"
		}

		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", row+startpoint), ratio)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", row+startpoint), lib.FormatSeconds(v.Data.SnmpUptime.Uptime))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", row+startpoint), lib.FormatSeconds(v.Data.SnmpUptime.Downtime))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", row+startpoint), lib.FormatSeconds(v.Data.SnmpUptime.UnknownTime))
		if i == len(data)-1 {
			lastPoint = fmt.Sprintf("K%d", row+startpoint)
		}
		row++
	}

	border, err := xlsx.NewStyle(`
	{
		"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":false,"vertical":"","wrap_text":false},
		"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]
	}`)
	if err != nil {
		logrus.Error("error set broder style : ", err)
	}
	xlsx.SetCellStyle(sheet1Name, "A1", lastPoint, border)

	filename := fmt.Sprintf("%s-%s.xlsx", "device_availability", time.Now().Format(time.RFC822))
	err = xlsx.SaveAs(fmt.Sprintf("./docs/%s", filename))
	if err != nil {
		fmt.Println(err)
	}

	return filename, nil
}
