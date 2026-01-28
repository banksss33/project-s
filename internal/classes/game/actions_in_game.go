package game

import (
	"encoding/json"
	"project-s/internal/types"
)

func InGameActionDispatcher() map[string]func(*GameRoom, *InGameState, types.Action) {
	InGameAction := make(map[string]func(*GameRoom, *InGameState, types.Action))

	InGameAction["ACCUSE_PLAYER"] = accusePlayer
	InGameAction["VOTE_ACCUSATION"] = voteAccusation
	InGameAction["DISAGREE_ACCUSATION"] = disagreeAccusation
	InGameAction["GUESS_LOCATION"] = spyGuessLocation
	InGameAction["PAUSE_TIMER"] = pauseTimer
	InGameAction["RESUME_TIMER"] = resumeTimer
	InGameAction["PLAYER_JOIN"] = inGamePlayerJoin
	InGameAction["PLAYER_LEFT"] = inGamePlayerLeft

	return InGameAction
}

func inGamePlayerJoin(room *GameRoom, inGame *InGameState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists {
		return
	}

	inGame.OnPlayerJoin(player)

	broadcastGameState(room, inGame)
}

func inGamePlayerLeft(room *GameRoom, inGame *InGameState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists {
		return
	}

	shouldCloseRoom := inGame.OnPlayerLeave(player)
	if shouldCloseRoom {
		room.GameClose <- true
		return
	}

	broadcastGameState(room, inGame)
}

func accusePlayer(room *GameRoom, inGame *InGameState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if exists {
		return
	}

	var accusePayload types.AccusePayload
	err := json.Unmarshal(payload.Payload, &accusePayload)
	if err != nil {
		return
	}

	targetPlayer, exists := room.GetPlayerByID(accusePayload.TargetUserID)
	if !exists {
		return
	}

	inGame.Accuse(player, targetPlayer)
	broadcastGameState(room, inGame)
}

func voteAccusation(room *GameRoom, inGame *InGameState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists || !inGame.isVoting {
		return
	}

	inGame.Vote(player)
	broadcastGameState(room, inGame)
}

func disagreeAccusation(room *GameRoom, inGame *InGameState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists || !inGame.isVoting {
		return
	}

	inGame.Disagree(player)
	broadcastGameState(room, inGame)

	disagreePlayerPayload := struct {
		PlayerID string `json:"player_id"`
	}{
		PlayerID: player.UserID,
	}
	disagreeJsonPayload, err := json.Marshal(disagreePlayerPayload)
	if err != nil {
		return
	}

	room.Broadcast <- types.ServerResponse{
		ResponseName: "PLAYER_WHO_DISAGREE",
		Payload:      disagreeJsonPayload,
	}
}

func spyGuessLocation(room *GameRoom, inGame *InGameState, payload types.Action) {
	fromPlayer := room.PlayerList[payload.CallerID]
	if status := inGame.playerStatus[fromPlayer]; status.Roles != "SPY" {
		return
	}

	var guessPayload types.SpyGuessPayload
	err := json.Unmarshal(payload.Payload, &guessPayload)
	if err != nil {
		return
	}

	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists {
		return
	}
	inGame.SpyVoteLocation(player, guessPayload.Location)
	broadcastGameState(room, inGame)
}

func pauseTimer(room *GameRoom, inGame *InGameState, payload types.Action) {
	inGame.PauseTimer()
	broadcastGameState(room, inGame)
}

func resumeTimer(room *GameRoom, inGame *InGameState, payload types.Action) {
	inGame.ResumeTimer()
	broadcastGameState(room, inGame)
}

// #region helper func
func broadcastGameState(room *GameRoom, inGame *InGameState) {
	playerStats := make(map[string]types.InGamePlayerStats)
	for player, status := range inGame.playerStatus {
		playerStats[player.UserID] = types.InGamePlayerStats{
			Score:       status.Score,
			AlreadyVote: status.AlreadyVote,
		}
	}

	gameStatus := types.GameStatus{
		IsTimeRunning: inGame.timer.IsRunning,
		IsVoting:      inGame.isVoting,
		IsRoundEnd:    inGame.isRoundEnd,
		RoundLeft:     inGame.roundLeft,
		PlayerList:    playerStats,
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

func SendPlayerRolesAndLocation(inGame *InGameState) {
	for player, status := range inGame.playerStatus {
		roleResponse := types.PlayerRoleResponse{
			Role: status.Roles,
		}
		if status.Roles == "SPY" {
			roleResponse.Location = "SPY"
		} else {
			roleResponse.Location = inGame.location
		}

		jsonRolePayload, _ := json.Marshal(roleResponse)
		response := types.ServerResponse{
			ResponseName: "ROLE_AND_LOCATION",
			Payload:      jsonRolePayload,
		}
		player.SendToPlayer <- response
	}
}
