package miner

import (
	"encoding/hex"
	"github.com/gnc-project/poc"
	"github.com/gnc-project/poc/chiapos"
	"github.com/gnc-project/poc/difficulty"
	"log"
	"math/big"
	"sync/atomic"
	"time"
)


func Mine(quit chan struct{},commit chan interface{},plots []*chiapos.DiskProver,challenge [32]byte,
		number uint64,lastBlockTime time.Time,diff *big.Int) error {

	blockTime := lastBlockTime.Add(1 * poc.PoCSlot * time.Second)
	var workSlot = uint64(blockTime.Unix()) / poc.PoCSlot

	var bestQuality = big.NewInt(0)
	var bestChiaQualityIndex int

	ticker := time.NewTicker(time.Second * poc.PoCSlot /4)
	defer ticker.Stop()

	chiaQualities := GetChiaQualities(plots,challenge)
	if len(chiaQualities) == 0 {
		log.Println("not found qualities")
		return nil
	}

	search:
		for{
			select {
			case <-quit:
				return nil
			case <-ticker.C:
				nowSlot := uint64(time.Now().Unix()) / poc.PoCSlot
				if workSlot > nowSlot+poc.AllowAhead  {
					log.Println( "mining too far in the future",
						"nowSlot",nowSlot, "workSlot", workSlot)
					continue search
				}

				// Try to solve, until workSlot reaches nowSlot+allowAhead
				for i := workSlot; i <= nowSlot+poc.AllowAhead ; i++ {

					select {
					case <-quit:
						return nil
					default:
					}

					// Ensure there are valid qualities
					qualities, err := GetGNCQualities(chiaQualities, workSlot, number)
					if err != nil {
						return err
					}
					if len(qualities) == 0 {
						continue search
					}

					// find best quality
					bestQuality.SetUint64(0)
					for iq, quality := range qualities {
						if quality.Cmp(bestQuality) > 0 {
							bestQuality = quality
							bestChiaQualityIndex = iq
						}
					}
					nextDiff := difficulty.CalcNextRequiredDifficulty(lastBlockTime,diff,blockTime)
					if bestQuality.Cmp(nextDiff) > 0 {
						bestChiaQuality := chiaQualities[bestChiaQualityIndex]
						proof, err := poc.GetGNCProof(challenge,bestChiaQuality.Index,bestChiaQuality.Plot)
						if err != nil {
							log.Println("get proof err",err.Error())
							return err
						}
						id := bestChiaQuality.Plot.ID()
						pid := hex.EncodeToString(id[:])

						commit <- &poc.Commit{
							pid,
							hex.EncodeToString(proof),
							bestChiaQuality.Plot.Size(),
							nextDiff,
							number,
							blockTime.Unix(),
						}
						return nil
					}

					// increase slot and header Timestamp
					atomic.AddUint64(&workSlot, 1)
					blockTime = blockTime.Add(poc.PoCSlot * time.Second)
				}
				continue search
			}
		}
}