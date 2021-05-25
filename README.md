# Airdrop
Airdrop an asset on algorand

Takes a list of addresses, one address per row and tries to submit an asset xfer on algorand for each one.

Run it with something like this: 
```
./Airdrop -airdrop-asset-id 10  \
        -airdrop-asset-amt 100 \
        -airdrop-from GBDUTCSTHDO3RFUACXJ2QOM6ABHGXEC246VSDB6RARB2EO5JBVL7HD273Q \
        -airdrop-wallet wbx
```


TODO:

 - Add a Dryrun option
 - Thread out txns for faster processing
 - Actually check to see if the txn went thru 
 - Query the accts prior to sending to make sure we're not double sending in the case of a fail/restart

