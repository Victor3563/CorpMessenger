<template>
  <div>
    <h2>Chats</h2>
    <ul>
      <li
        v-for="chat in chats"
        :key="chat.id"
        @click="select(chat)"
        class="chat-item"
      >
        <span>{{ chat.name }}</span>
        <span v-if="chat.unread_count > 0" class="unread">{{ chat.unread_count }}</span>
        <button class="info-btn" @click.stop="openChatInfo(chat.id)">⋮</button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  chats: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['selectChat', 'openChatInfo'])

const select = (chat) => {
  emit('selectChat', chat)
}

const openChatInfo = (chatId) => {
  emit('openChatInfo', chatId)
}
</script>

<style scoped>
.chat-item {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  align-items: center;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}
.chat-item:hover {
  background-color: #f0f0f0;
}
.unread {
  background-color: red;
  color: white;
  border-radius: 50%;
  padding: 2px 8px;
  font-size: 0.8rem;
  margin-left: 8px;
}
.info-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
}
</style>
