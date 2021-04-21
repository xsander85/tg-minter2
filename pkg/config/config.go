package config

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"time"
)

type Config struct {
	Minter   Minter       `json:"minter"`
	Telegram TelegramConf `json:"telegram"`
}

type TelegramConf struct {
	Token  string  `json:"token"`
	ChatId []int64 `json:"chatId"`
	Debug  bool    `json:"debug"`
}

type Minter struct {
	NodeList  []NodeItem `json:"nodelist"`
	Validator Validator  `json:"validator"`
}
type Validator struct {
	Testnet          bool   `json:"testnet"`
	OwnerSeed        string `json:"ownerSeed"`
	OwnerWallet      string `json:"ownerWallet"`
	PublicKey        string `json:"publicKey"`
	MaxMissedBlockes uint   `json:"maxMissedBlockes"`
}
type NodeItem struct {
	Url     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
}

func (c *Config) LoadConfig(configFile string) error {
	jsonFile, err := os.Open(configFile)

	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)

	return nil
}

//Generates random int as function of range
func GetRandomInt(Range int) int {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(Range)
}

func In_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}
