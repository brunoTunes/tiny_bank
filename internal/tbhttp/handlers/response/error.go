package response

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}
