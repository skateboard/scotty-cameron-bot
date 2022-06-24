package module

import (
	"Scotty/internal/task"
	"encoding/json"
	"github.com/drizzleaio/http"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type ProductResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Quantity int    `json:"Quantity"`
	Data     struct {
		Style                 string      `json:"style"`
		Price                 string      `json:"price"`
		Sku                   string      `json:"sku"`
		ProductId             int         `json:"productId"`
		AddOnMessage          interface{} `json:"addOnMessage"`
		IsOutOfStock          interface{} `json:"isOutOfStock"`
		QtyExceedingInventory interface{} `json:"qtyExceedingInventory"`
	} `json:"data"`
}

type InventoryProducts struct {
	Products []ProductInventorySku
}

type ProductInventorySku struct {
	Sku  int    `json:"sku"`
	Type string `json:"type"`
}

type InventoryResponse struct {
	Status bool `json:"status"`
	Data   []struct {
		SKU                        string        `json:"SKU"`
		PublishProductId           int           `json:"PublishProductId"`
		IsQuantityRestricted       bool          `json:"IsQuantityRestricted"`
		IsItemAvailableForPurchase bool          `json:"IsItemAvailableForPurchase"`
		Quantity                   int           `json:"Quantity"`
		ReOrderLevel               int           `json:"ReOrderLevel"`
		InventoryMessage           string        `json:"InventoryMessage"`
		ShowAddToCart              bool          `json:"ShowAddToCart"`
		CartQuantity               int           `json:"CartQuantity"`
		Inventory                  []interface{} `json:"Inventory"`
	} `json:"data"`
}

func (t *ScottyTask) getProduct() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/store/product/getproductprice/?productSKU="+t.options.Sku+"&parentProductSKU="+t.options.Sku+"&quantity=1&addOnIds=&parentProductId=6298&_=1655877004571", nil)
	if err != nil {
		log.Println("Get Product: Failed to create request.")
		return FailedGetProduct
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "*/*"},
		{"accept-language": "en-US,en;q=0.5"},
		{"accept-encoding": "zip, deflate, br"},
		{"referer": t.options.Url},
		{"content-type": "application/x-www-form-urlencoded; charset=UTF-8"},
		{"x-requested-with": "XMLHttpRequest"},
		{"content-length": ""},
		{"origin": "https://www.scottycameron.com"},
		{"dnt": "1"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "none"},
		{"te": "trailers"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		log.Println("Get Product: Failed to get response.")
		return FailedGetProduct
	}

	if response.StatusCode != 200 {
		log.Println("Get Product: Failed to get response.")
		return FailedGetProduct
	}

	var productResponse ProductResponse
	err = response.GetBodyJSON(&productResponse)
	if err != nil {
		log.Println("Get Product: Failed to get response.")
		return FailedGetProduct
	}

	if !productResponse.Success {
		log.Println("Get Product: Failed to get response.")
		return FailedGetProduct
	}

	if productResponse.Message != "In-stock" {
		log.Println("Get Product: Failed to get response.")
		return FailedGetProduct
	}
	log.Println("Product in stock.")

	t.parameters.Price = productResponse.Data.Price

	return AddToCart
}

func (t *ScottyTask) getProductInventory() task.State {
	skuNum, _ := strconv.Atoi(t.options.Sku)

	products := InventoryProducts{Products: []ProductInventorySku{
		{Sku: skuNum, Type: "SimpleProduct"},
	}}
	jsonBytes, _ := json.Marshal(products)

	form := url.Values{}
	form.Add("products", string(jsonBytes))
	form.Add("categoryId", t.options.CategoryID)

	request, err := http.NewRequest("POST", "https://www.scottycameron.com/store/scottyproduct/getproductinventory/", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Get Product Inventory: Failed to create request.")
		return FailedGetInventory
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "*/*"},
		{"accept-language": "en-US,en;q=0.5"},
		{"accept-encoding": "zip, deflate, br"},
		{"referer": t.options.Url},
		{"x-newrelic-id": t.parameters.RelicID},
		{"content-type": "application/x-www-form-urlencoded; charset=UTF-8"},
		{"x-requested-with": "XMLHttpRequest"},
		{"content-length": ""},
		{"origin": "https://www.scottycameron.com"},
		{"dnt": "1"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "none"},
		{"te": "trailers"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		log.Println("Get Product Inventory: Failed to get response.")
		return FailedGetInventory
	}

	if response.StatusCode != 200 {
		log.Println("Get Product Inventory: Failed to get response.")
		return FailedGetInventory
	}

	var inventoryResponse InventoryResponse
	err = response.GetBodyJSON(&inventoryResponse)
	if err != nil {
		log.Println("Get Product Inventory: Failed to get response.")
		return FailedGetInventory
	}

	if !inventoryResponse.Status {
		log.Println("Get Product Inventory: Failed to get response.")
		return FailedGetInventory
	}

	if len(inventoryResponse.Data) == 0 {
		log.Println("Get Product Inventory: Failed to get response.")
		return FailedGetInventory
	}

	if inventoryResponse.Data[0].IsQuantityRestricted {
		log.Println("Product is Quantity Restricted.")
	}

	t.parameters.ProductID = inventoryResponse.Data[0].PublishProductId

	return AddToCart
}
