module github.com/almog-t/light-client-poc

go 1.17

require github.com/algorand/go-algorand-sdk v1.17.0

require (
	github.com/algorand/falcon v0.0.0-20220419072721-9f9785b53dd1 // indirect
	github.com/algorand/go-codec/codec v1.1.8 // indirect
	github.com/algorand/go-sumhash v0.1.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
)

replace github.com/algorand/go-algorand-sdk v1.17.0 => github.com/almog-t/go-algorand-sdk v1.14.1-0.20220808083618-7c46fe8f3e64
