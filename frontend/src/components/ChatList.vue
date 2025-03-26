<template>
  <div>
    <ul>
      <li v-for="chat in chats" :key="chat.id" class="chat-item">
        <span @click="select(chat)">{{ chat.name }}</span>
        <button class="info-btn" @click.stop="openChatInfo(chat.id)">⋮</button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
const props = defineProps({
  chats: { type: Array, required: true }
})
const emit = defineEmits(['selectChat', 'openChatInfo'])
const select = (chat) => {
  emit('selectChat', chat)
}
const openChatInfo = (chat) => {
  emit('openChatInfo', chat)
}
</script>

<style scoped>
.chat-item {
  padding: 8px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}
.chat-item:hover {
  background-color: #f5f5f5;
}
.info-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
}
</style>
