package crawler

import (
	"fmt"
	"io"
)

type result struct {
	ID            *string
	Type          *string
	SecurityGroup *string
}

type reporter struct {
	writer  io.Writer
	results chan *result
	done    chan bool
}

func newReporter(results chan *result, writer io.Writer) *reporter {
	return &reporter{
		results: results,
		writer:  writer,
		done:    make(chan bool, 1),
	}
}

func (r *reporter) run() {
	for res := range r.results {
		r.print(res)
	}
	//result channel is closed and the reporter reports its work is done
	r.done <- true
}

func (r *reporter) print(res *result) {
	fmt.Fprintf(r.writer, "%v,%v,%v\n", *res.Type, *res.ID, *res.SecurityGroup)
}
