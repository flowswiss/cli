package flow

type KeyPair struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}
