package gophermart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/devldavydov/gophermart/internal/gophermart/storage"
	"github.com/sirupsen/logrus"
)

const (
	_httpClientTimeout    = 1 * time.Second
	_databaseOrderTimeout = 5 * time.Second
)

type AccrualStatus string

const (
	_accrualStatusRegistered AccrualStatus = "REGISTERED"
	_accrualStatusInvalid    AccrualStatus = "INVALID"
	_accrualStatusProcessing AccrualStatus = "PROCESSING"
	_accrualStatusProcessed  AccrualStatus = "PROCESSED"
)

type orderProcItem struct {
	order  string
	userID int
}

type AccrualDaemon struct {
	wg                   sync.WaitGroup
	accrualSystemAddress *url.URL
	stg                  storage.Storage
	maxThreadNum         int
	orderChan            chan orderProcItem
	dbScanInterval       time.Duration
	logger               *logrus.Logger
}

func NewAccrualDaemon(accrualSystemAddress *url.URL, stg storage.Storage, maxThreadNum int, dbScanInterval time.Duration, logger *logrus.Logger) *AccrualDaemon {
	return &AccrualDaemon{
		accrualSystemAddress: accrualSystemAddress,
		stg:                  stg,
		maxThreadNum:         maxThreadNum,
		orderChan:            make(chan orderProcItem, maxThreadNum*2),
		dbScanInterval:       dbScanInterval,
		logger:               logger,
	}
}

type accrualThread struct {
	httpClient           *http.Client
	accrualSystemAddress *url.URL
	threadID             int
	stg                  storage.Storage
	logger               *logrus.Logger
}

type accrualResponse struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual *float64      `json:"accrual,omitempty"`
}

func newAccrualThread(accrualSystemAddress *url.URL, threadID int, stg storage.Storage, logger *logrus.Logger) *accrualThread {
	client := &http.Client{
		Timeout: _httpClientTimeout,
	}

	return &accrualThread{
		httpClient:           client,
		accrualSystemAddress: accrualSystemAddress,
		threadID:             threadID,
		stg:                  stg,
		logger:               logger,
	}
}

func (ad *AccrualDaemon) Start(ctx context.Context) {
	ad.logger.Info("Accrual Daemon service started")

	var wg sync.WaitGroup

	ticker := time.NewTicker(ad.dbScanInterval)
	defer ticker.Stop()

	for i := 0; i < ad.maxThreadNum; i++ {
		wg.Add(1)
		go func(threadID int) {
			wg.Done()
			newAccrualThread(ad.accrualSystemAddress, threadID, ad.stg, ad.logger).start(ctx, ad.orderChan)
		}(i + 1)
	}

	for {
		select {
		case <-ticker.C:
			orders, err := ad.getOrdersToProcess()
			if err != nil {
				ad.logger.Errorf("failed to get orders to process: %s", err)
				continue
			}
			for _, i := range orders {
				ad.orderChan <- i
			}
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

func (ad *AccrualDaemon) getOrdersToProcess() ([]orderProcItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseOrderTimeout)
	defer cancel()

	orders, err := ad.stg.GetOrdersToProcess(ctx)
	if err != nil {
		return nil, err
	}

	orderNums := make([]orderProcItem, 0, len(orders))
	for _, v := range orders {
		orderNums = append(orderNums, orderProcItem{order: v.Number, userID: v.UserID})
	}
	return orderNums, nil
}

func (at *accrualThread) start(ctx context.Context, orderChan chan orderProcItem) {
	at.logger.Infof("[adthread #%d] started", at.threadID)
	for {
		select {
		case item := <-orderChan:
			at.logger.Infof("[adthread #%d] get order for process: %s", at.threadID, item.order)
			if err := at.processOrder(item.order, item.userID); err != nil {
				at.logger.Errorf("[adthread #%d] failed to process order [%s]: %s", at.threadID, item.order, err)
			}
		case <-ctx.Done():
			at.logger.Infof("[adthread #%d] finished", at.threadID)
			return
		}
	}
}

func (at *accrualThread) processOrder(orderNum string, userID int) error {
	ctxProcess, cancelProcess := context.WithTimeout(context.Background(), _databaseOrderTimeout)
	defer cancelProcess()

	err := at.stg.ProcessOrder(ctxProcess, orderNum)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodGet,
		at.accrualSystemAddress.JoinPath("api").JoinPath("orders").JoinPath(orderNum).String(),
		nil)
	if err != nil {
		return err
	}

	response, err := at.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("wrong response from accrual: %s", accrualResponseToString(response.StatusCode))
	}

	var accResp accrualResponse
	err = json.NewDecoder(response.Body).Decode(&accResp)
	if err != nil {
		return fmt.Errorf("failed to parse accrual response: %w", err)
	}

	ctxFinish, cancelFinish := context.WithTimeout(context.Background(), _databaseOrderTimeout)
	defer cancelFinish()

	switch accResp.Status {
	case _accrualStatusInvalid:
		err = at.stg.FinishOrder(ctxFinish, orderNum, userID, false, 0)
	case _accrualStatusProcessed:
		err = at.stg.FinishOrder(ctxFinish, orderNum, userID, true, *accResp.Accrual)
	default:
		at.logger.Infof("[adthread #%d] order %s still in process", at.threadID, orderNum)
	}

	if err != nil {
		return fmt.Errorf("failed to set accrual status: %w", err)
	}

	return nil
}

func accrualResponseToString(statusCode int) string {
	switch statusCode {
	case 204:
		return "not existing order"
	default:
		return http.StatusText(statusCode)
	}
}
