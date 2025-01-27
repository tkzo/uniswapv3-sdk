package helper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//

func SendTX(client *ethclient.Client, toAddress common.Address, value *big.Int,
	data []byte, w *Wallet) (*types.Transaction, error) {
	signedTx, err := TryTX(client, toAddress, value, data, w)
	if err != nil {
		return nil, err
	}
	return signedTx, client.SendTransaction(context.Background(), signedTx)
}

func TryTX(client *ethclient.Client, toAddress common.Address, value *big.Int,
	data []byte, w *Wallet) (*types.Transaction, error) {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     w.PublicKey,
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}
	gasLimit = gasLimit * 2
	nounc, err := client.NonceAt(context.Background(), w.PublicKey, nil)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nounc, toAddress, value,
		gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.PrivateKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
