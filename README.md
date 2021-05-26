# Ⱥirdrop
Airdrop an asset on algorand

Takes a list of addresses, one address per row and tries to submit an asset xfer on algorand for each one.


Flags
-----

```
Usage of ./airdrop:
  -airdrop-address-file string
        Path to file containing list of addresses to airdrop to (default "address.csv")
  -airdrop-algod-host string
        hostname of Algorand node (default "127.0.0.1")
  -airdrop-algod-port int
        port of Algorand node (default 4001)
  -airdrop-algod-token string
        token to use when making requests to Algorand node
  -airdrop-asset-amt int
        Amount of asset to transfer
  -airdrop-asset-id int
        Asset id to transfer
  -airdrop-dryrun
        Don't actually submit transactions
  -airdrop-kmd-host string
        hostname of Key Management Daemon (default "127.0.0.1")
  -airdrop-kmd-port int
        port of Key Management Daemon (default 7833)
  -airdrop-kmd-token string
        token to use when making requests to Key Management Daemon
  -airdrop-wallet-address string
        Address to use as sender
  -airdrop-wallet-name string
        Name of wallet to use
```



Try it:
------


cd into directory
```
go build 
```

Run a dryrun to check file and parameters
```
./airdrop -airdrop-asset-id 5  \
        -airdrop-asset-amt 10 \
        -airdrop-kmd-port 4002 \
        -airdrop-algod-token aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
        -airdrop-kmd-token aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
        -airdrop-wallet-address 7LQ7U4SEYEVQ7P4KJVCHPJA5NSIFJTGIEXJ4V6MFS4SL5FMDW6MYHL2JXM  \
        -airdrop-wallet-id "446f2cb0af0e478ce48af3c6c4c95693" \
        -airdrop-address-file addrs.csv
```

Then just lop off -airdrop-dryrun and it should just work ™


Tested this in my local sandbox and it works but YMMV

TODO:
 - Query the accts prior to sending to make sure we're not double sending in the case of a fail/restart

