package l10n

import (
	"fmt"

	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
)

func TrackerCallback(details common.TrackedDetails) future.Future {
	return future.Exec(func() error {
		fmt.Println(Phrase(details, false))
		return nil
	})
}
