package v1

type ManufacturerInput struct {
	ID      string   `json:"id"`
	Details *Details `json:"details"`
}

type Details struct {
	Name       string `json:"name"`
	NeedUpdate bool   `json:"needUpdate"`
}
