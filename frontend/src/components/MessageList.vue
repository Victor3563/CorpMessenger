<!-- src/components/MessageList.vue -->
<template>
  <div>
    <div v-for="(message, index) in messages" :key="index" class="message">
      <strong>{{ message.sender_id === currentUser ? 'You' : 'User ' + message.sender_id }}:</strong>
      {{ message.content }}
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
</script>

<style scoped>
.message {
  margin-bottom: 10px;
}
</style>
