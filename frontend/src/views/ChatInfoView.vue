<template>
  <div class="chat-info-container">
    <div class="chat-info-block">
      <h1>{{ chat.name }}</h1>
      <ul class="participants-list">
        <li v-for="user in participants" :key="user.id">
          {{ user.username }} ({{ user.role }})
          <button
            v-if="isAdmin && user.id !== userStore.user.id"
            @click="removeParticipant(user.id)"
          >
            Delete
          </button>
        </li>
      </ul>

      <div class="actions">
        <div v-if="isAdmin && !isPrivateChat" class="action-group">
          <button @click="goToAddParticipants">Add Participant</button>
        </div>

        <div v-if="isAdmin && !isPrivateChat" class="action-group">
          <button @click="deleteChat">Delete chat</button>
        </div>

        <div class="action-group">
          <button @click="leaveChat">Leave Chat</button>
          <button @click="goBack">Back</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../store'
import {
  getChatUsersAPI,
  removeMemberAPI,
  leaveChatAPI,
  getUserChats,
  deleteChatAPI
} from '../api/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const chat = ref({ id: parseInt(route.params.id), name: '', type: '' })
const participants = ref([])

onMounted(async () => {
  const chats = await getUserChats(userStore.user.id)
  const targetChat = chats.find(c => c.id === chat.value.id)
  if (targetChat) {
    chat.value.name = targetChat.name
    chat.value.type = targetChat.type
  }
  participants.value = await getChatUsersAPI(chat.value.id)
})

const isAdmin = computed(() => {
  const current = participants.value.find(u => u.id === userStore.user.id)
  return current?.role === 'admin'
})

const isPrivateChat = computed(() => {
  return chat.value.type === 'private' && participants.value.length >= 2
})

const removeParticipant = async (userId) => {
  await removeMemberAPI(chat.value.id, userId)
  participants.value = await getChatUsersAPI(chat.value.id)
}

const deleteChat = async () => {
  await deleteChatAPI(chat.value.id)
  router.push('/chat')
}

const leaveChat = async () => {
  await leaveChatAPI(chat.value.id, userStore.user.id)
  router.push('/chat')
}

const goToAddParticipants = () => {
  router.push({ name: 'Search', query: { chat_id: chat.value.id } })
}

const goBack = () => {
  router.push('/chat')
}
</script>

<style scoped>
.chat-info-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background-color: #f5f5f5;
}

.chat-info-block {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.1);
  padding: 32px;
  width: 100%;
  max-width: 600px;
}

h1 {
  text-align: center;
  margin-bottom: 24px;
  color: #333;
}

.participants-list {
  list-style: none;
  padding: 0;
  margin: 0 0 24px;
}

.participants-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin: 8px 0;
  background-color: #f8f9fa;
  border-radius: 6px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.action-group {
  display: flex;
  gap: 12px;
  justify-content: center;
}

button {
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
  transform: translateY(-1px);
}

button:active {
  transform: translateY(0);
}
</style>