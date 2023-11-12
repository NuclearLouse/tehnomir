package tehnomir

import (
	"time"

	"github.com/NuclearLouse/tehnomir/utilits"
)

const (
	EUR Currency = "EUR"
	USD Currency = "USD"

	TestConnect     ApiPath = "test/connect"
	PriceSearch     ApiPath = "price/search"
	GetSuppliers    ApiPath = "info/getSuppliers"
	GetBrands       ApiPath = "info/getBrands"
	GetBrandGroups  ApiPath = "info/getBrandGroups"
	GetProductInfo  ApiPath = "info/getProductInfo"
	GetUnloads      ApiPath = "unload/search"
	GetUnloadData   ApiPath = "unload/getData"
	GetBoxesReady   ApiPath = "unload/getBoxesReadyToSend"
	BasketAdd       ApiPath = "basket/add"
	GetPositionInfo ApiPath = "order/getPositionInfo"
)

type (
	Currency string
	ApiPath  string
)

type Config struct {
	Token          string        `cfg:"token"`
	Proto          string        `cfg:"proto"`
	Host           string        `cfg:"host"`
	Timeout        time.Duration `cfg:"timeout"`
	PriceAvia      float64       `cfg:"price_avia"`
	PriceContainer float64       `cfg:"price_container"`
	PriceVolume    float64       `cfg:"price_volume"`
}

type TokenRequestBody struct {
	Token string `json:"apiToken"`
}
type TestRequestBody struct {
	TokenRequestBody
	Phrase string `json:"string"`
}

type PriceSearchRequestBody struct {
	TokenRequestBody
	BrandID     int    `json:"brandId"`
	Code        string `json:"code"`
	ShowAnalogs int    `json:"isShowAnalogs"`
	Currency    string `json:"currency"`
}

type ProductInfoRequestBody struct {
	TokenRequestBody
	BrandID int    `json:"brandId"`
	Code    string `json:"code"`
}

type GetUnloadsRequestBody struct {
	TokenRequestBody
	FromDate string `json:"fromDate"` //"2023-11-08" for "unload/search"
	ToDate   string `json:"toDate"`   // for "unload/search"
}

type GetUnloadDataRequestBody struct {
	TokenRequestBody
	UnloadID int `json:"unloadId"` //for unload/getData
}

type BasketAddRequestBody struct {
	TokenRequestBody
	ProductID int64  `json:"productId"`
	PriceLogo string `json:"priceLogo"`
	Quantity  int    `json:"quantity"`
	Reference string `json:"reference"` //for order/getPositionInfo
	Comment   string `json:"comment"`
}

type PositionInfoRequestBody struct {
	TokenRequestBody
	Reference string `json:"reference"` //for order/getPositionInfo
}

type ResponseErrorBody struct {
	Success bool `json:"success"`
	Data    struct {
		Name       string `json:"name"`
		Status     int    `json:"status"`
		Message    string `json:"message"`
		TestString string `json:"testString,omitempty"`
	} `json:"data"`
}

// ответ на price/search
// можно использовать и для добавления детали в корзину пользователя
// проверка по PriceLogo
type PriceSearchResponse struct {
	Success bool          `json:"success"`
	Details []FoundDetail `json:"data"`
}

type FoundDetail struct {
	ProductID          int                   `json:"productId"`
	BrandID            int                   `json:"brandId"`
	BrandGroupID       int                   `json:"brandGroupId"`
	Brand              string                `json:"brand"`
	Code               string                `json:"code"`
	DescriptionRus     string                `json:"descriptionRus"`
	DescriptionUa      string                `json:"descriptionUa"`
	Weight             utilits.CustomFloat64 `json:"weight"`
	IsOriginal         utilits.CustomBool    `json:"isOriginal"`
	IsExistProductInfo utilits.CustomBool    `json:"isExistProductInfo"`
	Stocks             []OfferSupplier       `json:"rests"`
}

