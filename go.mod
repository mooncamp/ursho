module github.com/douglasmakey/ursho

go 1.13

require (
	github.com/dgraph-io/dgo/v2 v2.0.0
	github.com/lib/pq v1.0.0
	google.golang.org/grpc v1.23.1
	mooncamp.com/dgx v0.0.0-20191002071046-556ca06f1e0b
)

replace github.com/dgraph-io/dgraph v1.1.0 => github.com/dgraph-io/dgraph v1.1.1-0.20190910201909-32c09acf647c

replace github.com/dgraph-io/badger v1.6.0 => github.com/dgraph-io/badger v1.6.1-0.20190903194648-fbcd60898198
