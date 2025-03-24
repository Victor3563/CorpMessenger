<!-- src/components/MessageInput.vue -->
<template>
  <div class="input-container">
    <input v-model="message" @keyup.enter="handleSend" placeholder="Type your message..." />
    <button @click="handleSend">Send</button>
  </div>
</template>

<script setup>
import { ref, defineEmits } from 'vue'
const emit = defineEmits(['sendMessage'])
const message = ref('')

const handleSend = () => {
  if (message.value.trim() !== '') {
    emit('sendMessage', message.value.trim())
    message.value = ''
  }
}
</script>

<style scoped>
.input-container {
  display: flex;
}
input {
  flex: 1;
  padding: 8px;
}
button {
  padding: 8px 16px;
}
</style>
