package runtime_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
)

var (
	zeroAddr common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
	value                   = new(big.Int)
	gas      uint64         = 5000
	code                    = []byte{
		evm.PUSH1, 0x01, evm.PUSH1, 0x02, evm.ADD,
		evm.PUSH1, 0x00, evm.MSTORE8,
		evm.PUSH1, 0x01, evm.PUSH1, 0x00, evm.RETURN,
	}
	cfg = dbutils.NewConfigFromEnv()
)

func TestRuntime(t *testing.T) {
	testEvm := evm.NewEVM()
	contract := runtime.NewContract(1, zeroAddr, zeroAddr, zeroAddr, value, gas, code)
	config := &runtime.ForksInTime{
		EIP158: true,
	}

	stateDb, err := db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	// TODO: Improve when state transition is implemented

	// store := tree.NewPostgresStore(stateDb)
	// mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
	// host := runtime.NewMerkleTreeHost(tree.NewStateTree(mt, nil))

	res := testEvm.Run(contract, nil, config)
	assert.Equal(t, uint64(4976), res.GasLeft)
	assert.NoError(t, res.Err)
}