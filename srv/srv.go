/*
Package srv defines the database server. It implements the gRPC server
and is responsible for GET, PUT, DEL and WATCH operations on the
database. The database is internally powered by bboltDB
(https://github.com/etcd-io/bbolt).
*/
package srv

import (
	"fmt"
	"net"

	pb "github.com/murtaza-u/keye"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const bucket = "KEYE"

// Srv implements the gRPC API server and performs database operations
// on request.
type Srv struct {
	opts
	db *bbolt.DB

	pb.UnimplementedApiServer
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

	db, err := bbolt.Open(opts.path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize boltdb database: %w", err)
	}

	return &Srv{
		opts: opts,
		db:   db,
	}, nil
}

// Run starts the database server. This is a blocking method and only
// returns if there is an error.
func (s *Srv) Run() error {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	grpcS := grpc.NewServer()
	pb.RegisterApiServer(grpcS, s)
	if s.reflect {
		reflection.Register(grpcS)
	}
	defer s.db.Close()
	return grpcS.Serve(ln)
}

// close closes the underlying boltdb database.
func (s *Srv) close() error { return s.db.Close() }
