package client

import (
	"context"
	"errors"
	"github.com/drizzleaio/http"
	"github.com/drizzleaio/http/cookiejar"
	tls "gitlab.com/yawning/utls.git"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

type CustomClient struct {
	*http.Client
	transport *http.Transport

	// this practice is considered "bad"
	// but since we are spawning a new client for every task
	// I think it'll be ok
	context context.Context

	mu        sync.RWMutex
	preHooks  []func(req *http.Request) error
	postHooks []func(resp *http.Response) error
}

func New(context context.Context) *CustomClient {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true, ServerName: ""},
		Fingerprint:           "chrome91", // The fingerprint you wish to use,
	}

	return &CustomClient{
		Client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
			Jar:       jar,
		},
		context:   context,
		preHooks:  make([]func(req *http.Request) error, 0),
		postHooks: make([]func(resp *http.Response) error, 0),
		transport: transport,
	}
}

// Do this function overrides the http client Do function
// this allows us to have pre hooks and post hooks as well
// allows us to use a context to stop requests.
func (c *CustomClient) Do(r *http.Request) (*http.Response, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, hook := range c.preHooks {
		if err := hook(r); err != nil {
			return nil, err
		}
	}

	if c.context == nil {
		return c.Client.Do(r)
	}

	request, err := http.NewRequestWithContext(c.context, r.Method, r.URL.String(), r.Body)
	if err != nil {
		return nil, err
	}

	if r.ContentLength > 0 {
		request.ContentLength = r.ContentLength
	}

	request.Header = r.Header
	request.HeaderOrder = r.HeaderOrder

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	for _, hook := range c.postHooks {
		if err = hook(response); err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (c *CustomClient) SetProxy(proxyStr string) {
	parsedProxy, err := url.Parse(proxyStr)
	if err == nil {
		c.transport.Proxy = http.ProxyURL(parsedProxy)

		c.Transport = c.transport
	}
}

// UpdateServerName updates the transport server name
func (c *CustomClient) UpdateServerName(serverName string) {
	c.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true, ServerName: serverName}
	c.Transport = c.transport
}

// SetCookies sets an array of cookie to a specified domainName
func (c *CustomClient) SetCookies(domainName string, cookies []*http.Cookie) error {
	if c.Jar == nil {
		return errors.New("no cookie jar found")
	}

	urlParsed, err := url.Parse(domainName)
	if err != nil {
		return err
	}

	c.Client.Jar.SetCookies(urlParsed, cookies)
	return nil
}

// RemoveCookie removes a specified cookie to a specified domainName
func (c *CustomClient) RemoveCookie(domainName, cookie string) error {
	if c.Jar == nil {
		return errors.New("no cookie jar found")
	}

	urlParsed, err := url.Parse(domainName)
	if err != nil {
		return err
	}

	newCookie := &http.Cookie{
		Name:  cookie,
		Value: "",
	}

	c.Jar.SetCookies(urlParsed, []*http.Cookie{newCookie})

	return nil
}

// SetCookie sets 1 cookie to a specified domainName
func (c *CustomClient) SetCookie(domainName string, cookie *http.Cookie) error {
	if c.Jar == nil {
		return errors.New("no cookie jar found")
	}

	urlParsed, err := url.Parse(domainName)
	if err != nil {
		return err
	}

	c.Client.Jar.SetCookies(urlParsed, []*http.Cookie{cookie})
	return nil
}

// GetCookie gets a cookie with-out and specified domain.
func (c *CustomClient) GetCookie(name string, equals bool) *http.Cookie {
	if c.Jar == nil {
		return nil
	}

	cookies := c.Jar.GetCookies()
	if cookies == nil {
		return nil
	}

	for _, cookie := range cookies {
		if equals {
			if cookie.Name == name {
				return cookie
			}
		} else if strings.Contains(cookie.Name, name) {
			return cookie
		}
	}

	return nil
}

// GetCookieFromDomain gets a specific cookie from a specified domain
func (c *CustomClient) GetCookieFromDomain(domain, name string, equals bool) *http.Cookie {
	if c.Jar == nil {
		return nil
	}

	urlParsed, _ := url.Parse(domain)

	cookies := c.Client.Jar.Cookies(urlParsed)
	for _, cookie := range cookies {
		if equals {
			if cookie.Name == name {
				//fmt.Println("found cookie")
				return cookie
			}
		} else {
			if strings.Contains(cookie.Name, name) {
				return cookie
			}
		}
	}

	return nil
}

// DoRedirect sets if the client will do redirects.
func (c *CustomClient) DoRedirect(doRedirect bool) {
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if doRedirect {
			if len(via) > 10 {
				return errors.New("max redirects hit")
			}

			return nil
		}

		return http.ErrUseLastResponse
	}
}

// AddPreHook adds a pre-hook that is called in the Do function
// done before the request
func (c *CustomClient) AddPreHook(hook func(req *http.Request) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.preHooks = append(c.preHooks, hook)
}

// AddPostHook adds a post-hook that is called in the Do function
// done after the response.
func (c *CustomClient) AddPostHook(hook func(resp *http.Response) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.postHooks = append(c.postHooks, hook)
}
