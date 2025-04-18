const API_BASE = 'http://localhost:8080'

export async function loginAPI(username, password) {
  const response = await fetch(`${API_BASE}/auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Login failed: ${errorText}`);
  }
  return response.json();
}

export async function registerAPI(username, email, password) {
  const response = await fetch(`${API_BASE}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Registration failed: ${errorText}`);
  }
  return response.json();
}

export async function getUserChats(userId) {
  const response = await fetch(`${API_BASE}/getChats?user_id=${userId}`, { method: 'GET' });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch chats: ${errorText}`);
  }
  return response.json();
}

export async function createChatAPI(type, name) {
  const storedUser = localStorage.getItem('user');
  if (!storedUser) throw new Error('User not found in localStorage');
  const user = JSON.parse(storedUser);
  const response = await fetch(`${API_BASE}/createChat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    // Отправляем creator_id вместе с type и name, как требует серверная логика
    body: JSON.stringify({ type, name, creator_id: user.id })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to create chat: ${errorText}`);
  }
  return response.json();
}

export async function deleteChatAPI(chatId) {
  const response = await fetch(`${API_BASE}/deleteChat`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: chatId })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to delete chat: ${errorText}`);
  }
  return response.text();
}

export async function addMemberAPI(conversationId, userId, role) {
  const response = await fetch(`${API_BASE}/addMember`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId, role })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to add member: ${errorText}`);
  }
  return response.text();
}

export async function removeMemberAPI(conversationId, userId) {
  const response = await fetch(`${API_BASE}/removeMember`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to remove member: ${errorText}`);
  }
  return response.text();
}

export async function getMessagesAPI(chatId, limit) {
  const response = await fetch(`${API_BASE}/getMessage?chat_id=${chatId}&limit=${limit}`, { method: 'GET' });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch messages: ${errorText}`);
  }
  return response.json();
}

export async function findUserAPI(query) {
  const response = await fetch(`${API_BASE}/findUser?username=${encodeURIComponent(query)}`, { method: 'GET' });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to search user: ${errorText}`);
  }
  return response.json();
}

export async function getUserByIdAPI(userId) {
  const response = await fetch(`${API_BASE}/findUserbyID?user_id=${userId}`, { method: 'GET' });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch user info: ${errorText}`);
  }
  return response.json();
}

export async function getChatUsersAPI(chatId) {
  const response = await fetch(`${API_BASE}/getUsersFromChat?chat_id=${chatId}`, { method: 'GET' });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch chat users: ${errorText}`);
  }
  return response.json();
}

export async function updateUserAPI(id, name, email, password) {
  const response = await fetch(`${API_BASE}/updateUser`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id, username: name, email, password })
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to update user: ${errorText}`);
  }
  return response.text();
}

export async function leaveChatAPI(chatId, userId) {
  const response = await fetch(`${API_BASE}/leaveChat`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: chatId, user_id: userId })
  })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Failed to leave chat: ${errorText}`)
  }
  return response.text()
}

export async function getUnreadCountsAPI(userId) {
  const response = await fetch(`${API_BASE}/unreadCounts?user_id=${userId}`, { method: 'GET' })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Failed to fetch unread counts: ${errorText}`)
  }
  return response.json() // { [chatId]: count }
}

export async function resetUnreadAPI(chatId, userId) {
  const response = await fetch(`${API_BASE}/resetUnread`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ chat_id: chatId, user_id: userId })
  })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Failed to reset unread counter: ${errorText}`)
  }
  return response.text()
}