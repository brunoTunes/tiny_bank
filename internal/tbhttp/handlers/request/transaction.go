package request

type Transaction struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Amount      int    `json:"amount"`
}

type Withdraw struct {
	Amount int `json:"amount"`
}

type Deposit struct {
	Amount int `json:"amount"`
}
