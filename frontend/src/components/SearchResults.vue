<template>
  <div class="search-results">
    <ul>
      <li v-for="user in results" :key="user.id" @click="select(user)">
        {{ user.username }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
const props = defineProps({
  results: { type: Array, required: true }
})
const emit = defineEmits(['select'])
const select = (user) => {
  emit('select', user)
}
</script>

<style scoped>
.search-results {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  width: 60%;
  background: white;
  border: 1px solid #ccc;
  max-height: 200px;
  overflow-y: auto;
  z-index: 5;
  color: black;
}
.search-results ul {
  list-style: none;
  padding: 0;
  margin: 0;
}
.search-results li {
  padding: 8px;
  cursor: pointer;
}
.search-results li:hover {
  background-color: #f0f0f0;
}
</style>
