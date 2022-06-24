package module

import (
	"Scotty/internal/task"
	"github.com/drizzleaio/http"
	"net/url"
	"strings"
)

type CartResponse struct {
	Status         bool   `json:"status"`
	TimerDuration  int    `json:"timerDuration"`
	IsTimerEnabled bool   `json:"isTimerEnabled"`
	CartCount      int    `json:"cartCount"`
	SuccessMessage string `json:"successMessage"`
}

func (t *ScottyTask) addToCart() task.State {
	form := url.Values{}
	form.Add("ProductId", "5934")
	form.Add("SKU", "33700")
	form.Add("ProductType", "SimpleProduct")
	form.Add("Quantity", "1")
	form.Add("ParentProductId", "5934")
	form.Add("ConfigurableProductSKUs", "")
	form.Add("AddOnProductSKUs", "")
	form.Add("PersonalisedCodes", "")
	form.Add("PersonalisedValues", "")
	form.Add("IsRedirectToCart", "False")
	form.Add("__RequestVerificationToken", "aioKIqJgWH488lBciTMVeIQzFMdun7Falpa2gJ0lEpAmoqbfD23HGRScdGxPuq0CmVZ9WYQhY0ew06ZCXardUMwLRyiGlNpDfMyJXdgIML25gBaYtk8Mu48iOKw8Yeht0")
	form.Add("IsProductEdit", "undefined")
	form.Add("X-Requested-With", "XMLHttpRequest")

	request, err := http.NewRequest("POST", "https://www.scottycameron.com/store/scottyproduct/addtocartproduct", strings.NewReader(form.Encode()))
	if err != nil {
		return FailedAddToCart
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "*/*"},
		{"accept-language": "en-US,en;q=0.5"},
		{"accept-encoding": "zip, deflate, br"},
		{"referer": t.options.Url},
		{"x-newrelic-id": t.parameters.RelicID},
		{"x-requested-with": "XMLHttpRequest"},
		{"content-type": "application/x-www-form-urlencoded; charset=UTF-8"},
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
		return FailedAddToCart
	}

	if response.StatusCode != 200 {
		return FailedAddToCart
	}
	//6Kx4gxYsywi78HS
	var cartResponse CartResponse
	err = response.GetBodyJSON(&cartResponse)
	if err != nil {
		return FailedAddToCart
	}

	if strings.Contains(cartResponse.SuccessMessage, "Product is added to the cart") {
		return AddToCart
	}

	return 0
}
