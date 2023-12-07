/*
Package srv defines the database server. It implements the gRPC server
and is responsible for GET, PUT, DEL and WATCH operations on the
database. The database is internally powered by bboltDB
(https://github.com/etcd-io/bbolt).
*/
package srv

import "github.com/murtaza-u/keye"

// Srv implements the gRPC API server and performs database operations
// on request.
type Srv struct {
	opts

	keye.UnimplementedApiServer
}

// New creates a new database server from the given options.
func New(optfuncs ...optFunc) (*Srv, error) {
	opts := defaultOpts()
	for _, fn := range optfuncs {
		err := fn(&opts)
		if err != nil {
			return nil, err
		}
	}
	return &Srv{
		opts: opts,
	}, nil
}
