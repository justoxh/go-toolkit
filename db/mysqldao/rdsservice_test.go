package mysqldao

import (
	"testing"
)

type Order struct {
	ID        int    `gorm:"column:id;primary_key;"`
	OrderHash string `gorm:"column:order_hash;type:varchar(82)"`
	Market    string `gorm:"column:market;type:varchar(40)"`
	Side      string `gorm:"column:side;type:varchar(40)"`
	OrderType string `gorm:"column:order_type;type:varchar(40)"`
}

func newTestRdsService(t *testing.T) (*RdsService, error) {

	config := MysqlConifg{
		Hostname:    "127.0.0.1",
		Port:        "3306",
		User:        "root",
		Password:    "root",
		DbName:      "test",
		TablePrefix: "",
		Debug:       false,
	}

	rdsService, err := NewRdsService(config)
	if err != nil {
		t.Fatal(err)
	}

	return rdsService, err
}

func TestAdd(t *testing.T) {
	rdsService, err := newTestRdsService(t)
	if err != nil {
		t.Fatal(err)
	}

	rdsService.RegistTable(&Order{})

	order := Order{
		OrderHash: "0xeeec6ca7c19ff79562e98c30c5daa9fb88b377a1957c6d4acd48fe47c047c257",
		Market:    "SEELE-WETH",
		Side:      "sell",
		OrderType: "market_order",
	}

	err = rdsService.Add(&order)
	if err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	rdsService, err := newTestRdsService(t)
	if err != nil {
		t.Fatal(err)
	}

	rdsService.RegistTable(&Order{})

	order := Order{
		OrderHash: "0xeeec6ca7c19ff79562e98c30c5daa9fb88b377a1957c6d4acd48fe47c047c257",
		Market:    "SEELE-WETH",
		Side:      "sell",
		OrderType: "market_order",
	}

	err = rdsService.Add(&order)
	if err != nil {
		t.Error(err)
	}

	err = rdsService.Del(&order)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSave(t *testing.T) {
	rdsService, err := newTestRdsService(t)
	if err != nil {
		t.Fatal(err)
	}

	rdsService.RegistTable(&Order{})

	order := Order{
		OrderHash: "0xeeec6ca7c19ff79562e98c30c5daa9fb88b377a1957c6d4acd48fe47c047c257",
		Market:    "SEELE-WETH",
		Side:      "sell",
		OrderType: "market_order",
	}

	err = rdsService.Add(&order)
	if err != nil {
		t.Error(err)
	}

	order.Side = "buy"
	err = rdsService.Save(&order)
	if err != nil {
		t.Error(err)
	}
}
