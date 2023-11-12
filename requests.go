package tehnomir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	return &Client{cfg: cfg, client: http.DefaultClient}
}

func (c *Client) newRequest(path ApiPath, body ...any) (*http.Response, error) {
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
	return c.client.Do(req)
}

func (c *Client) makeApiPath(path ApiPath) string {
	u := url.URL{
		Scheme: c.cfg.Proto,
		Host:   c.cfg.Host,
		Path:   string(path),
	}
	return u.String()
}

func (c *Client) makeRequestBody(path ApiPath, body any) any {
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
	}
	return body
}

func makeError(body io.ReadCloser) error {
	var resp ResponseErrorBody
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
	resp, err := c.newRequest(TestConnect, &TestRequestBody{Phrase: phrase})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return makeError(resp.Body)
	}

	return nil
}

func (c *Client) priceSearch(code string, currency Currency, brand int, analog bool) (*PriceSearchResponse, error) {
	resp, err := c.newRequest(
		PriceSearch,
		&PriceSearchRequestBody{
			BrandID:     brand,
			Code:        utilits.ClearString(code),
			ShowAnalogs: utilits.BoolToInt(analog),
			Currency:    string(currency),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res PriceSearchResponse
	if err := json.Unmarshal(body, &res); err != nil {
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
	resp, err := c.newRequest(
		GetSuppliers,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res SuppliersResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBrands() (*BrandsResponse, error) {
	resp, err := c.newRequest(
		GetBrands,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BrandsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBrandGroups() (*BrandGroupsResponse, error) {
	resp, err := c.newRequest(
		GetBrandGroups,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BrandGroupsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetProductInfo(code string, brand int) (*ProductInfoResponse, error) {
	resp, err := c.newRequest(
		GetProductInfo,
		&ProductInfoRequestBody{
			BrandID: brand,
			Code:    code,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res ProductInfoResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetUnloads(from, to time.Time) (*UnloadsResponse, error) {
	timeFormat := "2006-02-01"
	resp, err := c.newRequest(
		GetUnloads,
		&GetUnloadsRequestBody{
			FromDate: from.Format(timeFormat),
			ToDate:   to.Format(timeFormat),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res UnloadsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetUnloadData(unloadID int) (*UnloadResponse, error) {
	resp, err := c.newRequest(
		GetUnloadData,
		&GetUnloadDataRequestBody{
			UnloadID: unloadID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res UnloadResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBoxesReady() (*BoxesReadyToSendResponse, error) {
	resp, err := c.newRequest(
		GetBoxesReady,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BoxesReadyToSendResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) BasketAdd(prodid int64, priceLogo string, quantity int, reference string, comment ...string) (*BasketAddResponse, error) {
	var com string
	if comment != nil {
		com = comment[0]
	}
	resp, err := c.newRequest(
		BasketAdd,
		&BasketAddRequestBody{
			ProductID: prodid,
			PriceLogo: priceLogo,
			Quantity:  quantity,
			Reference: reference,
			Comment:   com,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BasketAddResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetPositionInfo(reference string) (*PositionInfoResponse, error) {
	resp, err := c.newRequest(
		GetPositionInfo,
		&PositionInfoRequestBody{
			Reference: reference,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res PositionInfoResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBasketPositions() (*BasketPositionsResponse, error) {
	resp, err := c.newRequest(
		GetBasketPositions,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BasketPositionsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) BasketDeletePosition(basketid int) error {
	resp, err := c.newRequest(
		BasketDeletePosition,
		&BasketDeletePositionRequestBody{
			BasketID: basketid,
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res SuccessResponseBody
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("failed request")
	}

	return nil
}

func (c *Client) BasketClear() error {
	resp, err := c.newRequest(
		BasketClear,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res SuccessResponseBody
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("failed request")
	}

	return nil
}

func (c *Client) GetCurrencies() (*CurrenciesResponse, error) {
	resp, err := c.newRequest(
		GetCurrencies,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res CurrenciesResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) BrandsByCode(code string) (*BrandsByCodeResponse, error) {
	resp, err := c.newRequest(
		GetBrandsByCode,
		&BrandsByCodeRequestBody{
			Code: code,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res BrandsByCodeResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) PositionStatuses() (*PositionStatusesResponse, error) {
	resp, err := c.newRequest(
		GetPositionStatuses,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, makeError(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res PositionStatusesResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}