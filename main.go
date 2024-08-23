/*
what to do:
exchange ['5j6Ak1GJaTy8XoC', '56kC8jyGoXTAa1J']
reverse
rc4 hAGMmLFnoa0
exchange = ['PUoVzgdK5FLZt', 'FVogUPtKzdZL5']
reverse
rc4 oUHxby23izOI5
exchange ['PEQmieNvWhrOX', 'OEehvmXQrWiPN']
reverse
rc4 tX6D4K8mPrq3V
base64
*/
package main

import (
	"cinezonescraper/api"
	"cinezonescraper/utils"
)

func main()  {
	utils.SetupKeys()
	go api.StartServer()
	select {}
}