<template>
  <div class="main-container">
    <!-- Верхняя синяя панель -->
    <div class="top-bar">
      <div class="user-icon" @click="goToProfile">
        <span>{{ userInitial }}</span>
      </div>
      <div class="search-container">
        <input v-model="searchQuery" @input="handleSearchInput" @keyup.enter="executeSearch" placeholder="Search users..." />
        <SearchResults v-if="searchResults.length > 0" :results="searchResults" @select="selectSearchedUser" />
      </div>
      <div class="signout-icon" @click="toggleSignOutModal">
        <span>Sign Out</span>
      </div>
    </div>

    <!-- Модальное окно выхода -->
    <div v-if="showSignOutModal" class="modal-overlay" @click.self="toggleSignOutModal">
      <div class="modal">
        <p>Are you sure you want to sign out?</p>
        <button @click="signOut">Yes</button>
        <button @click="toggleSignOutModal">No</button>
      </div>
    </div>

    <div class="content">
      <!-- Левая колонка: список чатов -->
      <div class="sidebar">
        <h2>Your Chats</h2>
        <ChatList :chats="chats" @selectChat="selectChat" @openChatInfo="goToChatInfo" />
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
      </div>

      <!-- Правая колонка: окно чата -->
      <div class="chat-window" v-if="selectedChat">
        <div class="chat-header">
          <h2 @click="goToChatInfo(selectedChat.id)">{{ selectedChat.name }}</h2>
        </div>
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
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useUserStore } from '../store'
import { useRouter } from 'vue-router'
import ChatList from '../components/ChatList.vue'
import MessageList from '../components/MessageList.vue'
import MessageInput from '../components/MessageInput.vue'
import SearchResults from '../components/SearchResults.vue'
import {
  getUserChats,
  createChatAPI,
  getMessagesAPI,
  findUserAPI,
  getUnreadCountsAPI,
  resetUnreadAPI
  
} from '../api/api'

const WS_URL = 'ws://localhost:8080/ws'
const userStore = useUserStore()
const router = useRouter()

const chats = ref([])
const selectedChat = ref(null)
const messages = ref([])
const unreadCounts = ref({})

const newChatType = ref('private')
const newChatName = ref('')
const searchQuery = ref('')
const searchResults = ref([])

const ws = ref(null)
const showSignOutModal = ref(false)

const fetchChats = async () => {
  try {
    const allChats = await getUserChats(userStore.user.id)
    const counts = await getUnreadCountsAPI(userStore.user.id)
    unreadCounts.value = counts

    chats.value = allChats.map(chat => ({
      ...chat,
      unread_count: counts[chat.id] || 0
    }))
  } catch (error) {
    alert(error.message)
  }
}

const fetchMessages = async (chatId) => {
  try {
    const msgs = await getMessagesAPI(chatId, 20)
    messages.value = msgs.reverse()
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

const selectChat = async (chat) => {
  selectedChat.value = chat
  messages.value = []
  if (ws.value) ws.value.close()

  await fetchMessages(chat.id)
  await resetUnreadCount(chat.id)

  connectWebSocket()
}

const connectWebSocket = () => {
  if (!selectedChat.value) return
  ws.value = new WebSocket(`${WS_URL}?chat_id=${selectedChat.value.id}&user_id=${userStore.user.id}`)
  ws.value.onopen = () => console.log('WebSocket connected')
  ws.value.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (selectedChat.value && msg.chat_id === selectedChat.value.id) {
        messages.value.push(msg)
      } else {
        const targetChat = chats.value.find(c => c.id === msg.chat_id)
        if (targetChat) targetChat.unread_count = (targetChat.unread_count || 0) + 1
      }
    } catch (e) {
      console.error('Error parsing WS message:', e)
    }
  }
  ws.value.onerror = (error) => console.error('WebSocket error:', error)
  ws.value.onclose = (event) => console.log('WebSocket closed:', event)
}

const resetUnreadCount = async (chatId) => {
  try {
    await resetUnreadAPI(chatId, userStore.user.id)
    const chat = chats.value.find(c => c.id === chatId)
    if (chat) chat.unread_count = 0
  } catch (err) {
    console.error('Failed to reset unread count:', err)
  }
}

const sendMessage = (text) => {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    const messageData = {
      chat_id: selectedChat.value.id,
      content: text,
      sender_id: userStore.user.id
    }
    ws.value.send(JSON.stringify(messageData))
  }
}


const handleSearchInput = () => {
  if (searchQuery.value.trim().length === 0) {
    searchResults.value = []
    return
  }
  executeSearch()
}

const executeSearch = async () => {
  try {
    const results = await findUserAPI(searchQuery.value)
    searchResults.value = results
  } catch (error) {
    console.error(error)
  }
}

const selectSearchedUser = (user) => {
  // При выборе из верхней панели переходим на профиль пользователя
  router.push({ name: 'Profile', query: { user_id: user.id } })
  searchResults.value = []
}

const toggleSignOutModal = () => {
  showSignOutModal.value = !showSignOutModal.value
}

const signOut = () => {
  userStore.clearUser()
  router.push('/')
}

const goToProfile = () => {
  router.push({ name: 'Profile', query: { user_id: userStore.user.id } })
}

const goToChatInfo = (chatId) => {
  router.push({ name: 'ChatInfo', params: { id: chatId } })
}

const userInitial = computed(() => {
  return userStore.user && userStore.user.username ? userStore.user.username.charAt(0).toUpperCase() : '?'
})

onMounted(() => {
  if (!userStore.user) {
    router.push('/')
    return
  }
  fetchChats()
})

onBeforeUnmount(() => {
  if (ws.value) ws.value.close()
})
</script>

<style scoped>
.main-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

/* Верхняя панель */
.top-bar {
  background-color: #007bff;
  color: white;
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  position: relative;
}
.top-bar .user-icon,
.top-bar .signout-icon {
  cursor: pointer;
  font-size: 1.5rem;
  width: 60px;
  text-align: center;
}
.top-bar .search-container {
  flex: 1;
  display: flex;
  justify-content: center;
  position: relative;
}
.top-bar input {
  width: 60%;
  padding: 8px 12px;
  font-size: 1rem;
  border: none;
  border-radius: 4px;
}

/* Модальное окно выхода */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}
.modal {
  background: white;
  padding: 20px;
  border-radius: 6px;
  text-align: center;
}

/* Основной контент */
.content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* Sidebar */
.sidebar {
  width: 300px;
  padding: 10px;
  border-right: 1px solid #ccc;
  background-color: #f9f9f9;
  overflow-y: auto;
}
.sidebar h2 {
  margin-bottom: 10px;
  font-size: 1.8rem;
}

/* Chat window */
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}
.chat-header h2 {
  cursor: pointer;
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
</style>
