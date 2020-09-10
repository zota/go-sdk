package zotapay

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// OrdersReport represents order report request
type OrdersReport struct {
	MerchantID  string `url:"merchantID"`
	DateType    string `url:"dateType"`
	EndpointIds string `url:"endpointIds"`
	FromDate    string `url:"fromDate"`
	RequestID   string `url:"requestID"`
	Statuses    string `url:"statuses"`
	Timestamp   string `url:"timestamp"`
	ToDate      string `url:"toDate"`
	Types       string `url:"types"`
	Signature   string `url:"signature"`
}

// OrdersReportResult represents the response from Zotapay API
type OrdersReportResult struct {
	Code         string `json:"code"`
	OrdersReport string `json:"data"`
	Message      string `json:"message"`
}

var mockedOrdersReportResult *OrdersReportResult

// OrdersReport init validation of the SDK struct and the OrdersReport
// generate sign and
// init a orders report request to Zotapay API
func (s *SDK) OrdersReport(d OrdersReport) (res OrdersReportResult, err error) {

	//validate that SDK is properly initialized
	err = s.validate()
	if err != nil {
		return
	}

	//validate that OrdersReport is properly initialized
	err = d.validate()
	if err != nil {
		return
	}

	//set Merchant id
	d.MerchantID = s.MerchantID

	//set current timestamp
	d.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	//generate new UUID
	d.RequestID = uuid.New().String()

	//if mockedOrdersReportResult is set return it as response
	//only for testing
	if mockedOrdersReportResult != nil {
		res = *mockedOrdersReportResult
		mockedOrdersReportResult = nil
		return
	}

	//generate signature
	d.Signature = s.sign(d.MerchantID, d.DateType, d.EndpointIds, d.FromDate, d.RequestID, d.Statuses, d.Timestamp, d.ToDate, d.Types)

	v, err := query.Values(d)
	if err != nil {
		return
	}

	code, body, err := s.httpDo(http.MethodGet, fmt.Sprintf("%v/api/v1/query/orders-report/csv/?%v", s.ApiBaseURL, v.Encode()), []byte(""))
	if err != nil {
		return
	}

	if code != 200 {
		err = json.Unmarshal(body, &res)
		if err != nil {
			err = fmt.Errorf("json Unmarshal err:%v", err)
			return
		}
		return
	}
	res = OrdersReportResult{
		Code:         "200",
		OrdersReport: string(body),
	}
	return
}

// SetMockResponse will set *OrdersReportResult as mocked
// response on the next order report request
// To be used only for test purposes
func (mock *OrdersReportResult) SetMockResponse() {
	mockedOrdersReportResult = mock
}

// validate the instance of OrdersReport
// if not valid returns an error
func (d *OrdersReport) validate() error {
	required := []string{"DateType", "FromDate", "ToDate"}
	for _, fieldName := range required {
		r := reflect.ValueOf(d)
		f := reflect.Indirect(r).FieldByName(fieldName)
		value := f.String()
		if value == "" {
			return fmt.Errorf("%v is required", fieldName)
		}
	}
	return nil
}
