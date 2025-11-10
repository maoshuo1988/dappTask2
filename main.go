package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Sepolia测试网络节点URL - 需要替换为您的实际项目ID
	sepoliaURL = "https://sepolia.infura.io/v3/33045a1d50ff4aa2ba2321a76281a2ed"

	// 示例私钥 - 在实际使用中请替换为您的私钥，并确保安全存储
	privateKeyHex = "92a23b579cfc51fb9579dfb9aadd9ceb03832216f2e96ede1cbf623ab20e3778"
)

func main() {
	// 连接到 Sepolia 测试网络
	client, err := ethclient.Dial(sepoliaURL)
	if err != nil {
		log.Fatal(err)
	}

	// 加载私钥（测试网专用）
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	_, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 创建认证对象
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 合约地址（部署后获取）
	contractAddress := common.HexToAddress("0x759F4d7D6Bbc4b11b418bcb59d23ba603f5e5C89")

	// 创建合约实例 - 修复 NewCounter 未定义问题
	instance, err := NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 读取当前计数器值
	currentCount, err := instance.Count(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("当前计数器值: %d\n", currentCount)

	// 调用增加计数器方法
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == 1 {
		fmt.Println("计数器增加成功！")

		// 再次读取更新后的值
		newCount, err := instance.Count(&bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("新的计数器值: %d\n", newCount)
	} else {
		fmt.Println("交易失败")
	}
}
