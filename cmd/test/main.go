package main

import (
	"fmt"

	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
)

var (
	nodeUrl     string
	walletOwner string
	seed        string
	PrivateKey  string

	PublicKey string
)

func main() {
	nodeUrl = "http://62.182.156.133:8843/v2/"
	walletOwner = "Mx8e37e797422abb977914e2878e14c8a3a67ba5a1"
	seed = "argue spider slab admit cheese local cherry cool crane sea fish term"
	PublicKey = "Mpbf64d15c692dee2914b941041a86b385f9d7b0b6ab2d2debe2fd86e362eba0e4"

	seed1, err := wallet.Seed(seed)
	if err != nil {
		panic(err)
	}
	prKey, err := wallet.PrivateKeyBySeed(seed1)
	if err != nil {
		panic(err)
	}
	PrivateKey = prKey

	SendTransactionCandidateOffTest()
}

func SendTransactionCandidateOffTest() error {
	client, _ := http_client.NewConcise(nodeUrl)
	nonce, _ := client.Nonce(walletOwner)
	tx, _ := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(
		transaction.
			NewSetCandidateOffData().
			MustSetPubKey(PublicKey),
	)
	sign, _ := tx.
		SetNonce(nonce).
		SetGasPrice(1).
		SetGasCoin(0).
		Sign(PrivateKey)
	encode, _ := sign.Encode()
	hash, _ := sign.Hash()
	_, err := client.SendTransaction(encode)
	if err != nil {
		_, body, err := http_client.ErrorBody(err)
		json, _ := http_client.Marshal(body)
		fmt.Println(json)
		return err
	}

	fmt.Println(hash)
	return nil
}
func SendTransactionCandidateOnTest() error {
	client, _ := http_client.NewConcise(nodeUrl)
	nonce, _ := client.Nonce(walletOwner)
	tx, _ := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(
		transaction.
			NewSetCandidateOnData().
			MustSetPubKey(PublicKey),
	)
	sign, _ := tx.
		SetNonce(nonce).
		SetGasPrice(1).
		SetGasCoin(0).
		Sign(PrivateKey)
	encode, _ := sign.Encode()
	hash, _ := sign.Hash()
	_, err := client.SendTransaction(encode)
	if err != nil {
		_, body, err := http_client.ErrorBody(err)
		json, _ := http_client.Marshal(body)
		fmt.Println(json)
		return err
	}

	fmt.Println(hash)
	return nil
}
