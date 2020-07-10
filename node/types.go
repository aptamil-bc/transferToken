package node

import (
	"github.com/shinecloudfoundation/shinecloudnet/codec"
	sdk "github.com/shinecloudfoundation/shinecloudnet/types"
	"github.com/shinecloudfoundation/shinecloudnet/x/auth"
	authcutils "github.com/shinecloudfoundation/shinecloudnet/x/auth/client/utils"
	"github.com/shinecloudfoundation/shinecloudnet/x/auth/exported"
	"github.com/shinecloudfoundation/transfertoken/key"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Node struct {
	chainID    string
	rpc        rpcclient.Client
	cdc        *codec.Codec
	keyManager key.KeyManager
	txEncoder  sdk.TxEncoder
}

func NewNode(chainID, url string, keyManager key.KeyManager, cdc *codec.Codec) *Node {
	rpc := rpcclient.NewHTTP(url, "/websocket")
	return &Node{
		chainID:    chainID,
		rpc:        rpc,
		cdc:        cdc,
		keyManager: keyManager,
		txEncoder:  authcutils.GetTxEncoder(cdc),
	}
}

func (node *Node) BuildAndSign(account exported.Account, memo string, msgs []sdk.Msg, fees auth.StdFee) ([]byte, error) {
	stdSignMsg := auth.StdSignMsg{
		ChainID:       node.chainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
		Memo:          memo,
		Msgs:          msgs,
		Fee:           fees,
	}

	sigBytes, err := node.keyManager.Sign(stdSignMsg)
	if err != nil {
		return nil, err
	}
	stdSignature := auth.StdSignature{
		PubKey:    node.keyManager.GetPrivKey().PubKey(),
		Signature: sigBytes,
	}
	return node.txEncoder(auth.NewStdTx(msgs, fees, []auth.StdSignature{stdSignature}, memo))
}
