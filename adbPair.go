package main

import (
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/mdp/qrterminal"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)
var letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$"
func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
var serviceId string
var password string

func main() {
	fmt.Println("🚀 如果设备已配对 可能已经自动连接")

	fmt.Println("🚀 若未配对,请使用Android 11 以上设备在开发者设置中的无线调试选项中`二维码配对`选项扫描如下二维码")
	serviceId=fmt.Sprintf("studio-%s",RandStringRunes(8))
	password=RandStringRunes(8)
	qrterminal.GenerateHalfBlock(fmt.Sprintf("WIFI:T:ADB;S:%s;P:%s;;",serviceId,password), qrterminal.M, os.Stdout)
	cmd := exec.Command("adb","mdns", "check")
	var res []byte
	var err error

	if res, err = cmd.CombinedOutput(); err != nil {
		if strings.Contains(string(res),"unknown command"){
			fmt.Println(string(res))
		}
	//	fmt.Println(err)
	//	if(err.)
	//	os.Exit(1)
	}
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		fmt.Println("❌ Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			//log.Println(entry)
			if entry.ServiceRecord.Instance==serviceId {
				fmt.Println("🚀 找到一个设备 正在配对....")


				cmd := exec.Command("adb","pair", fmt.Sprintf("%s:%d",entry.AddrIPv4,entry.Port),password)
				fmt.Println(fmt.Sprintf("%s:%d",entry.AddrIPv4,entry.Port),password)
				var res []byte
				var err error

				if res, err = cmd.CombinedOutput(); err != nil {
					//if strings.Contains(string(res),"unknown command"){
					//}
						fmt.Println(string(res))

			}

			if strings.Contains(string(res),"Failed:") {
				fmt.Println("❌ 配对失败:"+string(res))
				os.Exit(1)
			}
				fmt.Println(string(res))

				if strings.Contains(string(res),"Successfully") {
					fmt.Println("🚀 配对成功 "+string(res))
					os.Exit(0)
				}
			}

		}
		fmt.Println("No more entries.")
	}(entries)
	err = resolver.Browse(ctx, "_adb-tls-pairing._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("查找设备失败", err.Error())
	}

	<-ctx.Done()

	fmt.Println(string(res))

}
