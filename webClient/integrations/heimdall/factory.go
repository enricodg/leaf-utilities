package leafHeimdall

import (
	leafWebClient "github.com/enricodg/leaf-utilities/webClient/webClient"
	leafCircuitBreaker "github.com/enricodg/leaf-utilities/webClient/webClient/circuitBreaker"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/hystrix"
	"time"
)

type webClientFactory struct{}

func (cf *webClientFactory) Create(opts ...leafWebClient.Option) leafWebClient.WebClient {
	wcOption := leafWebClient.WebClientOption{}
	for _, opt := range opts {
		opt(&wcOption)
	}

	if 0 != wcOption.GetRetryCount() {
		return &WebClient{
			Doer: *httpclient.NewClient(
				httpclient.WithHTTPTimeout(wcOption.GetTimeout()),
				httpclient.WithRetryCount(wcOption.GetRetryCount())),
			EnableCircuitBreaker: false,
		}
	}

	if nil != wcOption.GetCircuitBreaker() {
		client := leafCircuitBreaker.CircuitBreakerOption{}
		for _, opt := range wcOption.GetCircuitBreaker()() {
			opt(&client)
		}
		options := make([]hystrix.Option, 0)

		durationZeroValue, _ := time.ParseDuration("0s")
		if durationZeroValue != client.GetTimeout() {
			options = append(options, hystrix.WithHTTPTimeout(client.GetTimeout()))
		}
		if durationZeroValue != client.GetCircuitBreakerTimeout() {
			options = append(options, hystrix.WithHystrixTimeout(client.GetCircuitBreakerTimeout()))
		}
		if "" != client.GetCommandName() {
			options = append(options, hystrix.WithCommandName(client.GetCommandName()))
		}
		if 0 != client.GetMaxConcurrentRequests() {
			options = append(options, hystrix.WithMaxConcurrentRequests(client.GetMaxConcurrentRequests()))
		}
		if 0 != client.GetRequestVolumeThreshold() {
			options = append(options, hystrix.WithRequestVolumeThreshold(client.GetRequestVolumeThreshold()))
		}
		if 0 != client.GetSleepWindow() {
			options = append(options, hystrix.WithSleepWindow(client.GetSleepWindow()))
		}
		if 0 != client.GetErrorPercentThreshold() {
			options = append(options, hystrix.WithErrorPercentThreshold(client.GetErrorPercentThreshold()))
		}
		if nil != client.GetFallbackFunc() {
			options = append(options, hystrix.WithFallbackFunc(client.GetFallbackFunc()))
		}
		if nil != client.GetStatsDCollector() {
			options = append(options, hystrix.WithStatsDCollector(client.GetStatsDCollector().StatsdAddr, client.GetStatsDCollector().Prefix))
		}
		if 0 != client.GetRetryCount() {
			options = append(options, hystrix.WithRetryCount(client.GetRetryCount()))
		}
		if nil != client.GetRetryBackoffOption() {
			options = append(options, hystrix.WithRetrier(convertToHeimdallRetryBackoff(*client.GetRetryBackoffOption())))
		}

		return &WebClient{
			HystrixDoer:          *hystrix.NewClient(options...),
			EnableCircuitBreaker: true,
		}
	}

	return &WebClient{
		Doer:                 *httpclient.NewClient(httpclient.WithHTTPTimeout(wcOption.GetTimeout())),
		EnableCircuitBreaker: false,
	}
}

func NewClientFactory() leafWebClient.Factory {
	return &webClientFactory{}
}
