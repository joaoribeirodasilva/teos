package requests

type HistoriesCreate struct {
	AuthenticationKey string `json:"authenticationKey"`
	Table             string `json:"table"`
	OriginalID        uint   `json:"originalId"`
	Data              string `json:"data"`
	CreatedBy         uint   `json:"createdBy"`
	UpdatedBy         uint   `json:"updatedBy"`
	DeletedBy         *uint  `json:"deletedBy"`
}
