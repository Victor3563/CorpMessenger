<template>
    <div>
      <h2>Chats</h2>
      <ul>
        <li v-for="chat in chats" :key="chat.id" @click="select(chat)">
          {{ chat.name }}
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
  const emit = defineEmits(['selectChat'])
  
  const select = (chat) => {
    emit('selectChat', chat)
  }
  </script>
  
  <style scoped>
  li {
    cursor: pointer;
    padding: 8px;
    border-bottom: 1px solid #eee;
  }
  li:hover {
    background-color: #f5f5f5;
  }
  </style>
  