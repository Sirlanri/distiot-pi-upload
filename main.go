package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/d2r2/go-dht"
)

func gethumi() {
	//引脚号使用的物理号
	temperature, humidity, err :=
		dht.ReadDHTxx(dht.DHT11, 4, false)
	if err != nil {
		fmt.Println(err.Error())
	}
	if temperature != -1 {
		// Print temperature and humidity
		fmt.Printf("Temperature = %v*C, Humidity = %v%%\n",
			temperature, humidity)
	}

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
	coreTemp := getCodeTemp()
	fmt.Println(coreTemp)
}
