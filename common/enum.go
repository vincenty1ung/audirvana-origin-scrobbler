package common

var (
	AppSystemEvents    = "System Events"
	AppAudirvanaOrigin = "Audirvana Origin"
)

type PlayerState string

const (
	PlayerStateDefault = ""
	PlayerStateStopped = "Stopped"
	PlayerStatePlaying = "Playing"
	PlayerStatePaused  = "Paused"
)
