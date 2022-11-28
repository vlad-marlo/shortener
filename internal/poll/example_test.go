package poll_test

import "github.com/vlad-marlo/shortener/internal/poll"

func Example() {
	// ...
	// init poller
	poller := poll.New(nil)
	// always defer poller close
	defer poller.Close()
}
