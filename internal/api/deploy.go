package api

type DeployRequest struct {
	Name string `json:"name"`

	Code string `json:"code"`
}
