package harvester

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"time"
)

type Harvester struct {
	Name string

	Url string

	isLoaded bool
	Queue    *Queue

	ctx context.Context
}

type Solve struct {
	SiteKey string
	Channel chan SolveResult
}

type SolveResult struct {
	Error error
	Token string
}

func New(name, url string) *Harvester {
	return &Harvester{
		Name:  name,
		Url:   url,
		Queue: NewQueue(),
	}
}

func (h *Harvester) Initialize() error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:101.0) Gecko/20100101 Firefox/101.0"),
		chromedp.WindowSize(200, 700),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	if err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(fmt.Sprintf(`window.location = "%v"`, h.Url), nil),
		chromedp.EvaluateAsDevTools(headHtml, nil),
		chromedp.EvaluateAsDevTools(scriptLoader, nil),
	); err != nil {
		return err
	}

	//err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(fmt.Sprintf(`window.location = "%v"`, h.Url), nil))
	//if err != nil {
	//	return err
	//}
	//
	//err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(headHtml, nil))
	//if err != nil {
	//	return err
	//}
	//
	//err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(scriptLoader, nil))
	//if err != nil {
	//	return err
	//}

	h.ctx = ctx

	go h.clearQueue()
	select {
	case <-ctx.Done():
	}

	return nil
}

func (h *Harvester) Harvest(siteKey string) (string, error) {

	resultChannel := make(chan SolveResult)
	h.Queue.Push(Solve{
		SiteKey: siteKey,
		Channel: resultChannel,
	})

	resultParsed := <-resultChannel

	if resultParsed.Error != nil {
		return "", resultParsed.Error
	}

	return resultParsed.Token, nil
}

func (h *Harvester) executeHarvest(solve Solve) SolveResult {
	var result string
	if err := chromedp.Run(h.ctx,
		chromedp.Evaluate(fmt.Sprintf(`document.harv.harvest("%s")`, solve.SiteKey), &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	); err != nil {
		return SolveResult{
			Error: err,
		}
	}

	return SolveResult{
		Token: result,
		Error: nil,
	}
}

func (h *Harvester) Destroy() error {
	err := chromedp.Cancel(h.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (h *Harvester) clearQueue() {
	for {
		if h.Queue.Len() != 0 {
			firstElement := h.Queue.Pop()
			if firstElement == nil {
				time.Sleep(250 * time.Millisecond)
				continue
			}

			parsedElement := firstElement.(Solve)
			if parsedElement.SiteKey == "" {
				time.Sleep(250 * time.Millisecond)
				continue
			}

			result := h.executeHarvest(parsedElement)
			parsedElement.Channel <- result
		} else {
			time.Sleep(250 * time.Millisecond)
			continue
		}
	}
}
