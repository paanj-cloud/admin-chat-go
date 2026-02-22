package chatadmin

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/paanj-cloud/paanj-go/admin"
)

type AdminChat struct {
	admin         *admin.PaanjAdmin
	Conversations *AdminConversationsResource
	Users         *AdminUsersResource
	Messages      *AdminMessagesResource
}

func NewAdminChat(a *admin.PaanjAdmin) *AdminChat {
	adminChat := &AdminChat{
		admin: a,
	}

	adminChat.Conversations = NewAdminConversationsResource(a)
	adminChat.Users = NewAdminUsersResource(a)
	adminChat.Messages = NewAdminMessagesResource(a)

	return adminChat
}

func (c *AdminChat) Conversation(conversationID string) *AdminConversationContext {
	return c.Conversations.Conversation(conversationID)
}

func (c *AdminChat) User(blockerID string) *AdminUserContext {
	return c.Users.User(blockerID)
}

// Conversations
type AdminConversationsResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminConversationsResource(a *admin.PaanjAdmin) *AdminConversationsResource {
	return &AdminConversationsResource{admin: a}
}

func (r *AdminConversationsResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[k] = v
	}

	if memberIDs, ok := data["memberIds"]; ok {
		switch ids := memberIDs.(type) {
		case []string:
			members := make([]map[string]interface{}, 0, len(ids))
			for _, id := range ids {
				members = append(members, map[string]interface{}{
					"userId": normalizeUserID(id),
					"role":   "member",
				})
			}
			payload["members"] = members
		case []interface{}:
			members := make([]map[string]interface{}, 0, len(ids))
			for _, id := range ids {
				members = append(members, map[string]interface{}{
					"userId": normalizeUserID(id),
					"role":   "member",
				})
			}
			payload["members"] = members
		}
		delete(payload, "memberIds")
	}

	if members, ok := data["members"]; ok {
		payload["members"] = normalizeMembers(members)
	}

	return r.admin.GetHttpClient().Request("POST", "/admin/conversations", payload)
}

func (r *AdminConversationsResource) List(filters ...map[string]interface{}) (map[string]interface{}, error) {
	path := "/admin/conversations"
	if len(filters) > 0 && filters[0] != nil {
		query := buildConversationQuery(filters[0])
		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}
	}

	return r.admin.GetHttpClient().Request("GET", path, nil)
}

func (r *AdminConversationsResource) Get(conversationID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/conversations/%s", conversationID), nil)
}

func (r *AdminConversationsResource) Update(conversationID string, updates map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("PATCH", fmt.Sprintf("/admin/conversations/%s", conversationID), updates)
}

func (r *AdminConversationsResource) Delete(conversationID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("DELETE", fmt.Sprintf("/admin/conversations/%s", conversationID), nil)
}

func (r *AdminConversationsResource) OnCreate(callback func(interface{})) {
	r.subscribeGlobal([]string{"conversation.create"})
	r.admin.On("conversation.create", callback)
}

func (r *AdminConversationsResource) OnUpdate(callback func(interface{})) {
	r.subscribeGlobal([]string{"conversation.update"})
	r.admin.On("conversation.update", callback)
}

func (r *AdminConversationsResource) OnDelete(callback func(interface{})) {
	r.subscribeGlobal([]string{"conversation.delete"})
	r.admin.On("conversation.delete", callback)
}

func (r *AdminConversationsResource) Conversation(conversationID string) *AdminConversationContext {
	return &AdminConversationContext{
		conversationID: conversationID,
		resource:       r,
	}
}

func (r *AdminConversationsResource) subscribeGlobal(events []string) {
	_ = r.admin.Subscribe(map[string]interface{}{
		"type":     "admin.subscribe",
		"resource": "global",
		"events":   events,
	})
}

type AdminConversationContext struct {
	conversationID string
	resource       *AdminConversationsResource
}

func (c *AdminConversationContext) Send(content string, metadata ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"content": content,
	}
	if len(metadata) > 0 && metadata[0] != nil {
		payload["metadata"] = metadata[0]
	}

	return c.resource.admin.GetHttpClient().Request(
		"POST",
		fmt.Sprintf("/admin/conversations/%s/messages", c.conversationID),
		payload,
	)
}

func (c *AdminConversationContext) AddParticipant(userID string, role ...string) (map[string]interface{}, error) {
	participantRole := "member"
	if len(role) > 0 && role[0] != "" {
		participantRole = role[0]
	}

	return c.resource.admin.GetHttpClient().Request(
		"POST",
		fmt.Sprintf("/admin/conversations/%s/participants", c.conversationID),
		map[string]interface{}{
			"userId": normalizeUserID(userID),
			"role":   participantRole,
		},
	)
}

func (c *AdminConversationContext) RemoveParticipant(userID string) (map[string]interface{}, error) {
	return c.resource.admin.GetHttpClient().Request(
		"DELETE",
		fmt.Sprintf("/admin/conversations/%s/participants/%s", c.conversationID, userID),
		nil,
	)
}

func (c *AdminConversationContext) OnMessage(callback func(interface{})) {
	_ = c.resource.admin.Subscribe(map[string]interface{}{
		"type":     "admin.subscribe",
		"resource": "conversation",
		"id":       c.conversationID,
		"events":   []string{"message.create"},
	})
	c.resource.admin.On(fmt.Sprintf("conversation:%s:message.create", c.conversationID), callback)
}

// Users
type AdminUsersResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminUsersResource(a *admin.PaanjAdmin) *AdminUsersResource {
	return &AdminUsersResource{admin: a}
}

func (r *AdminUsersResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", "/admin/users", data)
}

