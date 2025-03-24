// Просто функции для ручек
const API_BASE = 'http://localhost:8080'

export async function loginAPI(username, password) {
  const response = await fetch(`${API_BASE}/auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  })
  if (!response.ok) throw new Error('Login failed')
  return response.json()
}

export async function registerAPI(username, email, password) {
  const response = await fetch(`${API_BASE}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password })
  })
  if (!response.ok) throw new Error('Registration failed')
  return response.json()
}

export async function getUserChats(userId) {
  const response = await fetch(`${API_BASE}/getChats?user_id=${userId}`, { method: 'GET' })
  if (!response.ok) throw new Error('Failed to fetch chats')
  return response.json()
}

export async function createChatAPI(type, name) {
  const response = await fetch(`${API_BASE}/createChat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ type, name })
  })
  if (!response.ok) throw new Error('Failed to create chat')
  return response.json()
}

export async function deleteChatAPI(chatId) {
  const response = await fetch(`${API_BASE}/deleteChat`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: chatId })
  })
  if (!response.ok) throw new Error('Failed to delete chat')
  return response.text()
}

export async function addMemberAPI(conversationId, userId, role) {
  const response = await fetch(`${API_BASE}/addMember`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId, role })
  })
  if (!response.ok) throw new Error('Failed to add member')
  return response.text()
}

export async function removeMemberAPI(conversationId, userId) {
  const response = await fetch(`${API_BASE}/removeMember`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId })
  })
  if (!response.ok) throw new Error('Failed to remove member')
  return response.text()
}

export async function deleteMessageAPI(messageId, senderId) {
  const response = await fetch(`${API_BASE}/deleteMessage`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message_id: messageId, sender_id: senderId })
  })
  if (!response.ok) throw new Error('Failed to delete message')
  return response.text()
}

export async function getMessagesAPI(chatId, limit) {
  const response = await fetch(`${API_BASE}/getMessage?chat_id=${chatId}&limit=${limit}`, { method: 'GET' })
  if (!response.ok) throw new Error('Failed to fetch messages')
  return response.json()
}
