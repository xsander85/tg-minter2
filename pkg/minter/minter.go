package minter

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/xsander85/tg-minter2/pkg/config"
)

type Valid struct {
	last_status     uint64
	last_block      uint64
	validator       config.Validator
	clientService   map[int]http_client.Client
	privateKey      string
	chainId         transaction.ChainID
	shutdownChannel chan interface{}
}

func New(c *config.Config) *Valid {

	clients := make(map[int]http_client.Client)

	for index, nodeItem := range c.Minter.NodeList {
		client, _ := http_client.NewConcise(nodeItem.Url)

		if nodeItem.Headers != nil {
			client = client.WithHeaders(nodeItem.Headers)

		}
		clients[index] = *client

	}
	seed1, err := wallet.Seed(c.Minter.Validator.OwnerSeed)
	if err != nil {
		panic(err)
	}
	prKey, err := wallet.PrivateKeyBySeed(seed1)
	if err != nil {
		panic(err)
	}

	chainID := transaction.TestNetChainID
	if !c.Minter.Validator.Testnet {
		chainID = transaction.MainNetChainID
	}

	return &Valid{
		last_status:     0,
		last_block:      0,
		validator:       c.Minter.Validator,
		clientService:   clients,
		privateKey:      prKey,
		chainId:         chainID,
		shutdownChannel: make(chan interface{}),
	}

}

func (v *Valid) Run(с chan string) {
	fmt.Println("Run()")
	go func() {
		for {
			select {
			case <-v.shutdownChannel:
				return
			default:
			}
			v.handleCheck(с)
			time.Sleep(time.Second)
		}
	}()
}

func (v *Valid) Exit() {
	close(v.shutdownChannel)
	fmt.Println("Valid stop")
}

func (v *Valid) handleCheck(сh1 chan string) {
	randNodeIndex := config.GetRandomInt(len(v.clientService))

	lastBlock := v.GetLastBlock(0)
	if v.last_block != lastBlock {
		fmt.Printf("lastblock!= [%d]   ", lastBlock)
		v.last_block = lastBlock
		v_data := v.GetValidatorData(randNodeIndex)
		fmt.Printf("#NodeIndex [%d]  status_new=[%d] status_old [%d]\n", randNodeIndex, v_data.Status, v.last_status)
		if v.last_status != v_data.Status {
			if v.last_status != 0 {
				status := "*ВЫКЛЮЧЕН*"
				if v.last_status == 2 {
					status = "*Включен*"
				}
				message := fmt.Sprintf("Валидатор [%s] изменил статус на [%s]", v.validator.PublicKey, status)
				fmt.Printf(">---------- %s\n", message)
				сh1 <- message
			} else {
				message := fmt.Sprintf("Запущена система мониторинга на блоке [%d]", v.last_block)
				fmt.Printf(">---------- %s\n", message)
				сh1 <- message

				message = v.StatusValidator(v_data)
				fmt.Printf(">---------- %s\n", message)
				сh1 <- message
			}
			v.last_status = v_data.Status
		}
		if v.last_status == 2 {
			if v.isMaxLimitErrors(randNodeIndex) {
				message := fmt.Sprintf("*Валидатор пропустил максимально допустимое количество блоков!*\n [%d]", v.last_block)
				сh1 <- message
				v.SendTransactionOff(randNodeIndex)
			}
		}
	}
}

func (v *Valid) GetValidatorData(nodeItemId int) *models.CandidateResponse {

	client := v.clientService[nodeItemId]
	data, err := client.Candidate(v.validator.PublicKey)
	if err != nil {
		panic(err)
	}
	return data
}

// Last block height
func (v *Valid) GetLastBlock(nodeItemId int) uint64 {
	client := v.clientService[nodeItemId]
	data, err := client.Status()

	if err != nil {
		panic(err)
	}
	return data.LatestBlockHeight
}

// Send Transaction set Validator Off
func (v *Valid) SendTransactionOff(nodeItemId int) string {
	client := v.clientService[nodeItemId]
	nonce, _ := client.Nonce(v.validator.OwnerWallet)

	tx, _ := transaction.
		NewBuilder(v.chainId).
		NewTransaction(
			transaction.
				NewSetCandidateOffData().
				MustSetPubKey(v.validator.PublicKey),
		)
	sign, _ := tx.
		SetNonce(nonce).
		SetGasPrice(1).
		SetGasCoin(0).
		Sign(v.privateKey)
	encode, _ := sign.Encode()
	hash, _ := sign.Hash()
	_, err := client.SendTransaction(encode)
	if err != nil {
		_, body, _ := http_client.ErrorBody(err)
		json, _ := http_client.Marshal(body)
		fmt.Println(json)
		return json
	}

	return hash
}

// Send Transaction set Validator On
func (v *Valid) SendTransactionOn(nodeItemId int) string {
	client := v.clientService[nodeItemId]
	nonce, _ := client.Nonce(v.validator.OwnerWallet)
	tx, _ := transaction.
		NewBuilder(v.chainId).
		NewTransaction(
			transaction.
				NewSetCandidateOnData().
				MustSetPubKey(v.validator.PublicKey),
		)
	sign, _ := tx.
		SetNonce(nonce).
		SetGasPrice(1).
		SetGasCoin(0).
		Sign(v.privateKey)
	encode, _ := sign.Encode()
	hash, _ := sign.Hash()
	_, err := client.SendTransaction(encode)
	if err != nil {
		_, body, _ := http_client.ErrorBody(err)
		json, _ := http_client.Marshal(body)
		fmt.Println(json)
		return json
	}

	return hash
}

func (v *Valid) StatusValidator(data *models.CandidateResponse) string {
	totalStake, _ := strconv.Atoi(data.TotalStake)
	totalStakeBip := uint64(totalStake / 10000000)

	return fmt.Sprintf(
		"*Валидатор*:\n _%s_\n"+
			"*Размер комиссии*: _%d_\n"+
			"*Адрес управляющего*:\n _%s_\n"+
			"*Адрес владельца*:\n _%s_\n"+
			"*Минимальный стейк*: _%s_\n"+
			"*Общий стейк*:\n _%d_\n"+
			"*Уникальных пользователей*: _%d_\n"+
			"*Количество занятых слотов*: _%d_\n",
		v.validator.PublicKey,
		data.Commission,
		data.ControlAddress,
		data.OwnerAddress,
		data.MinStake,
		totalStakeBip,
		data.UniqUsers,
		data.UsedSlots)

}

func (v *Valid) isMaxLimitErrors(nodeItemId int) bool {
	count := v.GetMissedBlockes(nodeItemId)

	return count >= int64(v.validator.MaxMissedBlockes)
}

func (v *Valid) GetMissedBlockes(nodeItemId int) int64 {
	client := v.clientService[nodeItemId]

	data, err := client.MissedBlocks(v.validator.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	return data.MissedBlocksCount

}
