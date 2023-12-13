/*
Package srv defines the database server. It implements the gRPC server
and is responsible for GET, PUT, DEL and WATCH operations on the
database. The database is internally powered by bboltDB
(https://github.com/etcd-io/bbolt).
*/
package srv

import (
	"errors"
	"fmt"
	"net"

	"github.com/murtaza-u/keye/internal/pb"
	"github.com/murtaza-u/keye/watch"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// ErrKeyNotFound is returned when the specified key does not exist
	// in the database.
	ErrKeyNotFound = errors.New("key not found")
	// ErrNoKeysMatchRegex is returned when none of the keys the the
	// database match the given regex.
	ErrNoKeysMatchRegex = errors.New("no keys match the given regex")
)

const bucket = "KEYE"

// Srv implements the gRPC API server and performs database operations
// on request.
type Srv struct {
	opts
	db      *bbolt.DB
	watcher *watch.W

	pb.UnimplementedApiServer
}

// New creates a new database server from the given options.
func New(optfns ...optFunc) (*Srv, error) {
	opts := defaultOpts()
	for _, fn := range optfns {
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
		opts:    opts,
		db:      db,
		watcher: watch.New(opts.eventQueueSize),
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

	go s.watcher.Listen()
	defer s.watcher.Close()

	defer s.close()

	return grpcS.Serve(ln)
}

// close closes the underlying boltdb database.
func (s *Srv) close() error { return s.db.Close() }
