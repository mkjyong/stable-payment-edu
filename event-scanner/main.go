package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Event ABI for PayrollBasic.ProofOfPayment
const payrollBasicEventABI = `[
  {
    "anonymous": false,
    "inputs": [
      {"indexed": true, "internalType": "address", "name": "employee", "type": "address"},
      {"indexed": false, "internalType": "uint256", "name": "amount", "type": "uint256"}
    ],
    "name": "ProofOfPayment",
    "type": "event"
  }
]`

type modeType string

const (
	modeHTTP modeType = "http" // FilterLogs once
	modeWS   modeType = "ws"   // SubscribeFilterLogs
)

func main() {
	// Flags
	var (
		modeStr         string
		rpcURL          string
		contractAddress string
		fromBlock       uint64
		employeeFilter  string
	)

	flag.StringVar(&modeStr, "mode", "http", "mode: http|ws")
	flag.StringVar(&rpcURL, "rpc", os.Getenv("RPC"), "RPC endpoint (HTTP for http mode, WS for ws mode). Env: RPC")
	flag.StringVar(&contractAddress, "contract", os.Getenv("CONTRACT_ADDRESS"), "Contract address. Env: CONTRACT_ADDRESS")
	flag.Uint64Var(&fromBlock, "from", 0, "From block (http mode only)")
	flag.StringVar(&employeeFilter, "employee", "", "Filter by employee address (optional)")
	flag.Parse()

	if rpcURL == "" || contractAddress == "" {
		log.Fatal("rpc and contract are required (via flags or env)")
	}

	m := modeType(modeStr)
	switch m {
	case modeHTTP:
		if err := runHTTP(rpcURL, contractAddress, fromBlock, employeeFilter); err != nil {
			log.Fatalf("http scan error: %v", err)
		}
	case modeWS:
		if err := runWS(rpcURL, contractAddress, employeeFilter); err != nil {
			log.Fatalf("ws subscribe error: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s", modeStr)
	}
}

func runHTTP(rpcURL, contractHex string, from uint64, employeeHex string) error {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer client.Close()

	parsedABI, err := abi.JSON(strings.NewReader(payrollBasicEventABI))
	if err != nil {
		return fmt.Errorf("parse abi: %w", err)
	}
	event := parsedABI.Events["ProofOfPayment"]

	contractAddr := common.HexToAddress(contractHex)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		Addresses: []common.Address{contractAddr},
		Topics:    [][]common.Hash{{event.ID}},
	}
	if employeeHex != "" {
		emp := common.HexToAddress(employeeHex)
		// topic for indexed address is left-padded 32 bytes
		query.Topics = [][]common.Hash{{event.ID}, {common.BytesToHash(common.LeftPadBytes(emp.Bytes(), 32))}}
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("filter logs: %w", err)
	}

	for _, vLog := range logs {
		if err := handleLog(parsedABI, vLog); err != nil {
			log.Printf("handle log: %v", err)
		}
	}
	return nil
}

func runWS(rpcURL, contractHex, employeeHex string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer client.Close()

	parsedABI, err := abi.JSON(strings.NewReader(payrollBasicEventABI))
	if err != nil {
		return fmt.Errorf("parse abi: %w", err)
	}
	event := parsedABI.Events["ProofOfPayment"]

	contractAddr := common.HexToAddress(contractHex)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    [][]common.Hash{{event.ID}},
	}
	if employeeHex != "" {
		emp := common.HexToAddress(employeeHex)
		query.Topics = [][]common.Hash{{event.ID}, {common.BytesToHash(common.LeftPadBytes(emp.Bytes(), 32))}}
	}

	logsCh := make(chan types.Log, 64)
	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("listening for ProofOfPayment events...")

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case vLog := <-logsCh:
			if err := handleLog(parsedABI, vLog); err != nil {
				log.Printf("handle log: %v", err)
			}
		case <-sig:
			log.Println("shutting down...")
			return nil
		}
	}
}

func handleLog(parsedABI abi.ABI, vLog types.Log) error {
	// topic[0] = event signature, topics[1] = indexed employee
	employee := topicToAddress(vLog.Topics[1])
	// unpack non-indexed fields from data
	var out struct {
		Amount *big.Int
	}
	if err := parsedABI.UnpackIntoInterface(&out, "ProofOfPayment", vLog.Data); err != nil {
		return fmt.Errorf("unpack: %w", err)
	}

	payload := map[string]any{
		"block":    vLog.BlockNumber,
		"txHash":   vLog.TxHash.Hex(),
		"logIndex": vLog.Index,
		"employee": employee.Hex(),
		"amount":   out.Amount.String(),
	}
	b, _ := json.Marshal(payload)
	fmt.Println(string(b))
	return nil
}

func topicToAddress(h common.Hash) common.Address {
	return common.BytesToAddress(h.Bytes()[12:])
}
