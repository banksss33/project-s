package game

import (
	"encoding/json"
	"project-s/internal/types"
)

func LobbyActionDispatcher() map[string]func(*GameRoom, *LobbyState, types.Action) {
	lobbyAction := make(map[string]func(*GameRoom, *LobbyState, types.Action))

	lobbyAction["PLAYER_JOIN"] = joinGame
	lobbyAction["SPECTATOR_JOIN"] = joinSpectator
	lobbyAction["PLAYER_LEFT"] = leaveGame
	lobbyAction["EDIT_GAME_SETTING"] = editGameSetting

	return lobbyAction
}

func joinGame(room *GameRoom, lobby *LobbyState, payload types.Action) {
	player := room.PlayerList[payload.CallerID]
	if player == nil {
		return
	}
	lobby.PlayerJoin(player)

	broadcastLobbyState(room, lobby)
	sendGameSettings(player, lobby)
}

func joinSpectator(room *GameRoom, lobby *LobbyState, payload types.Action) {
	player, exists := room.GetPlayerByID(payload.CallerID)
	if !exists {
		return
	}

	if !lobby.PlayerList[player] { //case when player are already spectator
		return
	}

	lobby.SpectatorJoin(player)

	broadcastLobbyState(room, lobby)
}

func leaveGame(room *GameRoom, lobby *LobbyState, payload types.Action) {
	player, exists := room.PlayerList[payload.CallerID]
	if !exists {
		return
	}

	lobby.PlayerLeft(player)
	delete(room.PlayerList, payload.CallerID)

	if len(room.PlayerList) == 0 {
		room.GameClose <- true
		return
	}

	broadcastLobbyState(room, lobby)
}

func editGameSetting(room *GameRoom, lobby *LobbyState, payload types.Action) {
	player := room.PlayerList[payload.CallerID]
	if player == nil || lobby.Host != player {
		return
	}

	var payloadGameSetting types.GameSettingPayload
	if err := json.Unmarshal(payload.Payload, &payloadGameSetting); err != nil {
		return
	}

	setting := types.GameSetting{
		Round:     payloadGameSetting.Round,
		Spies:     payloadGameSetting.Spies,
		Timer:     payloadGameSetting.Timer,
		Locations: payloadGameSetting.Locations,
	}
	lobby.EditSetting(setting)

	response := types.ServerResponse{
		ResponseName: "GAME_SETTING",
		Payload:      payload.Payload,
	}

	room.Broadcast <- response
}

// #region helper func
// Helper to broadcast UPDATE_LOBBY_PLAYER_LIST
func broadcastLobbyState(room *GameRoom, lobby *LobbyState) {
	responsePayload := types.UpdateLobbyPlayerListResponse{
		Host:       "",
		Players:    make([]string, 0),
		Spectators: make([]string, 0),
	}
	if lobby.Host != nil {
		responsePayload.Host = lobby.Host.UserID
	}

	//if isPlayer = true it's player if not that player is spectator
	for player, isPlayer := range lobby.PlayerList {
		if isPlayer {
			responsePayload.Players = append(responsePayload.Players, player.UserID)
			continue
		}

		responsePayload.Spectators = append(responsePayload.Spectators, player.UserID)

	}

	jsonResponsePayload, err := json.Marshal(responsePayload)
	if err != nil {
		return
	}

	room.Broadcast <- types.ServerResponse{
		ResponseName: "UPDATE_LOBBY_PLAYER_LIST",
		Payload:      jsonResponsePayload,
	}
}

// Helper to send current settings to specific player
func sendGameSettings(player *Player, lobby *LobbyState) {
	settingJSON := types.GameSettingPayload{
		Round:     lobby.setting.Round,
		Spies:     lobby.setting.Spies,
		Timer:     lobby.setting.Timer,
		Locations: lobby.setting.Locations,
	}

	gameSettingJson, err := json.Marshal(settingJSON)
	if err != nil {
		return
	}

	player.SendToPlayer <- types.ServerResponse{
		ResponseName: "GAME_SETTING",
		Payload:      gameSettingJson,
	}
}
