// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"strings"
// 	"time"

// 	"funding-watch/models"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// // 合约ABI
// const abiJSON = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"timestamp","type":"uint256"}],"name":"Funded","type":"event"}]`

// func main() {
// 	// ===================== 你的配置 =====================
// 	rpcURL := "http://127.0.0.1:8545"
// 	contractAddr := "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
// 	pgDSN := "host=localhost user=postgres password=123 dbname=Funding port=5432 sslmode=disable TimeZone=UTC"
// 	// ====================================================

// 	FundRecord := models.FundRecord{}

// 	// 1. 连接DB
// 	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("DB连接失败:", err)
// 	}
// 	fmt.Println("✅ PostgreSQL 连接成功")

// 	// 2. 连接链
// 	client, err := ethclient.Dial(rpcURL)
// 	if err != nil {
// 		log.Fatal("链连接失败:", err)
// 	}
// 	fmt.Println("✅ 链连接成功")

// 	// 3. 解析ABI
// 	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
// 	if err != nil {
// 		log.Fatal("ABI解析错误:", err)
// 	}

// 	contractAddress := common.HexToAddress(contractAddr)

// 	// 4. 从当前最新块开始监听
// 	lastBlock, err := client.BlockNumber(context.Background())
// 	if err != nil {
// 		log.Fatal("获取块高失败:", err)
// 	}
// 	fmt.Printf("✅ 开始监听 Funded 事件，从块高 %d 开始\n", lastBlock)

// 	// 5. 循环轮询（兼容 Hardhat 本地节点）
// 	for {
// 		time.Sleep(1 * time.Second)

// 		// 当前最新块
// 		currentBlock, err := client.BlockNumber(context.Background())
// 		if err != nil {
// 			fmt.Println("获取最新块错误:", err)
// 			continue
// 		}

// 		if currentBlock <= lastBlock {
// 			continue
// 		}

// 		// 查询区间日志
// 		query := ethereum.FilterQuery{
// 			Addresses: []common.Address{contractAddress},
// 			FromBlock: big.NewInt(int64(lastBlock + 1)),
// 			ToBlock:   big.NewInt(int64(currentBlock)),
// 		}

// 		logs, err := client.FilterLogs(context.Background(), query)
// 		if err != nil {
// 			fmt.Println("获取日志错误:", err)
// 			lastBlock = currentBlock
// 			continue
// 		}

// 		// 处理每条日志
// 		for _, vLog := range logs {
// 			fmt.Printf("\n📥 收到交易：%s\n", vLog.TxHash.Hex())

// 			event := struct {
// 				Amount    *big.Int
// 				Timestamp *big.Int
// 			}{}
// 			sender := common.HexToAddress(vLog.Topics[1].Hex())
// 			err := parsedABI.UnpackIntoInterface(&event, "Funded", vLog.Data)
// 			if err != nil {
// 				fmt.Println("解析错误:", err)
// 				continue
// 			}

// 			// 转 ETH
// 			amount := new(big.Float).Quo(
// 				new(big.Float).SetInt(event.Amount),
// 				big.NewFloat(1e18),
// 			)

// 			// 写入数据库
// 			record := FundRecord{
// 				Sender:   sender.Hex(),
// 				Amount:   amount.String(),
// 				TxHash:   vLog.TxHash.Hex(),
// 				BlockNum: int64(vLog.BlockNumber),
// 			}
// 			db.Create(&record)
// 			fmt.Println("✅ 已写入 PostgreSQL")
// 		}

// 		lastBlock = currentBlock
// 	}
// }

package main

import (
	"fmt"
	"funding-watch/abi"
	"funding-watch/chain"
	"funding-watch/config"
	"funding-watch/routes"
	"funding-watch/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.InitDB()
	chain.ChainInit()
	abi.InitABI()
	go service.StartWatch()
	r := gin.Default()
	routes.InitRoutes(r)
	r.Run(fmt.Sprintf(":%s", config.C.App.Port))
}
