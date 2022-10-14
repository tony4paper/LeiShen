package check

import (
	"leishen/pkg/etypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// aave transaction example:
	// 0x7fd5d3af1b30b64751d07426ccdebf4a0b7f34e693bf86498cb368a852f377d7 13208700
	aave = &Feature{
		Name:      "AAVE",
		StartFunc: Function{FuncName: "flashLoan(address,address,uint256,bytes)", FuncSign: common.Hex2Bytes("0x5cffe9de"[2:])},
		Events: []Event{
			{EventName: "FlashLoan(address,address,uint256,uint256,uint256,uint256)", Topic: common.HexToHash("0x5b8f46461c1dd69fb968f1a003acee221ea3e19540e350233b612ddb43433b55")},
		},
	}

	// dydx transaction example:
	// 0xb5c8bd9430b6cc87a0e2fe110ece6bf527fa4f170a4bc8cd032f768fc5219838 9484688
	dydx = &Feature{
		Name:      "dYdX",
		StartFunc: Function{FuncName: "operate((address,uint256)[],(uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[])", FuncSign: common.Hex2Bytes("0xa67a6a45"[2:])},
		Events: []Event{
			{EventName: "LogOperation(address)", Topic: common.HexToHash("0x91b01baeee3a24b590d112613814d86801005c7ef9353e7fc1eaeaf33ccf83b0")},
			{EventName: "LogWithdraw(address,uint256,uint256,((bool,uint256),(bool,uint128)),address)", Topic: common.HexToHash("0xbc83c08f0b269b1726990c8348ffdf1ae1696244a14868d766e542a2f18cd7d4")},
			{EventName: "LogCall(address,uint256,address)", Topic: common.HexToHash("0xab38cdc4a831ebe6542bf277d36b65dbc5c66a4d03ec6cf56ac38de05dc30098")},
			{EventName: "LogDeposit(address,uint256,uint256,((bool,uint256),(bool,uint128)),address)", Topic: common.HexToHash("0x2bad8bc95088af2c247b30fa2b2e6a0886f88625e0945cd3051008e0e270198f")},
		},
	}

	// uniswapv2 transaction example:
	// 0x46a03488247425f845e444b9c10b52ba3c14927c687d38287c0faddc7471150a 11256673
	uniswapV2 = &Feature{
		Name:      "Uniswap V2",
		StartFunc: Function{FuncName: "swap(uint256,uint256,address,bytes)", FuncSign: common.Hex2Bytes("0x022c0d9f"[2:])},
		Events: []Event{
			{EventName: "Swap(address,uint256,uint256,uint256,uint256,address)", Topic: common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")},
		},
		InternalFuncs: []Function{
			{FuncName: "uniswapV2Call(address,uint256,uint256,bytes)", FuncSign: common.Hex2Bytes("0x10d1e85c"[2:])},
		},
	}

	// uniswapv3 transaction example:
	// 0xebd70cbcb41e39fe54d650716ec0ca0c7d9bf57b926f1829d2720d03898ecb6f 13666739
	uniswapV3 = &Feature{
		Name:      "Uniswap V3",
		StartFunc: Function{FuncName: "flash(address,uint256,uint256,bytes)", FuncSign: common.Hex2Bytes("0x490e6cbc"[2:])},
		Events: []Event{
			{EventName: "Flash(address,address,uint256,uint256,uint256,uint256)", Topic: common.HexToHash("0xbdbdb71d7860376ba52b25a5028beea23581364a40522f6bcfb86bb1f2dca633")},
		},
	}

	AllFeatures = map[string]*Feature{aave.Name: aave, dydx.Name: dydx, uniswapV2.Name: uniswapV2, uniswapV3.Name: uniswapV3}
)

type Feature struct {
	Name          string
	StartFunc     Function
	Events        []Event
	InternalFuncs []Function
}

func (f *Feature) MayBeInBlock(block *types.Block) bool {
	bloom := block.Bloom()
	for _, event := range f.Events {
		if !bloom.Test(event.Topic[:]) {
			return false
		}
	}

	return true
}

func (f *Feature) MayBeInTransaction(receipt *types.Receipt) bool {
	bloom := receipt.Bloom
	for _, event := range f.Events {
		if !bloom.Test(event.Topic[:]) {
			return false
		}
	}

	return true
}

func (f *Feature) IsFlashLoan(receipt *types.Receipt, internalTxs []*etypes.InternalTransaction) bool {
	for eventIndex, logIndex := 0, 0; eventIndex < len(f.Events); eventIndex++ {
		haveTopic := false
		for ; logIndex < len(receipt.Logs); logIndex++ {
			if len(receipt.Logs[logIndex].Topics) != 0 &&
				receipt.Logs[logIndex].Topics[0] == f.Events[eventIndex].Topic {
				haveTopic = true
				break
			}
		}
		if !haveTopic {
			return false
		}
	}

	if len(f.InternalFuncs) == 0 {
		return true
	}

	itxIndex := 0
	for _, interFunc := range f.InternalFuncs {
		haveFunc := false
		for ; itxIndex < len(internalTxs); itxIndex++ {
			itx := internalTxs[itxIndex]
			if itx.CallingFunction != nil && *itx.CallingFunction == common.Bytes2Hex(interFunc.FuncSign) {
				haveFunc = true
				break
			}
		}
		if !haveFunc {
			return false
		}
	}

	return true
}

func getBorrowerAddress(name string, receipt *types.Receipt, internalTxs []*etypes.InternalTransaction) *common.Address {
	switch name {
	case aave.Name:
		for logIndex := 0; logIndex < len(receipt.Logs); logIndex++ {
			if len(receipt.Logs[logIndex].Topics) >= 2 &&
				receipt.Logs[logIndex].Topics[0] == aave.Events[0].Topic {
				topic := receipt.Logs[logIndex].Topics[1]
				addr := common.HexToAddress(topic.Hex())
				return &addr
			}
		}

	case dydx.Name:
		for logIndex := 0; logIndex < len(receipt.Logs); logIndex++ {
			if len(receipt.Logs[logIndex].Topics) != 0 &&
				receipt.Logs[logIndex].Topics[0] == dydx.Events[0].Topic {
				data := receipt.Logs[logIndex].Data
				addr := common.BytesToAddress(data)
				return &addr
			}
		}

	case uniswapV2.Name:
		interFunc := uniswapV2.InternalFuncs[0]
		for itxIndex := 0; itxIndex < len(internalTxs); itxIndex++ {
			itx := internalTxs[itxIndex]
			if itx.CallingFunction != nil && *itx.CallingFunction == common.Bytes2Hex(interFunc.FuncSign) {
				return &itx.To
			}
		}

	case uniswapV3.Name:

		for logIndex := 0; logIndex < len(receipt.Logs); logIndex++ {
			if len(receipt.Logs[logIndex].Topics) == 3 &&
				receipt.Logs[logIndex].Topics[0] == uniswapV3.Events[0].Topic {
				topic := receipt.Logs[logIndex].Topics[2]
				addr := common.HexToAddress(topic.Hex())
				return &addr
			}
		}
	}

	return nil
}

type Function struct {
	FuncName string
	FuncSign []byte
}

type Event struct {
	EventName string
	Topic     common.Hash
}
