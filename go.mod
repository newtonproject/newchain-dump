module github.com/newtonproject/newchain-dump

go 1.14

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.10.15
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/peterh/liner v1.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	upper.io/db.v3 v3.7.1+incompatible
)

replace github.com/ethereum/go-ethereum => github.com/newtonproject/newchain v1.10.15-newton
