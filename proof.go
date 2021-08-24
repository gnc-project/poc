package poc

import (
	"github.com/gnc-project/poc/chiapos"
	"math"
	"math/big"
)

const (
	// PoCSlot represents the unit Slot of PoC
	PoCSlot = 3

	AllowAhead = 1

	ToleranceSlot = 4

	// KiB represents KiByte
	KiB = 1024

	// MiB represents MiByte
	MiB = 1024 * KiB

	// MinValidDefaultBitLength represents smallest default proof BitLength
	MinValidDefaultBitLength = 24

	// MaxValidDefaultBitLength represents biggest default BitLength
	MaxValidDefaultBitLength = 40

	// MinValidDefaultBitLength represents smallest default proof BitLength
	MinValidChiaBitLength = chiapos.MinPlotSize

	// MaxValidDefaultBitLength represents biggest default BitLength
	MaxValidChiaBitLength = chiapos.MaxPlotSize

	// QualityConstantMASSIP0002 multiplies the original quality to fix proof filter
	QualityConstantGNC = 512

	// QualityConstantMASSValidity represents the fix between mass and chia proof
	QualityConstantGNCValidity = 0.329
)

type ProofType uint8

const (
	ProofTypeChia    ProofType = 1
	ProofTypeEmpty   ProofType = math.MaxUint8
)

// EnsureBitLength returns whether it is a valid bitLength.
func (pt ProofType) EnsureBitLength(bl int) bool {
	switch pt {
	case ProofTypeChia:
		return MinValidChiaBitLength <= bl && bl <= MaxValidChiaBitLength
	default:
		return false
	}
}

func (pt ProofType) PlotSize(bl int) uint64 {
	if !pt.EnsureBitLength(bl) {
		return 0
	}
	switch pt {
	case ProofTypeChia:
		return ChiaPlotSize(bl)
	default:
		return 0
	}
}

func PlotSize(proofType ProofType, bl int) uint64 {
	return proofType.PlotSize(bl)
}

type Proof interface {
	Type() ProofType
	BitLength() int
	Encode() []byte
	Decode([]byte) error
	Quality(slot, height uint64) *big.Int
	Verify(seed [32]byte, challenge [32]byte, filter bool) error
	VerifiedQuality(seed [32]byte, challenge [32]byte, filter bool, slot, height uint64) (*big.Int, error)
}

// VerifyProof verifies proof.
func VerifyProof(proof Proof, pubKeyHash [32]byte, challenge [32]byte, filter bool) error {
	return proof.Verify(pubKeyHash, challenge, filter)
}

// GetQuality produces the relative quality with factor Q1
func GetQuality(Q1 *big.Float, hashVal [32]byte) *big.Int {
	// Note: FH = H in big.Float
	FH := new(big.Float).SetInt(new(big.Int).SetBytes(hashVal[:]))
	F64H, _ := FH.Float64()

	// Note: log2FH = log2(H)
	log2FH := big.NewFloat(math.Log2(F64H))

	// Note: Q2 = 256 - log2(H)
	Q2 := big.NewFloat(256)
	Q2.Sub(Q2, log2FH)
	if Q2.Cmp(big.NewFloat(0)) <= 0 {
		return big.NewInt(0)
	}

	quality := big.NewInt(0)
	big.NewFloat(0).Quo(Q1, Q2).Int(quality)
	return quality
}

func IsValidProofType(pt ProofType) bool {
	return pt == ProofTypeChia
}
