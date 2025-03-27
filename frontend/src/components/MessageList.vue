<template>
  <div>
    <div v-for="(message, index) in messages" :key="message.id" class="message" :class="{ mine: message.sender_id === currentUser }">
      <div class="sender">{{ message.sender_id === currentUser ? 'You' : 'User ' + message.sender_id }}</div>
      <div class="content">
        <p v-for="(para, i) in splitContent(message.content)" :key="i">{{ para }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import { useUserStore } from '../store'
const props = defineProps({
  messages: { type: Array, required: true }
})
const userStore = useUserStore()
const currentUser = userStore.user ? userStore.user.id : null
const splitContent = (text) => {
    const chunks = []
    for (let i = 0; i < text.length; i += 70) {
      chunks.push(text.substring(i, i + 70))
    }
    return chunks
  }
</script>

<style scoped>
.message {
  margin-bottom: 10px;
  padding: 10px;
  border-radius: 6px;
  max-width: 70%;
}
.message.mine {
  align-self: flex-end;
  background-color: #d1e7dd;
}
.message:not(.mine) {
  align-self: flex-start;
  background-color: #f8d7da;
}
.sender {
  font-weight: bold;
  margin-bottom: 5px;
}
.content p {
  margin: 0;
}
</style>
