package internal

import tele "gopkg.in/telebot.v3"

func EndpointFormat(s string) string {
	switch s[0] {
	case '\a':
		return endpointName(s)
	case '\f':
		return "CallbackUnique(" + s[1:] + ")"
	}
	return s
}

func endpointName(endpoint string) string {
	switch endpoint {
	case tele.OnText:
		return "OnText"
	case tele.OnEdited:
		return "OnEdited"
	case tele.OnPhoto:
		return "OnPhoto"
	case tele.OnAudio:
		return "OnAudio"
	case tele.OnAnimation:
		return "OnAnimation"
	case tele.OnDocument:
		return "OnDocument"
	case tele.OnSticker:
		return "OnSticker"
	case tele.OnVideo:
		return "OnVideo"
	case tele.OnVoice:
		return "OnVoice"
	case tele.OnVideoNote:
		return "OnVideoNote"
	case tele.OnContact:
		return "OnContact"
	case tele.OnLocation:
		return "OnLocation"
	case tele.OnVenue:
		return "OnVenue"
	case tele.OnDice:
		return "OnDice"
	case tele.OnInvoice:
		return "OnInvoice"
	case tele.OnPayment:
		return "OnPayment"
	case tele.OnGame:
		return "OnGame"
	case tele.OnPoll:
		return "OnPoll"
	case tele.OnPollAnswer:
		return "OnPollAnswer"
	case tele.OnPinned:
		return "OnPinned"
	case tele.OnChannelPost:
		return "OnChannelPost"
	case tele.OnEditedChannelPost:
		return "OnEditedChannelPost"
	case tele.OnTopicCreated:
		return "OnTopicCreated"
	case tele.OnTopicReopened:
		return "OnTopicReopened"
	case tele.OnTopicClosed:
		return "OnTopicClosed"
	case tele.OnTopicEdited:
		return "OnTopicEdited"
	case tele.OnGeneralTopicHidden:
		return "OnGeneralTopicHidden"
	case tele.OnGeneralTopicUnhidden:
		return "OnGeneralTopicUnhidden"
	case tele.OnWriteAccessAllowed:
		return "OnWriteAccessAllowed"
	case tele.OnAddedToGroup:
		return "OnAddedToGroup"
	case tele.OnUserJoined:
		return "OnUserJoined"
	case tele.OnUserLeft:
		return "OnUserLeft"
	case tele.OnUserShared:
		return "OnUserShared"
	case tele.OnChatShared:
		return "OnChatShared"
	case tele.OnNewGroupTitle:
		return "OnNewGroupTitle"
	case tele.OnNewGroupPhoto:
		return "OnNewGroupPhoto"
	case tele.OnGroupPhotoDeleted:
		return "OnGroupPhotoDeleted"
	case tele.OnGroupCreated:
		return "OnGroupCreated"
	case tele.OnSuperGroupCreated:
		return "OnSuperGroupCreated"
	case tele.OnChannelCreated:
		return "OnChannelCreated"
	case tele.OnMigration:
		return "OnMigration"
	case tele.OnMedia:
		return "OnMedia"
	case tele.OnCallback:
		return "OnCallback"
	case tele.OnQuery:
		return "OnQuery"
	case tele.OnInlineResult:
		return "OnInlineResult"
	case tele.OnShipping:
		return "OnShipping"
	case tele.OnCheckout:
		return "OnCheckout"
	case tele.OnMyChatMember:
		return "OnMyChatMember"
	case tele.OnChatMember:
		return "OnChatMember"
	case tele.OnChatJoinRequest:
		return "OnChatJoinRequest"
	case tele.OnProximityAlert:
		return "OnProximityAlert"
	case tele.OnAutoDeleteTimer:
		return "OnAutoDeleteTimer"
	case tele.OnWebApp:
		return "OnWebApp"
	case tele.OnVideoChatStarted:
		return "OnVideoChatStarted"
	case tele.OnVideoChatEnded:
		return "OnVideoChatEnded"
	case tele.OnVideoChatParticipants:
		return "OnVideoChatParticipants"
	case tele.OnVideoChatScheduled:
		return "OnVideoChatScheduled"
	}
	return endpoint
}
