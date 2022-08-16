
Light Client PoC
====================


A Go implementation of a light client, able to ingest Algorand state proofs and transactions.

# Background

This implementation contains both code comprising the light client itself and encoded assets needed to demonstrate the light client's functionalities.

Assets include data required to initialize the Oracle, state proof data to advance the Oracle's state and data required to verify a transaction using the Oracle's state. All assets were generated using a private Algorand network.

The code is intended to be used as a guide to building light clients.


# Reading Order

1. main.go - read in its entirety for a general overview of the light client.
2. oracle.go - start from the commentary on the Oracle struct, followed by the commentary on AdvanceState. Branch out as needed.
3. transactionVerifier.go - start from the commentary on verifyTransaction. Branch out as needed.

# Running the Light Client
Building the light client using
```bash
go build
```
and running the resulting binary will cause the light client to operate on the assets in encodedassets.
An Oracle will be initialized using the data in the [genesis folder](encodedassets/genesis), its state will advance using the data in the [state proof folder](encodedassets/stateproofverification) and a transaction will be verified using the data in the [transaction verification folder](encodedassets/transactionverification).
