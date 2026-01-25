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
	//Mock Data supposed to be from data
	var mockLocations = map[string][]string{
		"Airplane":      {"Pilot", "Co-pilot", "Flight Attendant", "Passenger", "Air Marshal", "Mechanic"},
		"Hospital":      {"Surgeon", "Nurse", "Patient", "Therapist", "Security Guard", "Doctor"},
		"Bank":          {"Bank Manager", "Teller", "Robber", "Customer", "Armored Car Driver", "Consultant"},
		"Pirate Ship":   {"Captain", "Cabin Boy", "Pirate", "Gunner", "Bound Prisoner", "Cook"},
		"Space Station": {"Astronaut", "Scientist", "Commander", "Engineer", "Alien", "Space Tourist"},
	}

	lobby.LobbyStateInit(host, mockLocations)
}

func joinGame(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	lobby.PlayerJoin(player)

	responsePayload := types.UpdateLobbyPlayerListResponse{
		Players:    make([]string, 0),
		Spectators: make([]string, 0),
	}
	responsePayload.Host = lobby.Host.UserID
	for player, isPlayer := range lobby.PlayerList {
		if isPlayer {
			responsePayload.Players = append(responsePayload.Players, player.UserID)
			continue
		}

		responsePayload.Spectators = append(responsePayload.Spectators, player.UserID)
	}

	jsonResponsePayload, _ := json.Marshal(responsePayload)

	response := types.ServerResponse{
		ResponseName: "UPDATE_LOBBY_PLAYER_LIST",
		Payload:      jsonResponsePayload,
	}
	//Broadcast
	room.Broadcast <- response

	gameSettingJson, _ := json.Marshal(lobby.setting)

	response.ResponseName = "GAME_SETTING"
	response.Payload = gameSettingJson

	player.SendToPlayer <- response

}

func joinSpectator(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	lobby.SpectatorJoin(player)

	responsePayload := types.UpdateLobbyPlayerListResponse{
		Players:    make([]string, 0),
		Spectators: make([]string, 0),
	}
	responsePayload.Host = lobby.Host.UserID
	for player, isPlayer := range lobby.PlayerList {
		if isPlayer {
			responsePayload.Players = append(responsePayload.Players, player.UserID)
			continue
		}

		responsePayload.Spectators = append(responsePayload.Spectators, player.UserID)
	}

	jsonResponsePayload, err := json.Marshal(responsePayload)
	if err != nil {
		panic("ERROR JOIN GAME")
	}

	response := types.ServerResponse{
		ResponseName: "UPDATE_LOBBY_PLAYER_LIST",
		Payload:      jsonResponsePayload,
	}
	//Broadcast
	room.Broadcast <- response
}

func leaveGame(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	lobby.PlayerLeft(player)
	delete(room.PlayerList, player.UserID)

	responsePayload := types.UpdateLobbyPlayerListResponse{
		Players:    make([]string, 0),
		Spectators: make([]string, 0),
	}
	responsePayload.Host = lobby.Host.UserID
	for player, isPlayer := range lobby.PlayerList {
		if isPlayer {
			responsePayload.Players = append(responsePayload.Players, player.UserID)
			continue
		}

		responsePayload.Spectators = append(responsePayload.Spectators, player.UserID)
	}

	jsonResponsePayload, err := json.Marshal(responsePayload)
	if err != nil {
		panic("ERROR JOIN GAME")
	}

	response := types.ServerResponse{
		ResponseName: "UPDATE_LOBBY_PLAYER_LIST",
		Payload:      jsonResponsePayload,
	}
	//Broadcast
	room.Broadcast <- response
}

func editGameSetting(room *GameRoom, lobby *LobbyState, payload types.PlayerAction) {
	player := room.PlayerList[payload.UserID]
	if lobby.Host != player {
		return
	}

	var gameSetting types.GameSetting
	json.Unmarshal(payload.Payload, &gameSetting)
	lobby.EditSetting(gameSetting)

	json, _ := json.Marshal(lobby.setting)
	response := types.ServerResponse{
		ResponseName: "GAME_SETTING",
		Payload:      json,
	}

	room.Broadcast <- response
}
