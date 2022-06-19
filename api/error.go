package api

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

func checkError(err error) error {
	// レート制限
	if rateLimit, has := twitter.RateLimitFromError(err); has {
		t := rateLimit.Reset.Time().Local().Format("15:04:05")
		return fmt.Errorf("Rate limit exceeded (Reset time: %s)", t)
	}

	//calloutError := &json.UnsupportedValueError{}
	//decodeError := &twitter.ResponseDecodeError{}
	//httpError := &twitter.HTTPError{}

	//switch {
	//case errors.As(err, &calloutError):
	//	//
	//	log.Println("callout")
	//	return calloutError.Str
	//case errors.As(err, &decodeError):
	//	//
	//	log.Println("decode")
	//	return decodeError.Name
	//case errors.As(err, &httpError):
	//	//
	//	return fmt.Sprintf("%d %s", httpError.StatusCode, httpError.Status)
	//case err != nil:
	//	return fmt.Sprintf("unknown: %s", err.Error())
	//}

	return err
}
