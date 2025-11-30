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
	xlsx.SetCellValue(sheet1Name, "A1", "IP Address")
	xlsx.SetCellValue(sheet1Name, "B1", "Month")
	xlsx.SetCellValue(sheet1Name, "C1", "Percentage")
	xlsx.SetCellValue(sheet1Name, "D1", "Uptime")
	xlsx.SetCellValue(sheet1Name, "E1", "Downtime")

	row := 0
	for _, v := range data {
		if v.Data.IP == "" {
			continue
		}
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", row+2), v.Data.IP)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", row+2), v.Data.Month)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", row+2), v.Data.IcmpPing.Ratio)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", row+2), v.Data.IcmpPing.Uptime)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", row+2), v.Data.IcmpPing.Downtime)
		row++
	}

	filename := fmt.Sprintf("%s-%s.xlsx", "device_availability", time.Now().Format(time.RFC822))
	err = xlsx.SaveAs(fmt.Sprintf("./docs/%s", filename))
	if err != nil {
		fmt.Println(err)
	}

	return filename, nil
}
