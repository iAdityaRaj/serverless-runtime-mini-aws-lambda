package registry

//function model
import "time"

type Function struct {
	Name string `json:"name"`

	Image string `json:"image"`

	Runtime string `json:"runtime"`

	CreatedAt time.Time `json:"created_at"`

	InvocationCount int `json:"invocation_count"`
}
