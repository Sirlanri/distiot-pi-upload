package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/Sirlanri/distiot-pi-upload/sdk"
	"github.com/d2r2/go-dht"
)

func gethumi() (int, int, error) {
	//引脚号使用的物理号
	temperature, humidity, err :=
		dht.ReadDHTxx(dht.DHT11, 4, false)
	if err != nil {
		return 0, 0, err
	}
	if temperature == -1 {
		return 0, 0, errors.New("采集温湿度失败")
	}
	return int(temperature), int(humidity), nil

}

func getCodeTemp() float64 {
	file, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		fmt.Println("读取core温度失败", err.Error())
		return 0
	}
	strValue := string(file)[:2] + "." + string(file)[2:5]
	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		fmt.Println("转换温度失败 ", err.Error())
		return 0
	}

	return floatValue
}

func main() {

	for {
		time.Sleep(time.Second)
		coreTemp := getCodeTemp()
		fmt.Println("内核温度 ", coreTemp)
		temp, humi, err := gethumi()
		if err != nil {
			continue
		}
		fmt.Printf("潮湿为%d，温度为%d \n", humi, temp)
	}
}
