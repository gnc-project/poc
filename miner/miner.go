package miner

import (
	"encoding/hex"
	"fmt"
	"github.com/gnc-project/poc"
	"github.com/gnc-project/poc/chiapos"
	"github.com/gnc-project/poc/difficulty"
	"log"
	"math/big"
	"time"
)


func Mine(quit chan struct{},commit chan interface{},plots []*chiapos.DiskProver,challenge [32]byte,
		number uint64,lastBlockTime time.Time,diff *big.Int) error {

	blockTime := lastBlockTime.Add( poc.PoCSlot * time.Second)
	var workSlot = uint64(blockTime.Unix()) / poc.PoCSlot

	var bestQuality = big.NewInt(0)
	var bestChiaQualityIndex int

	chiaQualities := GetChiaQualities(plots,challenge)
	if len(chiaQualities) == 0 {
		log.Println("not found qualities")
		return nil
	}

	search:
		for{
			select {
			case <- quit:
				return nil
			default:
				select {
				case <- quit:
					return nil
				default:
				}
				if time.Now().Unix() < blockTime.Unix() {
					time.Sleep(500 * time.Millisecond)
					continue search
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

					dif,err := poc.VerifiedQuality(proof,id,challenge,uint64(blockTime.Unix()/poc.PoCSlot),number,uint64(bestChiaQuality.Plot.Size()))
					if err != nil {
						return err
					}
					if dif.Cmp(nextDiff) <= 0 {
						return fmt.Errorf("bestQuality verify is err dif=%v nextDiff=%v",dif,nextDiff)
					}

					commit <- &poc.Commit{
						Pid:        hex.EncodeToString(id[:]),
						Proof:      hex.EncodeToString(proof),
						K:          bestChiaQuality.Plot.Size(),
						Difficulty: nextDiff,
						Number:     number,
						Challenge:  hex.EncodeToString(challenge[:]),
						Timestamp:  blockTime.Unix(),
					}
					return nil
				}


				// increase slot and header Timestamp
				blockTime = blockTime.Add(poc.AllowAhead * time.Second)
				workSlot = uint64(blockTime.Unix()) / poc.PoCSlot

				log.Println("increase slot and header Timestamp","bestQuality=",bestQuality,"nextDiff",nextDiff,"blockTime",blockTime)
			}
		}
}
