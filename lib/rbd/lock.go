package rbd

type Lock struct {
	ID      string `json:"id"`
	Locker  string `json:"locker"`
	Address string `json:"address"`
}
