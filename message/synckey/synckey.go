package message

type SyncKey struct {
	Count int      `json:"Count"`
	List  []KeyVal `json:"List"`
}
