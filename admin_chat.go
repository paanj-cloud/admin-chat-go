package chatadmin

// Test sync: 2026-02-04

import (
	"fmt"

	"github.com/paanj-cloud/admin-go"
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
	return r.admin.GetHttpClient().Request("POST", "/admin/conversations", data)
}

func (r *AdminConversationsResource) List() (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", "/admin/conversations", nil)
}

func (r *AdminConversationsResource) Get(conversationId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/conversations/%s", conversationId), nil)
}

func (r *AdminConversationsResource) Update(conversationId string, data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("PATCH", fmt.Sprintf("/admin/conversations/%s", conversationId), data)
}

func (r *AdminConversationsResource) Delete(conversationId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("DELETE", fmt.Sprintf("/admin/conversations/%s", conversationId), nil)
}

func (r *AdminConversationsResource) AddParticipant(conversationId string, data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/conversations/%s/participants", conversationId), data)
}

func (r *AdminConversationsResource) RemoveParticipant(conversationId, userId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("DELETE", fmt.Sprintf("/admin/conversations/%s/participants/%s", conversationId, userId), nil)
}

func (r *AdminConversationsResource) SendMessage(conversationId string, data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/conversations/%s/messages", conversationId), data)
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

func (r *AdminUsersResource) Get(userId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/users/%s", userId), nil)
}

func (r *AdminUsersResource) Block(blockerId, blockedId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/users/%s/block", blockedId), map[string]interface{}{
		"blockerId": blockerId,
	})
}

func (r *AdminUsersResource) Unblock(userId, blockedId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("POST", fmt.Sprintf("/admin/users/%s/unblock", blockedId), nil)
}

func (r *AdminUsersResource) Update(userId string, data map[string]interface{}) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("PATCH", fmt.Sprintf("/admin/users/%s", userId), data)
}

func (r *AdminUsersResource) Delete(userId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("DELETE", fmt.Sprintf("/admin/users/%s", userId), nil)
}

func (r *AdminUsersResource) GetConversations(userId string) (map[string]interface{}, error) {
	return r.admin.GetHttpClient().Request("GET", fmt.Sprintf("/admin/users/%s/conversations", userId), nil)
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
