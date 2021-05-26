# Airdrop
Airdrop an asset on algorand

Takes a list of addresses, one address per row and tries to submit an asset xfer on algorand for each one.

Try it with something like this: 

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

Then just lop off -airdrop-dryrun 

TODO:
 
 - Actually test it
 - Query the accts prior to sending to make sure we're not double sending in the case of a fail/restart

