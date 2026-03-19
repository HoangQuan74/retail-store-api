package nats

const (
	StreamRetailStore = "RETAIL_STORE"

	SubjectOrderCreated  = "orders.created"
	SubjectOrderUpdated  = "orders.updated"
	SubjectProductViewed = "products.viewed"
	SubjectProductCreated = "products.created"
	SubjectProductUpdated = "products.updated"
	SubjectProductDeleted = "products.deleted"
)

var StreamSubjects = []string{
	"orders.>",
	"products.>",
}
