module github.com/almog-t/light-client-poc

go 1.17

require github.com/algorand/go-algorand-sdk v1.17.0

require (
	github.com/algorand/falcon v0.0.0-20220419072721-9f9785b53dd1 // indirect
	github.com/algorand/go-algorand v0.22.0-crypto-split // indirect
	github.com/algorand/go-codec/codec v1.1.8 // indirect
	github.com/algorand/go-deadlock v0.2.2 // indirect
	github.com/algorand/go-sumhash v0.1.0 // indirect
	github.com/algorand/msgp v1.1.52 // indirect
	github.com/aws/aws-sdk-go v1.16.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davidlazar/go-crypto v0.0.0-20170701192655-dcfb0a7ac018 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/olivere/elastic v6.2.14+incompatible // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	gopkg.in/sohlich/elogrus.v3 v3.0.0-20180410122755-1fa29e2f2009 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/algorand/go-algorand-sdk v1.17.0 => github.com/almog-t/go-algorand-sdk v1.14.1-0.20220726115802-16e2c7121e88

replace github.com/algorand/go-algorand v0.22.0-crypto-split => github.com/algonathan/go-algorand v0.22.0-crypto-split
