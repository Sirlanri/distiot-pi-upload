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
	man := sdk.NewManager("68b47b8e-8dcc-4f09-97eb-801123c58f59")
	man.MasterUrl = "http://192.168.1.150:8001/master"
	man.UserUrl = "http://192.168.1.150:8091/user"

	humiDevice, err := man.NewDevice(17)
	if err != nil {
		fmt.Println("创建humi设备失败 ", err.Error())
		return
	}
	tempDevice, err := man.NewDevice(16)
	if err != nil {
		fmt.Println("创建temp设备失败 ", err.Error())
		return
	}

	coreDevice, err := man.NewDevice(18)
	if err != nil {
		fmt.Println("创建CoreTemp设备失败", err.Error())
		return
	}
	for {
		time.Sleep(time.Second)
		coreTemp := getCodeTemp()
		fmt.Println("内核温度 ", coreTemp)
		err = coreDevice.UploadDataHttp(strconv.FormatFloat(coreTemp, 'E', -1, 64))
		if err != nil {
			fmt.Println("CoreTemp 上传温度失败！", err.Error())
			return
		}
		//获取上传温湿度数据
		temp, humi, err := gethumi()
		if err != nil {
			continue
		}
		err = tempDevice.UploadDataHttp(strconv.Itoa(temp))
		if err != nil {
			fmt.Println("上传temp失败 ", err.Error())
			return
		}
		err = humiDevice.UploadDataHttp(strconv.Itoa(humi))
		if err != nil {
			fmt.Println("上传humi失败 ", err.Error())
			return
		}
	}
}
