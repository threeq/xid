module github.com/threeq/xid

go 1.13

require (
	github.com/gin-gonic/gin v1.7.0
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/go-playground/assert.v1 v1.2.1
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
