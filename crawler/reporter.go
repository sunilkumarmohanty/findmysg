package crawler

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
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
	r := &reporter{
		results: results,
		writer:  writer,
		done:    make(chan bool, 1),
	}
	// Print Header

	r.print(&result{
		ID:            aws.String("ID"),
		Type:          aws.String("Type"),
		SecurityGroup: aws.String("SecurityGroup"),
	})
	return r

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
