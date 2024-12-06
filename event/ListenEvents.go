package event

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/haitao-sun03/go/abi/erc20"
	"github.com/haitao-sun03/go/config"
	log "github.com/sirupsen/logrus"
)

func ListenTransferEvent() {
	contract, _ := erc20.NewERC20(common.HexToAddress(config.Config.Geth.ContractAddress), config.GethWsClient)
	transferChan := make(chan *erc20.ERC20Transfer)
	myEventSubscription, err := contract.WatchTransfer(&bind.WatchOpts{Start: nil, Context: context.Background()}, transferChan, nil, nil)
	if err != nil {
		log.WithError(err).Error("WatchTransfer err")
	}
	defer myEventSubscription.Unsubscribe()
	// 处理事件
	for {

		select {
		// 处理错误
		case err := <-myEventSubscription.Err():
			log.Fatal(err)
		case event := <-transferChan:
			log.Infof("Transfer event: %+v", event)

		}
	}

}

func ListenEvents() {
	ListenTransferEvent()
}
