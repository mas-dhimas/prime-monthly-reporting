package handler

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	network "gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/service"
)

type Handler struct {
	service *network.Service
}

func NewHandler(service *network.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RouteService() error {
r:
	for {
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("|    Current available modules: 'availability'                   |")
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("|    Please type to generate report data following below format  |")
		fmt.Println("|    eg: module=availability start-month=1 end-month=2           |")
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("|    [exit] or [quit] or [ctrl + c] to stop program              |")
		fmt.Println("------------------------------------------------------------------")

		reader := bufio.NewReader(os.Stdin)
		data := make(map[string]string)

		fmt.Print("|    Input: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		fmt.Println()
		switch text {
		case "exit", "quit":
			logrus.Info("Stopping program....")
			break r
		default:
			parts := strings.Fields(text)

			for _, part := range parts {
				kv := strings.SplitN(part, "=", 2)
				if len(kv) != 2 {
					fmt.Printf("invalid segment: %s (use key=value)\n", part)
					continue
				}

				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])

				data[key] = value
			}
		}

		module := data["module"]
		start := data["start-month"]
		startMonth, err := strconv.Atoi(start)
		if err != nil {
			logrus.Warn("invalid start month period, using January as default start month")
			startMonth = 1
		}

		end := data["end-month"]
		endMonth, err := strconv.Atoi(end)
		if err != nil {
			logrus.Warn("invalid end month period, using current month as default end month")
			fmt.Println()
			endMonth = int(time.Now().Month())
		}
		flag.Parse()

		filename := ""
		switch strings.ToLower(module) {
		case "availability":
			fmt.Printf("|    Generating report for module %s....\n", module)
			filename, err = h.service.GenerateDeviceAvailabilityReporting(startMonth, endMonth)
			if err != nil {
				return err
			}
		default:
			logrus.Warn(fmt.Sprintf("specied module not implemented yet, please request to %s", "DEV NETMONK"))
			continue
		}

		fmt.Printf("|    Success generate report, please check %s\n", filename)
		fmt.Print("\n\n")
		continue
	}

	return nil
}
