package queue

import (
	"fmt"
	"os"
	"time"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/invoice"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}

type Billing struct{}

func (b *Billing) Run(qt QueueTask) error {
	id, ok := qt.Data.(string)
	if !ok {
		return fmt.Errorf("the data should be a stripe customer ID")
	}

	//我们将执行延迟2个小时，
	//以便在创建发票之间进行添加/删除操作，
	//因为我们正在执行 go routine，因此我们可以使用 time.Sleep。
	time.Sleep(2 * time.Hour)

	p := &stripe.InvoiceParams{Customer: &id}
	_, err := invoice.New(p)
	return err
}
