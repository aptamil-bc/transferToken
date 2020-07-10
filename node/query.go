package node

import (
	"fmt"

	"github.com/shinecloudfoundation/shinecloudnet/codec"
	sdk "github.com/shinecloudfoundation/shinecloudnet/types"
	"github.com/shinecloudfoundation/shinecloudnet/x/auth/exported"
	authtypes "github.com/shinecloudfoundation/shinecloudnet/x/auth/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)



func (node *Node)GetAccountWithHeight(cdc *codec.Codec, addr sdk.AccAddress) (exported.Account, error) {
	bs, err := cdc.MarshalJSON(authtypes.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	res, err := node.query(fmt.Sprintf("custom/%s/%s", "acc", "account"), bs)
	if err != nil {
		return nil, err
	}

	var account exported.Account
	if err := cdc.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func (node *Node)query(path string, key cmn.HexBytes) (res []byte, err error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	}

	result, err := node.rpc.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, fmt.Errorf(resp.Log)
	}

	return resp.Value, nil
}
