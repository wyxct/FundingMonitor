package chain

import (
	"fmt"
	"funding-watch/config"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var Client *ethclient.Client

func ChainInit() {
	Rpc_Url := config.C.Chain.RpcUrl

	var err error
	Client, err = ethclient.Dial(Rpc_Url)
	if err != nil {
		log.Fatal("链连接失败:", err)
		panic(err)
	}

	fmt.Println("✅ 链连接成功")

}
