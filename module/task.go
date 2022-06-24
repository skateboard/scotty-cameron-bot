package module

import (
	"Scotty/internal/account"
	"Scotty/internal/harvester"
	"Scotty/internal/task"
	"log"
	"regexp"
)

type ScottyTask struct {
	*task.TaskBase

	options Options

	parameters Parameters

	account account.Account
}

type Options struct {
	Sku        string
	CategoryID string
	Url        string

	harvester *harvester.Harvester
}

type Parameters struct {
	AkamaiSensor string

	RelicID   string
	UserAgent string

	Price string

	ProductID int

	RequestVerificationToken string
}

const (
	Session task.State = iota
	FailedSession

	GetAkamai
	FailedGetAkamai

	Login
	FailedLogin
	FailedLoginParameters

	GetProduct
	FailedGetProduct

	GetInventory
	FailedGetInventory

	AddToCart
	FailedAddToCart

	OutOfStock
)

var (
	requestVerifyTokenRegex = regexp.MustCompile(`(?m)<input name="__RequestVerificationToken" type="hidden" value="(.*)" />`)
)

func (t *ScottyTask) Next(state task.State) (task.State, error) {
	switch state {
	case task.Initialize:
		t.Client.UpdateServerName("www.scottycameron.com")
		t.Client.SetProxy("http://127.0.0.1:8888")
		t.Client.DoRedirect(false)

		t.parameters.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"

		return Session, nil
	case Session:
		log.Println("Session")
		return t.getSession(), nil
	case FailedSession:
		log.Println("Failed to get session")
	case GetAkamai:
		log.Println("Get Akamai")
		return t.getAkamaiSensorEP(), nil
	case FailedGetAkamai:
		log.Println("Failed to get Akamai sensor EP")
	case Login:
		log.Println("Login")
		return t.getLoginParameters(), nil
	case FailedLogin:
		log.Println("Failed to login")
	case GetProduct:
		log.Println("Get Product")
		return t.getProduct(), nil
	case FailedGetProduct:
		log.Println("Failed to get product")
	case GetInventory:
		log.Println("Get Inventory")
		return t.getProductInventory(), nil
	case FailedGetInventory:
		log.Println("Failed to get inventory")
	case OutOfStock:
		log.Println("Out of stock")
	case AddToCart:
		log.Println("Add to cart")
		return t.addToCart(), nil
	case FailedAddToCart:
		log.Println("Failed to add to cart")
	}

	t.Stop()
	return 0, nil
}
