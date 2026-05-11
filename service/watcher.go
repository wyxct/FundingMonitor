package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"funding-watch/abi"
	"funding-watch/config"
	"funding-watch/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartWatch() {
	// ===================== 你的配置 =====================
	rpcURL := config.C.Chain.RpcUrl
	contractAddr := config.C.Chain.Contract
	// ====================================================

	// 2. 连接链
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("链连接失败:", err)
	}
	fmt.Println("✅ 链连接成功")

	// 3. 解析ABI
	parsedABI := abi.FundingABI

	contractAddress := common.HexToAddress(contractAddr)
	fmt.Printf("✅ 合约地址为：%s\n", contractAddress.Hex())
	// 4. 从当前最新块开始监听
	lastBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("获取块高失败:", err)
	}
	fmt.Printf("✅ 开始监听 Funded 事件，从块高 %d 开始\n", lastBlock)

	// 5. 循环轮询（兼容 Hardhat 本地节点）
	for {
		time.Sleep(1 * time.Second)

		// 当前最新块
		currentBlock, err := client.BlockNumber(context.Background())
		if err != nil {
			fmt.Println("获取最新块错误:", err)
			continue
		}

		if currentBlock <= lastBlock {
			continue
		}

		// 查询区间日志
		query := ethereum.FilterQuery{
			Addresses: []common.Address{contractAddress},
			FromBlock: big.NewInt(int64(lastBlock + 1)),
			ToBlock:   big.NewInt(int64(currentBlock)),
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			fmt.Println("获取日志错误:", err)
			lastBlock = currentBlock
			continue
		}

		// 处理每条日志
		for _, vLog := range logs {
			fmt.Printf("\n📥 收到交易：%s\n", vLog.TxHash.Hex())

			event := struct {
				Amount    *big.Int
				Timestamp *big.Int
			}{}
			sender := common.HexToAddress(vLog.Topics[1].Hex())
			err := parsedABI.UnpackIntoInterface(&event, "Funded", vLog.Data)
			if err != nil {
				fmt.Println("解析错误:", err)
				continue
			}

			// 转 ETH
			amount := new(big.Float).Quo(
				new(big.Float).SetInt(event.Amount),
				big.NewFloat(1e18),
			)

			// 写入数据库
			// record := FundRecord{
			// 	Sender:   sender.Hex(),
			// 	Amount:   amount.String(),
			// 	TxHash:   vLog.TxHash.Hex(),
			// 	BlockNum: int64(vLog.BlockNumber),
			// }
			// db.Create(&record)
			// fmt.Println("✅ 已写入 PostgreSQL")

			record := models.NewFundRecord(sender.Hex(), amount.String(), vLog.TxHash.Hex(), int64(vLog.BlockNumber))
			err = record.Create()
			if err != nil {
				fmt.Println("写入数据库错误:", err)
				continue
			}
			fmt.Println("✅ 已写入 PostgreSQL")
		}

		lastBlock = currentBlock
	}
}
