package handlers

import (
	"fmt"
	"meteo_des_aeroports/internal/utils"
	"strconv"
	"time"
)

func GetValueOfDataTypeWithRange(iata string, start string, end string, dataType string) (string, error) {
	listProbes, err := utils.HGetAll(fmt.Sprintf("%s:probes:%s", iata, dataType))

	if err != nil {
		return err.Error(), err
	}

	var result string
	result = `{`

	for index, probeId := range listProbes {
		if index%2 == 0 {

			result += fmt.Sprintf("\"%s\":[", probeId)
			listData, err := utils.ZRANGEBYSCOREWITHSCORES(fmt.Sprintf("%s:probe:%s:%s", iata, dataType, probeId), start, end)

			if err != nil {
				return err.Error(), err
			}

			var value string
			for index, data := range listData {
				if index%2 == 0 {
					value = data
				} else {
					result += `{"` + data + `":` + value + `},`
				}
			}

			if len(listData) > 0 {
				result = result[:len(result)-1]
			}

			result += "],"
		}
	}
	if len(listProbes) > 0 {
		result = result[:len(result)-1]
	}
	result += "}"

	return result, nil
}

func GetAverageValueOfTheDay(iata string, date string) (string, error) {
	var t time.Time

	if date == "" {
		t = time.Now().UTC()
	} else {
		i, err := strconv.ParseInt(date, 10, 64)
		if err != nil {
			return err.Error(), err
		}
		t = time.Unix(i, 0).UTC()
	}

	var start time.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	var end time.Time = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())

	result := "{"

	for _, dataType := range utils.DataTypes {
		if value, err := GetAverageValueOfTheDayOfDataType(iata, string(dataType), start, end); err == nil {
			result += `"` + string(dataType) + `":` + value + ","
		} else {
			return err.Error(), err
		}
	}
	if len(utils.DataTypes) > 0 {
		result = result[:len(result)-1]
	}

	result += "}"

	return result, nil
}

func GetAverageValueOfTheDayOfDataType(iata string, dataType string, start time.Time, end time.Time) (string, error) {
	listProbes, err := utils.HGetAll(fmt.Sprintf("%s:probes:%s", iata, dataType))

	fmt.Println(dataType)

	if err != nil {
		return err.Error(), err
	}

	var sum float64 = 0.0
	var counter float64 = 0.0

	for index, probeId := range listProbes {
		if index%2 == 0 {

			listData, err := utils.ZRANGEBYSCORE(fmt.Sprintf("%s:probe:%s:%s", iata, dataType, probeId), strconv.Itoa(int(start.Unix())), strconv.Itoa(int(end.Unix())))

			fmt.Println(listData)

			if err != nil {
				return err.Error(), err
			}

			for _, data := range listData {
				if value, err := strconv.ParseFloat(data, 64); err == nil {
					sum += value
					counter += 1
				}
			}

		}
	}

	fmt.Printf("sum : %f \ncounter : %f", sum, counter)

	return fmt.Sprintf("%f", sum/counter), nil
}
