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
	"strings"
	"sync"

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
	kmdPort  = flag.Int("airdrop-kmd-port", 7833, "port of Key Management Daemon")
	kmdToken = flag.String("airdrop-kmd-token", "", "token to use when making requests to Key Management Daemon")

	fname = flag.String("airdrop-address-file", "address.csv", "Path to file containing list of addresses to airdrop to")

	asset = flag.Int("airdrop-asset-id", 0, "Asset id to transfer")
	amt   = flag.Int("airdrop-asset-amt", 0, "Amount of asset to transfer")

	sender = flag.String("airdrop-wallet-address", "", "Address to use as sender")
	wallet = flag.String("airdrop-wallet-id", "", "Id of wallet to use")

	dryrun = flag.Bool("airdrop-dryrun", false, "Don't actually submit transactions")
)

func main() {
	flag.Parse()

	if *sender == "" || *wallet == "" {
		fmt.Println("Must include sender and wallet name")
		flag.PrintDefaults()
		return
	}

	if *amt == 0 || *asset == 0 {
		fmt.Println("No asset id or amount specified")
		flag.PrintDefaults()
		return
	}

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
		log.Fatalf("Failed to make algod client: %+v", err)
	}

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
	pass := strings.TrimSpace(string(walletpw))

	handle, err := kmdClient.InitWalletHandle(*wallet, pass)
	if err != nil {
		log.Fatalf("Couldn't get a handle to the wallet: %+v", err)
	}

	// Read in file
	addrs := getAddrs(*fname)
	log.Printf("Got %d addresses from %s", len(addrs), *fname)

	// Get suggested params from algod
	suggestedParams, err := algodClient.SuggestedParams().Do(context.TODO())
	if err != nil {
		log.Fatalf("Failed to get suggested params: %+v", err)
	}

	var (
		wg      = &sync.WaitGroup{}
		txnBuff = bytes.NewBuffer(nil)
	)
	for idx, addr := range addrs {
		xfer, err := future.MakeAssetTransferTxn(*sender, addr, uint64(*amt), nil, suggestedParams, "", uint64(*asset))
		if err != nil {
			log.Printf("Failed to make asset xfer: %+v", err)
			continue
		}

		signed, err := kmdClient.SignTransactionWithSpecificPublicKey(handle.WalletHandleToken, pass, xfer, pk)
		if err != nil {
			log.Printf("Failed to create signed transaction: %+v", err)
			continue
		}

		txnBuff.Write(signed.SignedTransaction)

		// For every 16 (current max grouped txns), try to send
		if idx%16 == 0 {

			if *dryrun {
				log.Printf("DRYRUN: send %d bytes, addr index %d", txnBuff.Len(), idx)
			} else {
				wg.Add(1)
				go sendTxn(algodClient, txnBuff.Bytes(), wg)
			}

			txnBuff.Reset()

			//Refresh suggested params since they may have changed singe we last requested them
			suggestedParams, err = algodClient.SuggestedParams().Do(context.TODO())
			if err != nil {
				log.Fatalf("Failed to get suggested params: %+v", err)
			}

		}
	}

	wg.Wait()
}

func sendTxn(algodClient *algod.Client, rawb []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	txid, err := algodClient.SendRawTransaction(rawb).Do(context.TODO())
	if err != nil {
		log.Printf("Failed to send raw Transaction: %+v", err)
	}

	log.Printf("Sent %d bytes with txid %s", len(rawb), txid)

	for {
		resp, _, err := algodClient.PendingTransactionInformation(txid).Do(context.TODO())
		if err != nil {
			log.Printf("Failed to get pending txn info for %s: %+v", txid, err)
			break
		}

		if resp.ConfirmedRound != 0 {
			log.Printf("Txn %s success", txid)
			break
		}

		if len(resp.PoolError) > 0 {
			log.Printf("Pool Error: %s", resp.PoolError)
			break
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
		// Check that its a valid address
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
