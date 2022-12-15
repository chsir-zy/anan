package contract

import "net/http"

const KernelKey = "anan:kernel"

type Kernel interface {
	HttpEngine() http.Handler
}