type OfferSupplier struct {
	PriceLogo       string              `json:"priceLogo"`
	Price           float64             `json:"price"`
	Currency        string              `json:"currency"`
	Quantity        utilits.CustomInt64 `json:"quantity"`
	QuantityType    string              `json:"quantityType"`
	Multiplicity    int                 `json:"multiplicity"`
	PriceQuality    float64             `json:"priceQuality"`
	DeliveryTypeID  int                 `json:"deliveryTypeId"`
	DeliveryType    string              `json:"deliveryType"`
	DeliveryDays    int                 `json:"deliveryTime"`
	DeliveryDate    utilits.CustomTime  `json:"deliveryDate"`
	DeliveryPercent int                 `json:"deliveryPercent"`
	PriceChangeDate utilits.CustomTime  `json:"priceChangeDate"` //дата-время последнего обновления прайса поставщика
	IsReturn        utilits.CustomBool  `json:"isReturn"`
	IsPriceFinal    utilits.CustomBool  `json:"isPriceFinal"`
}

// ответ на info/getProductInfo для получения изображений и списка аналогов
type ProductInfo struct {
	Success bool `json:"success"`
	Data    struct {
		ProductID      int                   `json:"productId"`
		Brand          string                `json:"brand"`
		Code           string                `json:"code"`
		DescriptionRus string                `json:"descriptionRus"`
		DescriptionUa  string                `json:"descriptionUa"`
		Weight         utilits.CustomFloat64 `json:"weight"`
		Volume         int                   `json:"volume"`
		Images         []struct {
			Image string `json:"image"`
		} `json:"images"`
		Properties []interface{} `json:"properties"`
		Analogs    []struct {
			ProductID      int                   `json:"productId"`
			Brand          string                `json:"brand"`
			Code           string                `json:"code"`
			DescriptionRus string                `json:"descriptionRus"`
			DescriptionUa  string                `json:"descriptionUa"`
			Weight         utilits.CustomFloat64 `json:"weight"`
			Volume         float64               `json:"volume"`
		} `json:"analogs"`
	} `json:"data"`
}

type SuppliersResponse struct {
	Success   bool       `json:"success"`
	Suppliers []Supplier `json:"data"`
}

type Supplier struct {
	PriceLogo      string              `json:"priceLogo"`
	DeliveryTypeID int                 `json:"deliveryTypeId"`
	DeliveryType   string              `json:"deliveryType"`
	DeliveryDays   int                 `json:"deliveryTime"`
	DeliveryHours  utilits.CustomInt64 `json:"deliveryTimeHours"`
	DeliveryDate   utilits.CustomTime  `json:"deliveryDate"`
	Region         string              `json:"region"`
	RegionEn       string              `json:"regionEn"`
	RegionUa       string              `json:"regionUa"`
	IsReturn       utilits.CustomBool  `json:"isReturnFlag"`
	IsPriceFinal   utilits.CustomBool  `json:"isPriceFinalFlag"`
}

type BrandGroupsResponse struct {
	Success     bool         `json:"success"`
	BrandGroups []BrandGroup `json:"data"`
}

type BrandGroup struct {
	GroupID   int                 `json:"groupId"`
	GroupName string              `json:"group"`
	BrandIds  []utilits.CustomInt `json:"brandIds"`
}

type BrandsResponse struct {
	Success bool    `json:"success"`
	Brands  []Brand `json:"data"`
}

type Brand struct {
	BrandID    int                `json:"brandId"`
	BrandName  string             `json:"brand"`
	IsOriginal utilits.CustomBool `json:"isOriginal"`
}

// Список отгрузок за указанный период
type UnloadsResponse struct {
	Success bool     `json:"success"`
	Unloads []Unload `json:"data"`
}

type Unload struct {
	UnloadID       int                   `json:"unloadId"`
	CreateTime     utilits.CustomTime    `json:"createTime"` //"2023-11-08 17:19:00.882505",
	BoxQuantity    int                   `json:"boxQuantity"`
	SumPositions   utilits.CustomFloat64 `json:"sumPositions"`
	SumWorks       int                   `json:"sumWorks"`
	SumDelivery    utilits.CustomFloat64 `json:"sumDelivery"`
	SumTotal       utilits.CustomFloat64 `json:"sumTotal"`
	Carrier        string                `json:"carrier"`
	CarrierWaybill interface{}           `json:"carrierWaybill"` // maybe null
}

// информация по отгрузке
type UnloadResponse struct {
	Success bool       `json:"success"`
	Unload  UnloadData `json:"data"`
}
type UnloadData struct {
	Boxes     []UnloadBox      `json:"boxes"`
	Positions []UnloadPosition `json:"positions"`
}

