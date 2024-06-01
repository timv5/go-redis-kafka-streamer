package request

type RetrieveRequest struct {
	UUID string `json:"uuid" binding:"required"`
}
