module github.com/threeq/xid

go 1.13

require (
	github.com/fzipp/gocyclo v0.0.0-20150627053110-6acd4345c835 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
)

// demo:
//      replace golang.org/x/crypto => github.com/golang/crypto latest

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/net => github.com/golang/net v0.0.0-20190918130420-a8b05e9114ab
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190916202348-b4ddaad3f8a3
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190918234917-7baacfbe02f2
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
)
