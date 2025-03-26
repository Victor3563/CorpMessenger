<template>
  <div class="chat-info-container">
    <div class="chat-info-block">
      <h1>Chat Info</h1>
      <h2>{{ chat.name }}</h2>
      <div class="participants">
        <h3>Participants</h3>
        <ul>
          <li v-for="user in participants" :key="user.id">
            {{ user.username }} ({{ user.role }})
            <button @click="removeParticipant(user.id)">Delete</button>
          </li>
        </ul>
      </div>
      <button @click="goToAddParticipants">Add Participants</button>
      <button @click="goBack">Back to Chat</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getChatUsersAPI, removeMemberAPI, getUserChats } from '../api/api'
import { useUserStore } from '../store'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const chat = ref({ id: null, name: '' })
const participants = ref([])

onMounted(async () => {
  chat.value.id = parseInt(route.params.id)
  const userChats = await getUserChats(userStore.user.id)
  const thisChat = userChats.find(c => c.id === chat.value.id)
  if (thisChat) {
    chat.value.name = thisChat.name
  }
  try {
    const data = await getChatUsersAPI(chat.value.id)
    participants.value = data
  } catch (error) {
    alert(error.message)
  }
})

const removeParticipant = async (userId) => {
  try {
    await removeMemberAPI(chat.value.id, userId)
    alert('Member removed successfully')
    // Обновляем список участников после удаления
    const data = await getChatUsersAPI(chat.value.id)
    participants.value = data
  } catch (error) {
    alert(error.message)
  }
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
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: white;
}
.chat-info-block {
  background-color: #f2f2f2;
  padding: 40px;
  border-radius: 8px;
  width: 100%;
  max-width: 600px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  text-align: center;
}
.participants ul {
  list-style: none;
  padding: 0;
}
.participants li {
  margin: 8px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
button {
  padding: 10px 20px;
  font-size: 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 10px 5px;
}
button:hover {
  background-color: #0056b3;
}
</style>
