package game

import (
	"encoding/json"
	"project-s/internal/types"
)

func LobbyActionDispatcher() map[string]func(*GameRoom, *LobbyState, types.PlayerAction) {
	lobbyAction := make(map[string]func(*GameRoom, *LobbyState, types.PlayerAction))

	lobbyAction["GAME_CREATED"] = lobbyCreated
	lobbyAction["PLAYER_JOIN"] = joinGame
	lobbyAction["SPECTATOR_JOIN"] = joinSpectator
	lobbyAction["PLAYER_LEFT"] = leaveGame
	lobbyAction["EDIT_GAME_SETTING"] = editGameSetting

	return lobbyAction
}

func lobbyCreated(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	host := room.PlayerList[payload.UserID]
	lobby.LobbyStateInit(host)
}

func joinGame(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	if player == nil {
		return
	}
	lobby.PlayerJoin(player)

	broadcastLobbyState(room, lobby)
	sendGameSettings(player, lobby)
}

func joinSpectator(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	if player == nil {
		return
	}
	lobby.SpectatorJoin(player)

	broadcastLobbyState(room, lobby)
}

func leaveGame(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player, exists := room.PlayerList[payload.UserID]
	if !exists {
		return
	}

	lobby.PlayerLeft(player)
	delete(room.PlayerList, player.UserID)

	if len(room.PlayerList) == 0 {
		room.GameClose <- true
		return
	}

	broadcastLobbyState(room, lobby)
}

func editGameSetting(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
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

// Helper to broadcast UPDATE_LOBBY_PLAYER_LIST
func broadcastLobbyState(room *GameRoom, lobby *LobbyState) {
	responsePayload := types.UpdateLobbyPlayerListResponse{
		Host:       lobby.Host.UserID,
		Players:    make([]string, 0),
		Spectators: make([]string, 0),
	}

	for player, isPlayer := range lobby.PlayerList {
		if isPlayer {
			responsePayload.Players = append(responsePayload.Players, player.UserID)
		} else {
			responsePayload.Spectators = append(responsePayload.Spectators, player.UserID)
		}
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
