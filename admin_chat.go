package chatadmin

import (
	"fmt"

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

// Resources (Simplified)

// Conversations
type AdminConversationsResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminConversationsResource(a *admin.PaanjAdmin) *AdminConversationsResource {
	return &AdminConversationsResource{admin: a}
}

func (r *AdminConversationsResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", "/api/v1/admin/conversations", data)
}

func (r *AdminConversationsResource) List() (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", "/api/v1/admin/conversations", nil)
}

func (r *AdminConversationsResource) Get(conversationId string) (map[string]interface{}, error) {
	// JS SDK: GET /admin/conversations/:id (not /api/v1/admin/...)
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/conversations/%s", conversationId), nil)
}

// Users
type AdminUsersResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminUsersResource(a *admin.PaanjAdmin) *AdminUsersResource {
	return &AdminUsersResource{admin: a}
}

func (r *AdminUsersResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", "/api/v1/admin/users", data)
}

func (r *AdminUsersResource) Get(userId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/api/v1/admin/users/%s", userId), nil)
}

func (r *AdminUsersResource) Block(blockerId, blockedId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", "/api/v1/admin/users/block", map[string]interface{}{
		"blockerId": blockerId,
		"blockedId": blockedId,
	})
}

// Messages
type AdminMessagesResource struct {
	admin *admin.PaanjAdmin
}

func NewAdminMessagesResource(a *admin.PaanjAdmin) *AdminMessagesResource {
	return &AdminMessagesResource{admin: a}
}

func (r *AdminMessagesResource) OnCreate(callback func(interface{})) {
	r.admin.On("message.create", callback)
}
