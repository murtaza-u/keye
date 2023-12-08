/*
Package client is a library that provides abstracted methods for
interacting with the Keye database, including GET, PUT, DEL, and WATCH
operations.
*/
package client

import (
	pb "github.com/murtaza-u/keye"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// C represents the Keye database client.
type C struct {
	Config
	conn *grpc.ClientConn
	api  pb.ApiClient
}

// New creates a Keye client based on the provided configuration.
func New(conf Config) (*C, error) {
	conf.impute()
	c, err := grpc.Dial(conf.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &C{
		Config: conf,
		conn:   c,
		api:    pb.NewApiClient(c),
	}, nil
}

// Close terminates the connection to the Keye database server.
func (c *C) Close() error { return c.conn.Close() }
