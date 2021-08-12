package miner

import (
	"github.com/gnc-project/poc"
	"github.com/gnc-project/poc/chiapos"
	"math/big"
)

type WorkSpaceQuality struct {
	Plot          *chiapos.DiskProver
	Index         uint32
	Quality       []byte
	Error         error
}

func GetChiaQualities(plots []*chiapos.DiskProver,challenge [32]byte) []*WorkSpaceQuality {
	var ws []*WorkSpaceQuality
	for _,p := range plots {
		if pass := chiapos.PassPlotFilter(p.ID(), challenge); !pass {
			continue
		}
		posChallenge := chiapos.CalculatePosChallenge(p.ID(), challenge)
		q,e := p.GetQualitiesForChallenge(posChallenge)
		if e != nil {
			continue
		}
		if len(q) == 0 {
			continue
		}
		for k,v :=range q {
			ws = append(ws,&WorkSpaceQuality{
				p,
				uint32(k),
				v,
				nil,
			})
		}
	}
	return ws
}

func GetGNCQualities(chiaQualities []*WorkSpaceQuality,slot, height uint64)([]*big.Int, error)  {
	GNCQualities := make([]*big.Int, len(chiaQualities))
	for i, chiaQuality := range chiaQualities {
		chiaHashVal := poc.HashValChia(chiaQuality.Quality, slot, height)
		GNCQualities[i] = poc.GetQuality(poc.Q1FactorChia(chiaQuality.Plot.Size()), chiaHashVal)
	}
	return GNCQualities, nil
}