func (r *AdminUsersResource) Get(userID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/users/%s", userID), nil)
}

func (r *AdminUsersResource) Update(userID string, updates map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("PATCH", fmt.Sprintf("/admin/users/%s", userID), updates)
}

func (r *AdminUsersResource) Delete(userID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("DELETE", fmt.Sprintf("/admin/users/%s", userID), nil)
}

func (r *AdminUsersResource) Block(blockerID, blockedID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/users/%s/block", blockerID), map[string]interface{}{
		"blockedUserId": normalizeUserID(blockedID),
	})
}

func (r *AdminUsersResource) Unblock(blockerID, blockedID string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/users/%s/unblock", blockerID), map[string]interface{}{
		"blockedUserId": normalizeUserID(blockedID),
	})
}

func (r *AdminUsersResource) User(blockerID string) *AdminUserContext {
	return &AdminUserContext{
		blockerID: blockerID,
		resource:  r,
	}
}

func (r *AdminUsersResource) OnCreate(callback func(interface{})) {
	r.subscribeGlobal([]string{"user.create"})
	r.admin.On("user.create", callback)
}

func (r *AdminUsersResource) OnUpdate(callback func(interface{})) {
	r.subscribeGlobal([]string{"user.update"})
	r.admin.On("user.update", callback)
}

func (r *AdminUsersResource) OnDelete(callback func(interface{})) {
	r.subscribeGlobal([]string{"user.delete"})
	r.admin.On("user.delete", callback)
}

func (r *AdminUsersResource) subscribeGlobal(events []string) {
	_ = r.admin.Subscribe(map[string]interface{}{
		"type":     "admin.subscribe",
		"resource": "global",
		"events":   events,
	})
}

type AdminUserContext struct {
	blockerID string
	resource  *AdminUsersResource
}

func (u *AdminUserContext) Block(blockedID string) (map[string]interface{}, error) {
	return u.resource.Block(u.blockerID, blockedID)
}

func (u *AdminUserContext) Unblock(blockedID string) (map[string]interface{}, error) {
	return u.resource.Unblock(u.blockerID, blockedID)
}

// Messages
type AdminMessagesResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminMessagesResource(a *admin.PaanjAdmin) *AdminMessagesResource {
	return &AdminMessagesResource{admin: a}
}

func (r *AdminMessagesResource) OnCreate(callback func(interface{})) {
	r.subscribeGlobal([]string{"message.create"})
	r.admin.On("message.create", callback)
}

func (r *AdminMessagesResource) OnSend(callback func(interface{})) {
	r.subscribeGlobal([]string{"message.send"})
	r.admin.On("message.send", callback)
}

func (r *AdminMessagesResource) OnUpdate(callback func(interface{})) {
	r.subscribeGlobal([]string{"message.update"})
	r.admin.On("message.update", callback)
}

func (r *AdminMessagesResource) OnDelete(callback func(interface{})) {
	r.subscribeGlobal([]string{"message.delete"})
	r.admin.On("message.delete", callback)
}

func (r *AdminMessagesResource) subscribeGlobal(events []string) {
	_ = r.admin.Subscribe(map[string]interface{}{
		"type":     "admin.subscribe",
		"resource": "global",
		"events":   events,
	})
}

func normalizeMembers(members interface{}) []map[string]interface{} {
	switch m := members.(type) {
	case []map[string]interface{}:
		normalized := make([]map[string]interface{}, 0, len(m))
		for _, member := range m {
			userID, hasUserID := member["userId"]
			if !hasUserID {
				continue
			}

			role := "member"
			if rawRole, ok := member["role"].(string); ok && rawRole != "" {
				role = rawRole
			}

			normalized = append(normalized, map[string]interface{}{
				"userId": normalizeUserID(userID),
				"role":   role,
			})
		}
		return normalized
	case []interface{}:
		normalized := make([]map[string]interface{}, 0, len(m))
		for _, raw := range m {
			member, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			userID, hasUserID := member["userId"]
			if !hasUserID {
				continue
			}
			role := "member"
			if rawRole, ok := member["role"].(string); ok && rawRole != "" {
				role = rawRole
			}
			normalized = append(normalized, map[string]interface{}{
				"userId": normalizeUserID(userID),
				"role":   role,
			})
		}
		return normalized
	default:
		return nil
	}
}

func normalizeUserID(userID interface{}) interface{} {
	switch v := userID.(type) {
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
		return v
	case float64:
		return int(v)
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return userID
	}
}

func buildConversationQuery(filters map[string]interface{}) string {
	values := url.Values{}
	if rawLimit, ok := filters["limit"]; ok {
		appendIntQuery(values, "limit", rawLimit)
	}
	if rawOffset, ok := filters["offset"]; ok {
		appendIntQuery(values, "offset", rawOffset)
	}
	if rawUserID, ok := filters["userId"]; ok {
		appendStringQuery(values, "userId", rawUserID)
	}
	return values.Encode()
}

func appendIntQuery(values url.Values, key string, value interface{}) {
	switch v := value.(type) {
	case int:
		values.Set(key, strconv.Itoa(v))
	case int32:
		values.Set(key, strconv.Itoa(int(v)))
	case int64:
		values.Set(key, strconv.FormatInt(v, 10))
	case float64:
		values.Set(key, strconv.Itoa(int(v)))
	case string:
		if _, err := strconv.Atoi(v); err == nil {
			values.Set(key, v)
		}
	}
}

func appendStringQuery(values url.Values, key string, value interface{}) {
	if str, ok := value.(string); ok && str != "" {
		values.Set(key, str)
	}
}
