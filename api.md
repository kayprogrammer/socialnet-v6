## A Realtime Social Networking API built with FIBER & GORM ORM.

### WEBSOCKETS:

#### Notifications

- URL: `wss://{host}/api/v6/ws/notifications`

- Requires authorization, so pass in the Bearer Authorization header.

- You can only read and not send notification messages into this socket.


#### Chats

- URL: `wss://{host}/api/v6/ws/chats/{id}`
- Requires authorization, so pass in the Bearer Authorization header.
- Use chat_id as the ID for an existing chat or username if it's the first message in a DM.
- You cannot read realtime messages from a username that doesn't belong to the authorized user, but you can surely send messages.
- Only send a message to the socket endpoint after the message has been created or updated, and files have been uploaded.
- Fields when sending a message through the socket:

  ```json
  { "status": "CREATED", "id": "fe4e0235-80fc-4c94-b15e-3da63226f8ab" }
  ```
