package domain

type EventType string

var (
	POSTBACK_EVENT_TYPE      EventType = "postback"
	MESSAGE_EVENT_TYPE       EventType = "message"
	UNFOLLOW_EVENT_TYPE      EventType = "unfollow"
	LEAVE_EVENT_TYPE         EventType = "leave"
	MEMBER_LEFT_EVENT_TYPE   EventType = "memberLeft"
	MEMBER_JOINED_EVENT_TYPE EventType = "memberJoined"
	JOIN_EVENT_TYPE          EventType = "join"
	FOLLOW_EVENT_TYPE        EventType = "follow"
	THINGS_EVENT_TYPE        EventType = "things"
	BEACON_EVENT_TYPE        EventType = "beacon"
	ACCOUNT_LINK_EVENT_TYPE  EventType = "accountLink"
)

type JustifyContent string

var (
	SPACE_EVENLY_JUSTIFY_CONTENT JustifyContent = "space-evenly"
	FLEX_END_JUSTIFY_CONTENT     JustifyContent = "flex-end"
)

type AlignItems string

var (
	FLEX_START_ALIGN_ITEMS AlignItems = "flex-start"
	CENTER_ALIGN_ITEMS     AlignItems = "center"
)
