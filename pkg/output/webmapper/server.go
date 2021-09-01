package webmapper

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"net/http"
)

func Handler() http.Handler {
	// TODO: Implement
	return nil
}

func TrackerCallback(details common.TrackedDetails) future.Future {
	// TODO: Implement
	f := future.New()
	f.Fail(fmt.Errorf("not implemented"))

	return f
}
