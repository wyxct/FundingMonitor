package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"funding-watch/abi"
	"funding-watch/config"
	"funding-watch/dao"
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

			eventSig := vLog.Topics[0].Hex()

			fundedSig := parsedABI.Events["Funded"].ID.Hex()
			refundedSig := parsedABI.Events["ReFunded"].ID.Hex()
			distributeSig := parsedABI.Events["Distribute"].ID.Hex()

			switch eventSig {
			case fundedSig:
				fmt.Println("✅ 收到 Funded 事件")
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

				record := models.NewFundRecord(sender.Hex(), amount.String(), vLog.TxHash.Hex(), int64(vLog.BlockNumber))
				FundingRecordDao := dao.NewFundingRecordDao()
				err = FundingRecordDao.InsertRecord(record)
				if err != nil {
					fmt.Println("写入数据库错误:", err)
					continue
				}
				fundingTotalDao := dao.NewFundTotalRecordDao()
				fundingTotal, err := fundingTotalDao.GetRecordBySender(sender.Hex())
				fmt.Println(fundingTotal)
				if len(fundingTotal) == 0 {
					fundingTotalRecord := models.NewFundingTotalRecord(sender.Hex(), amount.String())
					err = fundingTotalDao.Insert(fundingTotalRecord)
				} else {
					fundingTotalRecord := fundingTotal[0]
					oldTotal, _ := new(big.Float).SetString(fundingTotalRecord.FundingTotal)
					newTotal := new(big.Float).Add(oldTotal, amount)
					fundingTotalRecord.FundingTotal = newTotal.String()
					err = fundingTotalDao.Update(&fundingTotalRecord)
				}
				fmt.Println("✅ 已写入 PostgreSQL")
			case refundedSig:
				fmt.Println("✅ 收到 Refunded 事件")
				event := struct {
					Amount    *big.Int
					Timestamp *big.Int
				}{}
				sender := common.HexToAddress(vLog.Topics[1].Hex())
				err := parsedABI.UnpackIntoInterface(&event, "ReFunded", vLog.Data)
				if err != nil {
					fmt.Println("解析错误:", err)
					continue
				}

				// 转 ETH
				amount := new(big.Float).Quo(
					new(big.Float).SetInt(event.Amount),
					big.NewFloat(1e18),
				)

				record := models.NewReFundRecord(sender.Hex(), amount.String(), vLog.TxHash.Hex(), int64(vLog.BlockNumber))
				ReFundingRecordDao := dao.NewReFundingRecordDao()
				err = ReFundingRecordDao.InsertRecord(record)
				if err != nil {
					fmt.Println("写入数据库错误:", err)
					continue
				}
				fundingTotalDao := dao.NewFundTotalRecordDao()
				fundingTotal, err := fundingTotalDao.GetRecordBySender(sender.Hex())
				fmt.Println(fundingTotal)
				if len(fundingTotal) == 0 {
					fmt.Println("❌ 无Funded记录，无法计算总资金")
				} else {
					fundingTotalRecord := fundingTotal[0]
					oldTotal, _ := new(big.Float).SetString(fundingTotalRecord.FundingTotal)
					newTotal := new(big.Float).Sub(oldTotal, amount)
					fundingTotalRecord.FundingTotal = newTotal.String()
					err = fundingTotalDao.Update(&fundingTotalRecord)
				}
				fmt.Println("✅ 已写入 PostgreSQL")
			case distributeSig:
				fmt.Println("✅ 收到 Distribute 事件")
				event := struct {
					Amount    *big.Int
					Timestamp *big.Int
				}{}
				DistributeAddr := common.HexToAddress(vLog.Topics[1].Hex())
				err := parsedABI.UnpackIntoInterface(&event, "Distribute", vLog.Data)
				if err != nil {
					fmt.Println("解析错误:", err)
					continue
				}

				// 转 ETH
				amount := new(big.Float).Quo(
					new(big.Float).SetInt(event.Amount),
					big.NewFloat(1e18),
				)

				record := models.NewDistributeRecord(DistributeAddr.Hex(), amount.String(), vLog.TxHash.Hex(), int64(vLog.BlockNumber))
				DistributeRecordDao := dao.NewDistributeRecordDao()
				err = DistributeRecordDao.InsertRecord(record)
				if err != nil {
					fmt.Println("写入数据库错误:", err)
					continue
				}
				fmt.Println("✅ 已写入 PostgreSQL")
			default:
				fmt.Println("❌ 未知事件类型")
				continue
			}

			lastBlock = currentBlock
		}
	}
}
