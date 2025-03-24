<template>
  <div class="chat-container">
    <!-- Левый сайдбар: список чатов и управление -->
    <div class="sidebar">
      <h2>Your Chats</h2>
      <ul>
        <li v-for="chat in chats" :key="chat.id" @click="selectChat(chat)" :class="{ active: selectedChat && selectedChat.id === chat.id }">
          {{ chat.name }}
          <button @click.stop="deleteChat(chat.id)">Delete</button>
        </li>
      </ul>
      <div class="create-chat">
        <h3>Create Chat</h3>
        <form @submit.prevent="createChat">
          <div>
            <label for="chatType">Type:</label>
            <select v-model="newChatType" id="chatType">
              <option value="private">Private</option>
              <option value="group">Group</option>
            </select>
          </div>
          <div>
            <label for="chatName">Name:</label>
            <input v-model="newChatName" type="text" id="chatName" required />
          </div>
          <button type="submit">Create</button>
        </form>
      </div>
      <div class="member-management" v-if="selectedChat">
        <h3>Manage Chat Members</h3>
        <form @submit.prevent="addMember">
          <div>
            <label for="addUserId">Add Member (User ID):</label>
            <input v-model="memberToAdd" type="number" id="addUserId" required />
          </div>
          <div>
            <label for="memberRole">Role:</label>
            <input v-model="memberRole" type="text" id="memberRole" required placeholder="e.g., member" />
          </div>
          <button type="submit">Add Member</button>
        </form>
        <form @submit.prevent="removeMember">
          <div>
            <label for="removeUserId">Remove Member (User ID):</label>
            <input v-model="memberToRemove" type="number" id="removeUserId" required />
          </div>
          <button type="submit">Remove Member</button>
        </form>
      </div>
    </div>

    <!-- Основная область: окно чата -->
    <div class="chat-window" v-if="selectedChat">
      <h2>{{ selectedChat.name }}</h2>
      <div class="message-list">
        <MessageList :messages="messages" />
      </div>
      <div class="message-input">
        <MessageInput @sendMessage="sendMessage" />
      </div>
    </div>
    <div v-else class="no-chat">
      <p>Select a chat to view messages</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useUserStore } from '../store'
import { useRouter } from 'vue-router'
import MessageList from '../components/MessageList.vue'
import MessageInput from '../components/MessageInput.vue'
import {
  getUserChats,
  createChatAPI,
  deleteChatAPI,
  addMemberAPI,
  removeMemberAPI,
  // добавим новую функцию для получения сообщений
  getMessagesAPI
} from '../api/api'

const WS_URL = 'ws://localhost:8080/ws'

const userStore = useUserStore()
const router = useRouter()

const chats = ref([])
const selectedChat = ref(null)
const messages = ref([])

const newChatType = ref('private')
const newChatName = ref('')
const memberToAdd = ref('')
const memberRole = ref('member')
const memberToRemove = ref('')

let ws = null

const fetchChats = async () => {
  try {
    const data = await getUserChats(userStore.user.id)
    chats.value = data
    if (!selectedChat.value && chats.value.length > 0) {
      selectChat(chats.value[0])
    }
  } catch (error) {
    alert(error.message)
  }
}

// Загружаем последние 20 сообщений выбранного чата
const fetchMessages = async (chatId) => {
  try {
    const msgs = await getMessagesAPI(chatId, 20)
    messages.value = msgs.reverse() // разворачиваем, чтобы старые были вверху
  } catch (error) {
    alert(error.message)
  }
}

const createChat = async () => {
  try {
    const newChat = await createChatAPI(newChatType.value, newChatName.value)
    chats.value.push(newChat)
    newChatName.value = ''
    selectChat(newChat)
  } catch (error) {
    alert(error.message)
  }
}

const deleteChat = async (chatId) => {
  try {
    await deleteChatAPI(chatId)
    chats.value = chats.value.filter(chat => chat.id !== chatId)
    if (selectedChat.value && selectedChat.value.id === chatId) {
      selectedChat.value = chats.value.length > 0 ? chats.value[0] : null
      messages.value = []
      if (ws) ws.close()
      if (selectedChat.value) {
        fetchMessages(selectedChat.value.id)
        connectWebSocket()
      }
    }
  } catch (error) {
    alert(error.message)
  }
}

const addMember = async () => {
  if (!selectedChat.value) return
  try {
    await addMemberAPI(selectedChat.value.id, parseInt(memberToAdd.value), memberRole.value)
    alert('Member added successfully')
    memberToAdd.value = ''
  } catch (error) {
    alert(error.message)
  }
}

const removeMember = async () => {
  if (!selectedChat.value) return
  try {
    await removeMemberAPI(selectedChat.value.id, parseInt(memberToRemove.value))
    alert('Member removed successfully')
    memberToRemove.value = ''
  } catch (error) {
    alert(error.message)
  }
}

const connectWebSocket = () => {
  if (!selectedChat.value) return
  ws = new WebSocket(`${WS_URL}?chat_id=${selectedChat.value.id}&user_id=${userStore.user.id}`)
  ws.onopen = () => console.log('WebSocket connected')
  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      messages.value.push(msg)
    } catch (e) {
      console.error('Error parsing WS message:', e)
    }
  }
  ws.onerror = (error) => console.error('WebSocket error:', error)
  ws.onclose = (event) => {
    console.log('WebSocket closed:', event)
  }
}

const sendMessage = (text) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    const messageData = { chat_id: selectedChat.value.id, content: text, sender_id: userStore.user.id }
    ws.send(JSON.stringify(messageData))
  }
}

const selectChat = async (chat) => {
  selectedChat.value = chat
  messages.value = []
  if (ws) ws.close()
  // Сначала загружаем последние 20 сообщений
  await fetchMessages(chat.id)
  connectWebSocket()
}

onMounted(() => {
  if (!userStore.user) {
    router.push('/')
    return
  }
  fetchChats()
})

onBeforeUnmount(() => {
  if (ws) ws.close()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
}

/* Sidebar */
.sidebar {
  width: 300px;
  padding: 10px;
  border-right: 1px solid #ccc;
  overflow-y: auto;
  background-color: #f9f9f9;
}
.sidebar h2 {
  margin-bottom: 10px;
  font-size: 1.8rem;
}
.sidebar ul {
  list-style: none;
  padding: 0;
}
.sidebar li {
  padding: 5px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.sidebar li.active {
  background-color: #e0e0e0;
}
.sidebar button {
  background-color: #dc3545;
  border: none;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}

/* Chat window */
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}
.chat-window h2 {
  margin: 0 0 10px;
}
.message-list {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #ccc;
  padding: 10px;
  margin-bottom: 10px;
}
.message-input {
  border-top: 1px solid #ccc;
  padding: 10px;
}

/* Create chat & member management */
.create-chat, .member-management {
  margin-top: 20px;
  border-top: 1px solid #ccc;
  padding-top: 10px;
}
</style>
