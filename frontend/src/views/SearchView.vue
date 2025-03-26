<template>
  <div class="search-container">
    <h1>Search Users</h1>
    <input v-model="query" @keyup.enter="search" placeholder="Enter username" />
    <button @click="search">Search</button>
    <ul>
      <li v-for="user in results" :key="user.id" @click="selectUser(user)">
        {{ user.username }} ({{ user.email }})
      </li>
    </ul>
    <button @click="goBack">Back</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { findUserAPI, addMemberAPI } from '../api/api'
const router = useRouter()
const route = useRoute()
const query = ref('')
const results = ref([])

const search = async () => {
  if (query.value.trim() === '') return
  try {
    results.value = await findUserAPI(query.value)
  } catch (error) {
    alert(error.message)
  }
}

const selectUser = async (user) => {
  // Если есть параметр chat_id, значит мы добавляем пользователя в чат
  const chatId = route.query.chat_id
  if (chatId) {
    try {
      await addMemberAPI(parseInt(chatId), user.id, 'member')
      alert(`User ${user.username} added successfully`)
      router.push({ name: 'ChatInfo', params: { id: chatId } })
    } catch (error) {
      alert(error.message)
    }
  } else {
    // Иначе переходим на профиль пользователя
    router.push({ name: 'Profile', query: { user_id: user.id } })
  }
}

const goBack = () => {
  const chatId = route.query.chat_id
  if (chatId) {
    router.push({ name: 'ChatInfo', params: { id: chatId } })
  } else {
    router.push('/chat')
  }
}
</script>

<style scoped>
.search-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: white;
  padding: 20px;
}
input {
  padding: 10px;
  font-size: 1.1rem;
  width: 300px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
button {
  padding: 10px 20px;
  font-size: 1.1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 5px;
}
button:hover {
  background-color: #0056b3;
}
ul {
  list-style: none;
  padding: 0;
  width: 300px;
}
li {
  padding: 8px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  color: black;
}
li:hover {
  background-color: #f5f5f5;
}
</style>
