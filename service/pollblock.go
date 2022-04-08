package service

import (
	"context"
	"errors"
	"math/big"
	"time"
)

func (s *Service) pollBlocks(ctx context.Context, sysErr chan<- error, ding chan<- *big.Int) error {
	glog.Info("Polling Blocks...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-ctx.Done():
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				glog.Error("Polling failed, retries exceeded")
				sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := s.conn.LatestBlock()
			if err != nil {
				glog.Error("Unable to get latest block", "block", s.currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, s.currentBlock).Cmp(BlockDelay) == -1 {
				//s.log.Debug("Block not ready, will retry", "target", s.currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			if big.NewInt(0).Mod(s.currentBlock, big.NewInt(15)).Cmp(big.NewInt(0)) == 0 {
				glog.Debug("pollBlocks", "currentBlock", s.currentBlock)
			}

			// Goto next block and reset retry counter
			if s.isTimeToCalculateRate() {
				blk := big.NewInt(0)
				ding <- blk.Add(blk, s.currentBlock)
			}
			s.currentBlock.Add(s.currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

func (s *Service) isTimeToCalculateRate() bool {
	m := big.NewInt(0)

	big.NewInt(0).DivMod(s.currentBlock, s.cfg.blockInterval, m)
	return m.Uint64() == 0
}
