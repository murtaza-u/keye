package client

import (
	"bytes"
	"context"
	"errors"

	"github.com/murtaza-u/keye/internal/pb"

	"google.golang.org/grpc/status"
)

// Backup requests and concatenates database snapshot.
func (c *C) Backup(chunkSize int64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout*2)
	defer cancel()

	stream, err := c.api.Backup(ctx, &pb.ChunkSize{
		Size: chunkSize,
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, errors.New(stat.Message())
		}
		return nil, err
	}

	buf := new(bytes.Buffer)

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return nil, errors.New(stat.Message())
			}
			return nil, err
		}

		data := chunk.GetData()
		if _, err = buf.Write(data); err != nil {
			return nil, err
		}

		if !chunk.GetMore() {
			stream.CloseSend()
			break
		}
	}

	return buf.Bytes(), nil
}
