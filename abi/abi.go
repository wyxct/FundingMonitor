package abi

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// 全局 ABI
var FundingABI abi.ABI

// 定义结构体，用来解析 Hardhat artifact
type HardhatArtifact struct {
	ABI json.RawMessage `json:"abi"` // 只提取 abi 字段
}

func InitABI() {
	// 1. 读取文件
	path := "./abi/Funding.json"
	data, err := os.ReadFile(path)
	if err != nil {
		panic("读取ABI文件失败: " + err.Error())
	}

	// 2. 解析外层 JSON，只拿 abi 字段
	var artifact HardhatArtifact
	err = json.Unmarshal(data, &artifact)
	if err != nil {
		panic("解析Hardhat artifact失败: " + err.Error())
	}

	// 3. 解析真正的 ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(artifact.ABI)))
	if err != nil {
		panic("解析ABI失败: " + err.Error())
	}

	FundingABI = parsedABI
	println("✅ ABI 初始化成功")
}