type UnloadBox struct {
	BoxID        int                   `json:"boxId"`
	SumPositions utilits.CustomFloat64 `json:"sumPositions"`
	SumWorks     utilits.CustomFloat64 `json:"sumWorks"` //стоимость доставки
	Length       int                   `json:"length"`
	Width        int                   `json:"width"`
	Height       int                   `json:"height"`
	Weight       utilits.CustomFloat64 `json:"weight"`
}

type UnloadPosition struct {
	BoxID           int     `json:"boxId"`
	OrderID         int     `json:"orderId"`     //номер заказа в ТМ
	OrderNumber     string  `json:"orderNumber"` //номер заказа в ТМ
	OrderPositionID int     `json:"orderPositionId"`
	PriceLogo       string  `json:"priceLogo"`
	Brand           string  `json:"brand"`
	BrandID         int     `json:"brandId"`
	Code            string  `json:"code"`
	DescriptionRus  string  `json:"descriptionRus"`
	DescriptionUa   string  `json:"descriptionUa"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	PriceFinal      float64 `json:"priceFinal"`
	Currency        string  `json:"currency"`
	Reference       string  `json:"reference"` //!отдает в верхнем регистре
	Comment         string  `json:"comment"`   //комментарий к клиенту
	AdminComment    string  `json:"adminComment"`
	Weight          float64 `json:"weight"`
	Sticker         string  `json:"sticker"`
}

// Получение коробок готовых к отгрузке
// возможно не понадобится
type BoxesReadyToSendResponse struct {
	Success    bool        `json:"success"`
	ReadyBoxes []UnloadBox `json:"data"`
}

type ProductInfoResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ProductID      int     `json:"productId"`
		Brand          string  `json:"brand"`          //Нужно
		Code           string  `json:"code"`           //Нужно
		DescriptionRus string  `json:"descriptionRus"` //Нужно
		DescriptionUa  string  `json:"descriptionUa"`
		Weight         float64 `json:"weight"`
		Volume         int     `json:"volume"`
		Images         []struct {
			Image string `json:"image"` //Нужно
		} `json:"images"`
		Properties []interface{} `json:"properties"`
		Analogs    []struct {
			ProductID      int     `json:"productId"`
			Brand          string  `json:"brand"`
			Code           string  `json:"code"`
			DescriptionRus string  `json:"descriptionRus"`
			DescriptionUa  string  `json:"descriptionUa"`
			Weight         float64 `json:"weight"`
			Volume         int     `json:"volume"`
		} `json:"analogs"`
	} `json:"data"`
}

type BasketAddResponse struct {
	Success bool `json:"success"`
	Data    struct {
		BasketID int `json:"basketId"`
	} `json:"data"`
}

type PositionInfoResponse struct {
	Success bool       `json:"success"`
	Data    []Position `json:"data"` //всегда 1 позиция
}

type Position struct {
	OrderID         int                   `json:"orderId"`
	OrderNumber     string                `json:"orderNumber"`
	OrderPositionID int                   `json:"orderPositionId"`
	PriceLogo       string                `json:"priceLogo"`
	BrandID         int                   `json:"brandId"`
	Brand           string                `json:"brand"`
	Code            string                `json:"code"`        //приходит неочищенный
	ReplaceCode     string                `json:"replaceCode"` //если они меняют номер при заказе
	DescriptionRus  string                `json:"descriptionRus"`
	DescriptionUa   string                `json:"descriptionUa"`
	Price           utilits.CustomFloat64 `json:"price"`
	Currency        string                `json:"currency"`
	Reference       string                `json:"reference"` //приходит в верхнем регистре
	Comment         string                `json:"comment"`
	AdminComment    string                `json:"adminComment"`
	States          []StatePosition       `json:"states"` //может быть несколько с отказными
}

type StatePosition struct {
	Quantity          int                `json:"quantity"`
	StatusID          int                `json:"statusId"`
	Status            string             `json:"status"`
	StatusChangedDate utilits.CustomTime `json:"statusChangedDate"` //"2023-10-24 17:28:02.716754"
}

// К крупногабаритным (объемным) деталям относятся:
// Детали кузова – капоты, крылья, двери, крыши, обшивка потолка крыши, крышка багажника, передняя панель радиатора, рейлинги верхней корзины, бампера, подкрыльники, корпуса панели приборов (торпеды), спойлер, накладки и спойлера бампера.
