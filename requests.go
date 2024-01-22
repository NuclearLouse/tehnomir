package tehnomir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/NuclearLouse/tehnomir/utilits"
)

type Client struct {
	cfg    *Config
	client *http.Client
}

func New(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
				TLSHandshakeTimeout:   time.Second,
				ResponseHeaderTimeout: time.Second * 3,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       60 * time.Second,
			},
		},
	}
}

func (c *Client) newRequest(path apiPath, body ...any) (*http.Response, error) {
	var reqbody any
	if body == nil {
		reqbody = TokenRequestBody{
			Token: c.cfg.Token,
		}
	} else {
		reqbody = c.makeRequestBody(path, body[0])
	}
	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(reqbody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.makeApiPath(path),
		buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}
	return res, nil
}

func (c *Client) makeApiPath(path apiPath) string {
	u := url.URL{
		Scheme: c.cfg.Proto,
		Host:   c.cfg.Host,
		Path:   string(path),
	}
	return u.String()
}

func (c *Client) makeRequestBody(path apiPath, body any) any {
	switch path {
	case TestConnect:
		b, ok := body.(*TestRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case PriceSearch:
		b, ok := body.(*PriceSearchRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetProductInfo:
		b, ok := body.(*ProductInfoRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetUnloads:
		b, ok := body.(*GetUnloadsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetUnloadData:
		b, ok := body.(*GetUnloadDataRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case BasketAdd:
		b, ok := body.(*BasketAddRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetPositionInfo:
		b, ok := body.(*PositionInfoRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case BasketDeletePosition:
		b, ok := body.(*BasketDeletePositionRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case OrderCreate:
		b, ok := body.(*OrderCreateRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case OrderSearch:
		b, ok := body.(*OrderSearchRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetChangedPositions:
		b, ok := body.(*GetChangedPositionsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetOrderPositions:
		b, ok := body.(*GetOrderPositionsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetOrderPositionsByStatus:
		b, ok := body.(*GetOrderPositionsByStatusRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	}
	return body
}

func (c *Client) requestAndDecode(request apiPath, response any, requestbody ...any) error {
	var (
		body any
		err  error
		resp *http.Response
	)
	if requestbody != nil {
		body = requestbody[0]
		resp, err = c.newRequest(
			request,
			body,
		)
	} else {
		resp, err = c.newRequest(request)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func makeError(body io.ReadCloser) error {
	var resp ResponseError
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return err
	}
	return fmt.Errorf("%d:%s - %s", resp.Data.Status, resp.Data.Name, resp.Data.Message)
}

func (c *Client) TestConnect(s ...string) error {
	var phrase string
	if s != nil {
		phrase = s[0]
	}
	var res TestConnectResponse
	if err := c.requestAndDecode(TestConnect, &res, &TestRequestBody{Phrase: phrase}); err != nil {
		return err
	}
	if s != nil && res.Data.TestString != s[0] {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) priceSearch(code string, currency Currency, brand int, analog bool) (*PriceSearchResponse, error) {
	var res PriceSearchResponse
	if err := c.requestAndDecode(PriceSearch, &res,
		&PriceSearchRequestBody{
			BrandID:     brand,
			Code:        utilits.ClearString(code),
			ShowAnalogs: utilits.BoolToInt(analog),
			Currency:    string(currency),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) SearchWithAnalogs(code string, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(code, cur, 0, true)
}

func (c *Client) SearchByBrandWithAnalogs(code string, brand int, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(code, cur, brand, true)
}

func (c *Client) SearchByBrandWithoutAnalogs(code string, brand int, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(code, cur, brand, false)
}

func (c *Client) GetSuppliers() (*SuppliersResponse, error) {
	var res SuppliersResponse
	if err := c.requestAndDecode(GetSuppliers, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetBrands() (*BrandsResponse, error) {
	var res BrandsResponse
	if err := c.requestAndDecode(GetBrands, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBrandGroups() (*BrandGroupsResponse, error) {
	var res BrandGroupsResponse
	if err := c.requestAndDecode(GetBrandGroups, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetProductInfo(code string, brand int) (*ProductInfoResponse, error) {
	var res ProductInfoResponse
	if err := c.requestAndDecode(GetProductInfo, &res, &ProductInfoRequestBody{
		BrandID: brand,
		Code:    code,
	}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetUnloads(from, to time.Time) (*UnloadsResponse, error) {
	timeFormat := "2006-01-02"
	var res UnloadsResponse
	if err := c.requestAndDecode(GetUnloads, &res, &GetUnloadsRequestBody{
		FromDate: from.Format(timeFormat),
		ToDate:   to.Format(timeFormat),
	}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetUnloadData(unloadID int) (*UnloadResponse, error) {
	var res UnloadResponse
	if err := c.requestAndDecode(GetUnloadData, &res,
		&GetUnloadDataRequestBody{
			UnloadID: unloadID,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetBoxesReady() (*BoxesReadyToSendResponse, error) {
	var res BoxesReadyToSendResponse
	if err := c.requestAndDecode(GetBoxesReady, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BasketAdd(prodid int64, priceLogo string, quantity int, reference string, comment ...string) (*BasketAddResponse, error) {
	var com string
	if comment != nil {
		com = comment[0]
	}
	var res BasketAddResponse
	if err := c.requestAndDecode(BasketAdd, &res,
		&BasketAddRequestBody{
			ProductID: prodid,
			PriceLogo: priceLogo,
			Quantity:  quantity,
			Reference: reference,
			Comment:   com,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetPositionInfo(reference string) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetPositionInfo, &res,
		&PositionInfoRequestBody{
			Reference: reference,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetBasketPositions() (*BasketPositionsResponse, error) {
	var res BasketPositionsResponse
	if err := c.requestAndDecode(GetBasketPositions, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BasketDeletePosition(basketid int) error {
	var res SuccessResponse
	if err := c.requestAndDecode(BasketDeletePosition, &res,
		&BasketDeletePositionRequestBody{
			BasketID: basketid,
		}); err != nil {
		return err
	}
	if !res.Success {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) BasketClear() error {
	var res SuccessResponse
	if err := c.requestAndDecode(BasketClear, &res); err != nil {
		return err
	}
	if !res.Success {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) GetCurrencies() (*CurrenciesResponse, error) {
	var res CurrenciesResponse
	if err := c.requestAndDecode(GetCurrencies, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BrandsByCode(code string) (*BrandsByCodeResponse, error) {
	var res BrandsByCodeResponse
	if err := c.requestAndDecode(GetBrandsByCode, &res,
		&BrandsByCodeRequestBody{
			Code: code,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) PositionStatuses() (*PositionStatusesResponse, error) {
	var res PositionStatusesResponse
	if err := c.requestAndDecode(GetPositionStatuses, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderCreate(ordernum string) (*OrderResponse, error) {
	var res OrderResponse
	if err := c.requestAndDecode(OrderCreate, &res,
		&OrderCreateRequestBody{
			OrderNumber: ordernum,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ActiveOrders() (*OrdersResponse, error) {
	var res OrdersResponse
	if err := c.requestAndDecode(GetActiveOrders, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderSearchByDate(fromDate, toDate time.Time) (*OrdersResponse, error) {
	timeFormat := "2006-01-02"
	var res OrdersResponse
	if err := c.requestAndDecode(OrderSearch, &res,
		&OrderSearchRequestBody{
			FromDate: fromDate.Format(timeFormat),
			ToDate:   toDate.Format(timeFormat),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderSearchByNumber(ordernum string) (*OrdersResponse, error) {
	var res OrdersResponse
	if err := c.requestAndDecode(OrderSearch, &res,
		&OrderSearchRequestBody{
			OrderNum: ordernum,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ChangedPositions(fromDateTime time.Time) (*PositionsInfoResponse, error) {
	timeFormat := "2006-01-02 15:04:05"
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetChangedPositions, &res,
		&GetChangedPositionsRequestBody{
			FromDate: fromDateTime.Format(timeFormat),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderPositions(orderid int) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetOrderPositions, &res,
		&GetOrderPositionsRequestBody{
			OrderID: orderid,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

// Отличаются статусы: В работе - StatusID: 2 и Закрыт - StatusID: 4
func (c *Client) OrderPositionsByStatus(statusid int) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetOrderPositionsByStatus, &res,
		&GetOrderPositionsByStatusRequestBody{
			StatusID: statusid,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) StockPrice() (*StockPriceResponse, error) {
	var res StockPriceResponse
	if err := c.requestAndDecode(GetStockPrice, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
