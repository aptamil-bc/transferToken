package main

import (
	"fmt"

	"github.com/shinecloudfoundation/shinecloudnet/app"
	sdk "github.com/shinecloudfoundation/shinecloudnet/types"
	"github.com/shinecloudfoundation/shinecloudnet/x/auth"
	"github.com/shinecloudfoundation/shinecloudnet/x/bank/types"
	"github.com/shinecloudfoundation/transfertoken/key"
	"github.com/shinecloudfoundation/transfertoken/node"
)

func main() {
	cdc := app.MakeCodec()

	keyManager, err := key.NewMnemonicKeyManager("filter cancel soul illness treat step input virus region item garage poet")
	if err != nil {
		fmt.Println(err)
	}

	chainID := "shinecloudnet-test"
	rpcNode := node.NewNode(chainID, "http://3.115.116.139:26657", keyManager, cdc)

	coins, err := sdk.ParseCoins("100uscds")
	if err != nil {
		fmt.Println(err)
		return
	}
	account, err := rpcNode.GetAccountWithHeight(cdc, keyManager.GetAddr())
	if err != nil {
		fmt.Println(err)
		return
	}
	to, err := sdk.AccAddressFromBech32("scloud1ksl2ct28stu43qwl388nxa3wtypju5xlk9u3np")
	if err != nil {
		fmt.Println(err)
		return
	}
	msgs := []sdk.Msg{types.NewMsgSend(keyManager.GetAddr(), to, coins)}
	memo := ""
	fees := auth.NewStdFee(200000, sdk.Coins{sdk.NewCoin("uscds", sdk.NewInt(1e6))})
	signedTxBytes, err := rpcNode.BuildAndSign(account, memo, msgs, fees)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := rpcNode.Broadcast(signedTxBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
