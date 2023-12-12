package srv

import (
	"bytes"

	pb "github.com/murtaza-u/keye"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Backup creates a database snapshot and streams it to the client in
// chunks.
func (s *Srv) Backup(in *pb.ChunkSize, stream pb.Api_BackupServer) error {
	size := in.Size
	if size == 0 {
		return status.Error(codes.InvalidArgument, "chunk size cannot be 0")
	}

	err := s.db.View(func(tx *bbolt.Tx) error {
		buf := new(bytes.Buffer)
		n, err := tx.WriteTo(buf)
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to create backup: %s", err.Error())
		}

		data := buf.Next(int(size))
		for i := size; data != nil && len(data) != 0; i += size {
			err := stream.Send(&pb.Chunk{
				Data: data,
				More: n > i,
			})
			if err != nil {
				return err
			}
			data = buf.Next(int(size))
		}

		return nil
	})

	return err
}
