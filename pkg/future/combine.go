package future

import "sync"

func All(futures ...Future) Future {
	complete := New()
	var wg sync.WaitGroup
	wg.Add(len(futures))
	success := true

	for _, f := range futures {
		go func(f Future) {
			select {
			case r := <-f:
				if !r.IsOK() {
					complete.Fail(r.Err)
					success = false
				}
				wg.Done()
			}
		}(f)
	}

	go func() {
		wg.Wait()
		if success {
			complete.Succeed()
		}
	}()

	return complete
}
