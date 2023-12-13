package client

import (
	"context"
	"errors"
	"time"

	pb "github.com/murtaza-u/keye"

	"google.golang.org/grpc/status"
)

// Measures represents statistics about the database.
type Measures struct {
	// global, ongoing stats
	TxStats *TxStats

	// freelist stats
	FreePageN     int32 // total number of free pages on the freelist
	PendingPageN  int32 // total number of pending pages on the freelist
	FreeAlloc     int32 // total bytes allocated in free pages
	FreelistInuse int32 // total bytes used by the freelist

	// transaction stats
	TxN     int32 // total number of started read transactions
	OpenTxN int32 // number of currently open read transactions
}

// TxStats represents statistics about the actions performed by the
// transaction.
type TxStats struct {
	// page statistics
	PageCount int64 // number of page allocations
	PageAlloc int64 // total bytes allocated

	// cursor statistics
	CursorCount int64 // number of cursors created

	// node statistics
	NodeCount int64 // number of node allocations
	NodeDeref int64 // number of node dereferences

	// rebalance statistics
	Rebalance     int64         // number of node rebalances
	RebalanceTime time.Duration // total time spent rebalancing

	// split/spill statistics
	Split     int64         // number of nodes split
	Spill     int64         // number of nodes spilled
	SpillTime time.Duration // total time spent spilling

	// write statistics
	Write     int64         // number of writes performed
	WriteTime time.Duration // total time spent writing to disk
}

// Stats calculates database performance statistics over a delta time
// range.
func (c *C) Stats(delta time.Duration) (*Measures, error) {
	ctx, cancel := context.WithTimeout(context.Background(), delta+c.Timeout)
	defer cancel()

	stats, err := c.api.Stats(ctx, &pb.Delta{
		Duration: int64(delta),
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, errors.New(stat.Message())
		}
		return nil, err
	}

	return &Measures{
		TxStats: &TxStats{
			PageCount:     stats.TxStats.PageCount,
			PageAlloc:     stats.TxStats.PageAlloc,
			CursorCount:   stats.TxStats.CursorCount,
			NodeCount:     stats.TxStats.NodeCount,
			NodeDeref:     stats.TxStats.NodeDeref,
			Rebalance:     stats.TxStats.Rebalance,
			RebalanceTime: time.Duration(stats.TxStats.RebalanceTime),
			Split:         stats.TxStats.Split,
			Spill:         stats.TxStats.Spill,
			SpillTime:     time.Duration(stats.TxStats.SpillTime),
			Write:         stats.TxStats.Write,
			WriteTime:     time.Duration(stats.TxStats.WriteTime),
		},
		FreePageN:     stats.FreePageN,
		PendingPageN:  stats.PendingPageN,
		FreeAlloc:     stats.FreeAlloc,
		FreelistInuse: stats.FreelistInuse,
		TxN:           stats.TxN,
		OpenTxN:       stats.OpenTxN,
	}, nil
}
