package game

import (
	"encoding/json"
	"project-s/internal/types"
)

func InGameActionDispatcher() map[string]func(*GameRoom, *InGameState, types.PlayerAction) {
	InGameAction := make(map[string]func(*GameRoom, *InGameState, types.PlayerAction))

	InGameAction["ACCUSE_PLAYER"] = accusePlayer
	InGameAction["VOTE_ACCUSATION"] = voteAccusation
	InGameAction["DISAGREE_ACCUSATION"] = disagreeAccusation
	InGameAction["GUESS_LOCATION"] = spyGuessLocation
	InGameAction["PAUSE_TIMER"] = pauseTimer
	InGameAction["RESUME_TIMER"] = resumeTimer

	return InGameAction
}

func accusePlayer(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	player := room.GetPlayerByID(payload.UserID)
	if player == nil {
		return
	}

	var accusePayload types.AccusePayload
	err := json.Unmarshal(payload.Payload, &accusePayload)
	if err != nil {
		return
	}

	targetPlayer := room.GetPlayerByID(accusePayload.TargetUserID)
	if targetPlayer == nil {
		return
	}

	inGame.Accruse(player, targetPlayer)
	broadcastGameState(room, inGame)
}

// c
func voteAccusation(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	player := room.GetPlayerByID(payload.UserID)
	if player == nil || !inGame.isVoting {
		return
	}

	inGame.Vote(player)
	broadcastGameState(room, inGame)
}

func disagreeAccusation(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	player := room.GetPlayerByID(payload.UserID)
	if player == nil || !inGame.isVoting {
		return
	}

	inGame.Disagree(player)
	broadcastGameState(room, inGame)

	disagreePlayerPayload := struct {
		PlayerID string `json:"player_id"`
	}{}
	disagreeJsonPayload, err := json.Marshal(disagreePlayerPayload)
	if err != nil {
		return
	}

	room.Broadcast <- types.ServerResponse{
		ResponseName: "PLAYER_WHO_DISAGREE",
		Payload:      disagreeJsonPayload,
	}
}

// c
func spyGuessLocation(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	var guessPayload types.SpyGuessPayload
	err := json.Unmarshal(payload.Payload, &guessPayload)
	if err != nil {
		return
	}

	player := room.GetPlayerByID(payload.UserID)
	inGame.SpyVoteLocation(player, guessPayload.Location)
	broadcastGameState(room, inGame)
}

func pauseTimer(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	inGame.PauseTimer()
	broadcastGameState(room, inGame)
}

func resumeTimer(room *GameRoom, inGame *InGameState, payload types.PlayerAction) {
	inGame.ResumeTimer()
	broadcastGameState(room, inGame)
}

func broadcastGameState(room *GameRoom, inGame *InGameState) {
	gameStatus := types.GameStatus{
		IsTimeRunning: inGame.timer.IsRunning,
		IsVoting:      inGame.isVoting,
		IsRoundEnd:    inGame.isRoundEnd,
		RoundLeft:     inGame.roundLeft,
	}
	jsonPayload, err := json.Marshal(gameStatus)
	if err != nil {
		return
	}

	response := types.ServerResponse{
		ResponseName: "GAME_STATUS_UPDATE",
		Payload:      jsonPayload,
	}

	room.Broadcast <- response
}

func broadcastPlayerRolesAndLocation(inGame *InGameState) {
	for player, status := range inGame.playerStatus {
		roleResponse := types.PlayerRoleResponse{
			Role: status.Roles,
		}
		jsonRolePayload, _ := json.Marshal(roleResponse)
		response := types.ServerResponse{
			ResponseName: "ROLE_AND_LOCATION",
			Payload:      jsonRolePayload,
		}
		player.SendToPlayer <- response
	}
}
