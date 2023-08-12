module github.com/threeq/xid

go 1.20

require (
	github.com/aceld/zinx v1.2.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/stretchr/testify v1.8.4
	github.com/valyala/fasthttp v1.48.0
	gopkg.in/go-playground/assert.v1 v1.2.1
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/klauspost/reedsolomon v1.11.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible // indirect
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
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
