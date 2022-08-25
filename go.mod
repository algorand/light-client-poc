module github.com/almog-t/light-client-poc

go 1.17

require (
	github.com/algorand/go-algorand-sdk v1.17.0
	github.com/algorand/go-stateproof-verification v0.0.0-20220824090814-9c32fdae6e05
)

require (
	github.com/algorand/falcon v0.0.0-20220727072124-02a2a64c4414 // indirect
	github.com/algorand/go-codec/codec v1.1.9 // indirect
	github.com/algorand/go-sumhash v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
)

replace (
	github.com/algorand/go-algorand-sdk v1.17.0 => github.com/almog-t/go-algorand-sdk v1.14.1-0.20220815130753-4f4f4a46360f
	github.com/algorand/go-stateproof-verification v0.0.0-20220824090814-9c32fdae6e05 => github.com/algorand/go-stateproof-verification v0.0.0-20220824112425-645eeea0c7d7
)
