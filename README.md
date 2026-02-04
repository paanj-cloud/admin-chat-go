# Paanj Chat Admin SDK for Go

Official Go Chat Admin SDK for Paanj - Server-side chat administration and management.

[![Go Reference](https://pkg.go.dev/badge/github.com/paanj-cloud/admin-chat-go.svg)](https://pkg.go.dev/github.com/paanj-cloud/admin-chat-go)

## Installation

```bash
go get github.com/paanj-cloud/admin-go@latest
go get github.com/paanj-cloud/admin-chat-go@latest
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    admin "github.com/paanj-cloud/admin-go"
    chatadmin "github.com/paanj-cloud/admin-chat-go"
)

func main() {
    // Initialize admin client
    paanjAdmin := admin.NewAdmin("sk_live_your_secret_key", admin.AdminOptions{
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })

    // Connect
    paanjAdmin.Connect()
    defer paanjAdmin.Disconnect()

    // Initialize chat admin
    chatAdmin := chatadmin.NewAdminChat(paanjAdmin)

    // Get conversation
    conv, err := chatAdmin.Conversations.Get("conversation_id")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Conversation: %+v\\n", conv)
}
```

## Features

- ✅ Server-side conversation management
- ✅ User administration
- ✅ Message monitoring
- ✅ Real-time admin events
- ✅ Moderation tools

## Complete Examples

### Example 1: Monitor All Conversations

```go
package main

import (
    "fmt"
    "log"
    
    admin "github.com/paanj-cloud/admin-go"
    chatadmin "github.com/paanj-cloud/admin-chat-go"
)

func main() {
    // Initialize
    paanjAdmin := admin.NewAdmin("sk_live_your_secret_key", admin.AdminOptions{
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })

    if err := paanjAdmin.Connect(); err != nil {
        log.Fatal(err)
    }
    defer paanjAdmin.Disconnect()

    chatAdmin := chatadmin.NewAdminChat(paanjAdmin)

    // Listen for new messages (admin can see all)
    chatAdmin.Messages.OnCreate(func(data interface{}) {
        msg := data.(map[string]interface{})
        fmt.Printf("📩 Message in %s: %s\\n", 
            msg["conversationId"], 
            msg["content"])
    })

    // List all conversations
    conversations, err := chatAdmin.Conversations.List()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total conversations: %d\\n", len(conversations.([]interface{})))

    // Keep alive
    select {}
}
```

### Example 2: Create and Manage Conversations

```go
package main

import (
    "fmt"
    "log"
    
    admin "github.com/paanj-cloud/admin-go"
    chatadmin "github.com/paanj-cloud/admin-chat-go"
)

func main() {
    // Initialize
    paanjAdmin := admin.NewAdmin("sk_live_your_secret_key", admin.AdminOptions{
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    paanjAdmin.Connect()
    defer paanjAdmin.Disconnect()

    chatAdmin := chatadmin.NewAdminChat(paanjAdmin)

    // Create conversation as admin
    conv, err := chatAdmin.Conversations.Create(map[string]interface{}{
        "name": "Official Announcement Channel",
        "members": []map[string]interface{}{
            {"userId": 123, "role": "admin"},
            {"userId": 456, "role": "member"},
        },
        "metadata": map[string]interface{}{
            "type": "announcement",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    conversationId := conv["id"].(string)
    fmt.Printf("✅ Conversation created: %s\\n", conversationId)

    // Get conversation details
    details, _ := chatAdmin.Conversations.Get(conversationId)
    fmt.Printf("Conversation details: %+v\\n", details)

    // List all conversations
    allConversations, _ := chatAdmin.Conversations.List()
    fmt.Printf("Total conversations: %d\\n", len(allConversations.(map[string]interface{})))
}
```

### Example 3: User Moderation

```go
package main

import (
    "fmt"
    "log"
    
    admin "github.com/paanj-cloud/admin-go"
    chatadmin "github.com/paanj-cloud/admin-chat-go"
)

func main() {
    // Initialize
    paanjAdmin := admin.NewAdmin("sk_live_your_secret_key", admin.AdminOptions{
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    paanjAdmin.Connect()
    defer paanjAdmin.Disconnect()

    chatAdmin := chatadmin.NewAdminChat(paanjAdmin)

    // Create a user
    user, err := chatAdmin.Users.Create(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "metadata": map[string]interface{}{
            "role": "moderator",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    userId := user["userId"].(string)
    fmt.Printf("✅ User created: %s\\n", userId)

    // Get user details
    userDetails, _ := chatAdmin.Users.Get(userId)
    fmt.Printf("User: %+v\\n", userDetails)

    // Block a user (moderation)
    result, err := chatAdmin.Users.Block("blocker_id", "blocked_user_id")
    if err != nil {
        log.Printf("Block failed: %v", err)
    } else {
        fmt.Println("✅ User blocked successfully")
    }
}
```

### Example 4: Real-time Monitoring Dashboard

```go
package main

import (
    "fmt"
    "log"
    "sync"
    
    admin "github.com/paanj-cloud/admin-go"
    chatadmin "github.com/paanj-cloud/admin-chat-go"
)

func main() {
    // Initialize
    paanjAdmin := admin.NewAdmin("sk_live_your_secret_key", admin.AdminOptions{
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    paanjAdmin.Connect()
    defer paanjAdmin.Disconnect()

    chatAdmin := chatadmin.NewAdminChat(paanjAdmin)

    // Statistics
    var mu sync.Mutex
    stats := struct {
        TotalMessages      int
        TotalConversations int
        ActiveUsers        map[string]bool
    }{
        ActiveUsers: make(map[string]bool),
    }

    // Monitor all messages
    chatAdmin.Messages.OnCreate(func(data interface{}) {
        msg := data.(map[string]interface{})
        mu.Lock()
        stats.TotalMessages++
        if senderId, ok := msg["senderId"].(string); ok {
            stats.ActiveUsers[senderId] = true
        }
        mu.Unlock()

        fmt.Printf("📊 Stats - Messages: %d, Active Users: %d\\n",
            stats.TotalMessages,
            len(stats.ActiveUsers))
    })

    // Subscribe to conversation events
    paanjAdmin.Subscribe(map[string]interface{}{
        "resource": "conversations",
        "events":   []string{"conversation.created", "conversation.updated"},
    })

    paanjAdmin.On("conversation.created", func(data interface{}) {
        mu.Lock()
        stats.TotalConversations++
        mu.Unlock()
        fmt.Printf("🆕 New conversation created\\n")
    })

    // Initial data load
    conversations, _ := chatAdmin.Conversations.List()
    if convMap, ok := conversations.(map[string]interface{}); ok {
        if convList, ok := convMap["conversations"].([]interface{}); ok {
            stats.TotalConversations = len(convList)
        }
    }

    fmt.Println("📊 Monitoring dashboard started...")
    select {}
}
```

## API Reference

### Conversations

#### Create Conversation

```go
conv, err := chatAdmin.Conversations.Create(map[string]interface{}{
    "name": "Conversation Name",
    "members": []map[string]interface{}{
        {"userId": 123, "role": "admin"},
    },
})
```

#### List All Conversations

```go
conversations, err := chatAdmin.Conversations.List()
```

#### Get Conversation

```go
conv, err := chatAdmin.Conversations.Get(conversationId)
```

### Users

#### Create User

```go
user, err := chatAdmin.Users.Create(map[string]interface{}{
    "name":  "User Name",
    "email": "user@example.com",
})
```

#### Get User

```go
user, err := chatAdmin.Users.Get(userId)
```

#### Block User

```go
result, err := chatAdmin.Users.Block(blockerId, blockedId)
```

### Messages

#### Listen for All Messages

```go
chatAdmin.Messages.OnCreate(func(data interface{}) {
    msg := data.(map[string]interface{})
    fmt.Printf("Message: %s\\n", msg["content"])
})
```

## Event Types

| Event | Description | Data |
|-------|-------------|------|
| `message.create` | New message (all conversations) | `{content, senderId, conversationId, ...}` |
| `conversation.created` | New conversation created | `{conversationId, name, ...}` |
| `conversation.updated` | Conversation updated | `{conversationId, ...}` |
| `user.created` | New user created | `{userId, name, ...}` |

## Security Best Practices

1. **Never expose secret keys** in client-side code or version control
2. **Use environment variables** for credentials
3. **Implement rate limiting** for admin operations
4. **Log all admin actions** for audit trails
5. **Use HTTPS/WSS** in production

```go
import "os"

secretKey := os.Getenv("PAANJ_SECRET_KEY")
if secretKey == "" {
    log.Fatal("PAANJ_SECRET_KEY environment variable not set")
}

paanjAdmin := admin.NewAdmin(secretKey, admin.AdminOptions{
    ApiUrl: "https://api1.paanj.com",
    WsUrl:  "wss://ws1.paanj.com",
})
```

## Error Handling

```go
conv, err := chatAdmin.Conversations.Get(conversationId)
if err != nil {
    log.Printf("Failed to get conversation: %v", err)
    return
}

user, err := chatAdmin.Users.Create(userData)
if err != nil {
    log.Printf("Failed to create user: %v", err)
    return
}
```

## Use Cases

- **Moderation Dashboard**: Monitor all conversations and messages
- **Analytics**: Track user activity and conversation metrics
- **User Management**: Create, update, and moderate users
- **Automated Moderation**: Implement content filtering and user blocking
- **Admin Notifications**: Get real-time alerts for important events

## License

MIT License - see LICENSE file for details.

## Support

- Documentation: https://docs.paanj.com
- Issues: https://github.com/paanj-cloud/admin-chat-go/issues
- Email: support@paanj.com
