package miner

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gnc-project/poc"
	"github.com/gnc-project/poc/chiapos"
	"github.com/gnc-project/poc/difficulty"
	"log"
	"math/big"
	"sync/atomic"
	"time"
)

func Mine(quit chan struct{},client *rpc.Client,plots []*chiapos.DiskProver,challenge [32]byte,
		number uint64,lastBlockTime time.Time,diff *big.Int) error {

	blockTime := lastBlockTime.Add(1 * poc.PoCSlot * time.Second)
	var workSlot = uint64(blockTime.Unix()) / poc.PoCSlot

	var bestQuality = big.NewInt(0)
	var bestChiaQualityIndex int

	ticker := time.NewTicker(time.Second * poc.PoCSlot / 4)
	defer ticker.Stop()

	chiaQualities := GetChiaQualities(plots,challenge)
	if len(chiaQualities) == 0 {
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
					for i, quality := range qualities {
						if quality.Cmp(bestQuality) > 0 {
							bestQuality = quality
							bestChiaQualityIndex = i
						}
					}

					if bestQuality.Cmp(difficulty.CalcNextRequiredDifficulty(lastBlockTime,diff,blockTime)) > 0 {
						bestChiaQuality := chiaQualities[bestChiaQualityIndex]
						proof, err := bestChiaQuality.Plot.GetFullProof(challenge,bestChiaQuality.Index)
						if err != nil {
							log.Println("get proof err",err.Error())
						}
						id := bestChiaQuality.Plot.ID()
						pid := hex.EncodeToString(id[:])

						return addPlot(client,pid,hex.EncodeToString(proof),bestChiaQuality.Plot.Size(),bestQuality,number,blockTime.Unix())
					}

					// increase slot and header Timestamp
					atomic.AddUint64(&workSlot, 1)
					blockTime = blockTime.Add(poc.PoCSlot * time.Second)
				}
			}
		}
}

//mining
func addPlot(client *rpc.Client,pid,proof string,k uint8,quality *big.Int,number uint64,timestamp int64) error {
	result := make(map[string]interface{})
	err := client.Call(&result,"eth_addPlot",pid,proof,k,quality,number,timestamp)
	if err != nil {
		return err
	}
	log.Printf("blockNumber=%v,deadline=%.2f",result["blockNumber"],result["deadline"].(float64))
	return nil
}