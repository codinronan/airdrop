package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	algodHost  = flag.String("airdrop-algod-host", "127.0.0.1", "hostname of Algorand node")
	algodPort  = flag.Int("airdrop-algod-port", 4001, "port of Algorand node")
	algodToken = flag.String("airdrop-algod-token", "", "token to use when making requests to Algorand node")

	kmdHost  = flag.String("airdrop-kmd-host", "127.0.0.1", "hostname of Key Management Daemon")
	kmdPort  = flag.Int("airdrop-kmd-port", 4001, "port of Key Management Daemon")
	kmdToken = flag.String("airdrop-kmd-token", "", "token to use when making requests to Key Management Daemon")

	asset = flag.Int("airdrop-asset-id", 0, "Asset id to transfer")
	amt   = flag.Int("airdrop-asset-amt", 0, "Asset id to transfer")
	fname = flag.String("airdrop-asset-file", "address.csv", "Path to file containing list of addresses to airdrop to")

	sender = flag.String("airdrop-from", "", "Address to use as sender")

	wallet = flag.String("airdrop-wallet", "", "Name of wallet to use")
)

func main() {
	flag.Parse()

	// Create a kmd client
	kmdAddress := fmt.Sprintf("http://%s:%d", *kmdHost, *kmdPort)
	kmdClient, err := kmd.MakeClient(kmdAddress, *kmdToken)
	if err != nil {
		log.Fatalf("failed to make kmd client: %s\n", err)
	}

	// Create an algod client
	algodAddress := fmt.Sprintf("http://%s:%d", *algodHost, *algodPort)
	algodClient, err := algod.MakeClient(algodAddress, *algodToken)
	if err != nil {
		log.Fatalf("Failed to make client: %+v", err)
	}

	// Read in file
	addrs := getAddrs(*fname)

	log.Printf("Got %d addresses from %s", len(addrs), *fname)

	// Get pub key from string
	from, err := types.DecodeAddress(*sender)
	if err != nil {
		log.Fatalf("Invalid address: %s", *sender)
	}
	pk := ed25519.PublicKey(from[:])

	// Get PW from user
	fmt.Println("Please enter the wallet password:")
	walletpw, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Failed to read password: %+v", err)
	}

	// Get suggested params from algod
	suggestedParams, err := algodClient.SuggestedParams().Do(context.TODO())
	if err != nil {
		log.Fatalf("Failed to get suggested params: %+v", err)
	}

	var txnBuff = bytes.NewBuffer(nil)

	for idx, addr := range addrs {
		xfer, err := future.MakeAssetTransferTxn(*sender, addr, uint64(*amt), nil, suggestedParams, "", uint64(*asset))
		if err != nil {
			log.Printf("Failed to make asset xfer: %+v", err)
			continue
		}

		signed, err := kmdClient.SignTransactionWithSpecificPublicKey(*wallet, string(walletpw), xfer, pk)
		if err != nil {
			log.Printf("Failed to create signed transaction: %+v", err)
			continue
		}

		txnBuff.Write(signed.SignedTransaction)

		if idx%16 == 0 {
			_, err = algodClient.SendRawTransaction(txnBuff.Bytes()).Do(context.TODO())
			if err != nil {
				log.Fatalf("Failed to send raw Transaction: %+v", err)
			}

			txnBuff.Reset()
			suggestedParams, err = algodClient.SuggestedParams().Do(context.TODO())
			if err != nil {
				log.Fatalf("Failed to get suggested params: %+v", err)
			}
		}
	}
}

func getAddrs(fname string) []string {
	var addrs []string

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := types.DecodeAddress(scanner.Text())
		if err != nil {
			log.Printf("Failed to decode address: %s", err)
			continue
		}

		addrs = append(addrs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return addrs
}
