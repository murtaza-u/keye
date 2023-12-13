package srv

import (
	"context"
	"time"

	"github.com/murtaza-u/keye/internal/pb"
)

// Stats implements the gRPC API service Stat method. It returns
// statistics over the given delta time range.
func (s *Srv) Stats(ctx context.Context, in *pb.Delta) (*pb.Measures, error) {
	prev := s.db.Stats()
	time.Sleep(time.Duration(in.GetDuration()))
	curr := s.db.Stats()
	diff := curr.Sub(&prev)
	return &pb.Measures{
		TxStats: &pb.TxStats{
			PageCount:     diff.TxStats.GetPageCount(),
			PageAlloc:     diff.TxStats.GetPageAlloc(),
			CursorCount:   diff.TxStats.GetCursorCount(),
			NodeCount:     diff.TxStats.GetNodeCount(),
			NodeDeref:     diff.TxStats.GetNodeDeref(),
			Rebalance:     diff.TxStats.GetRebalance(),
			RebalanceTime: int64(diff.TxStats.GetRebalanceTime()),
			Split:         diff.TxStats.GetSplit(),
			Spill:         diff.TxStats.GetSpill(),
			SpillTime:     int64(diff.TxStats.GetSpillTime()),
			Write:         diff.TxStats.GetWrite(),
			WriteTime:     int64(diff.TxStats.GetWriteTime()),
		},
		FreePageN:     int32(diff.FreePageN),
		PendingPageN:  int32(diff.PendingPageN),
		FreeAlloc:     int32(diff.FreeAlloc),
		FreelistInuse: int32(diff.FreelistInuse),
		TxN:           int32(diff.TxN),
		OpenTxN:       int32(diff.OpenTxN),
	}, nil
}
