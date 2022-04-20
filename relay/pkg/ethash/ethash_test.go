//go:build !ci

package ethash

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestEpoch1(t *testing.T) {
	ethash := New("./test", 0, 0)
	edata, err := ethash.GetEpochData(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, edata.MerkleNodes[0].String(), "102440233808978497570443060571864861601403660416248037251589510336707164571416")
}

func TestEpoch475(t *testing.T) {
	ethash := New("./test", 0, 0)
	edata, err := ethash.GetEpochData(475)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, edata.MerkleNodes[0].String(), "86794333638112600585166912134894146973886718995379343243688093087503049446227")
}

func TestLookup(t *testing.T) {
	ethash := New("./test", 0, 0)

	l1, l2, err := ethash.GetBlockLookups(14257704, 16957842275414857198, common.HexToHash("0x4e710542b313c81468f2787dd9d1324663b3370ad25cefc0f8be7754257f9b6a"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, l1[0].String(), "83421943766262542122791167114674843962167463147708128675835754176032138682042")
	assert.Equal(t, l2[0].String(), "44275787240451683914795652699717121364927282182489534770837402673149472584329")

}

func TestBytesToUint32Slice(t *testing.T) {
	a := []byte{0x10, 0x20, 0x30, 0x40}
	b := bytesToUint32Slice(a)
	assert.Equal(t, []uint32{1076895760}, b)
}
