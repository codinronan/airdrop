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
./airdrop -airdrop-asset-id 10  \
        -airdrop-asset-amt 100 \
        -airdrop-wallet-address GBDUTCSTHDO3RFUACXJ2QOM6ABHGXEC246VSDB6RARB2EO5JBVL7HD273Q \
        -airdrop-wallet-name wbx \ 
	-airdrop-dryrun 
```

Then just lop off -airdrop-dryrun and it should just work ™



TODO:
 
 - Actually test it
 - Query the accts prior to sending to make sure we're not double sending in the case of a fail/restart

