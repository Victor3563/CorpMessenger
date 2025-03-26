<template>
  <div class="profile-container">
    <div class="profile-block">
      <h1>User Profile</h1>
      <form v-if="isOwnProfile" @submit.prevent="updateProfile">
        <div class="form-group">
          <label for="username">Username</label>
          <input v-model="username" type="text" id="username" required />
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input v-model="email" type="email" id="email" required />
        </div>
        <!-- Поле password показывается только для редактирования своего профиля -->
        <div class="form-group">
          <label for="password">New Password</label>
          <input v-model="password" type="password" id="password" placeholder="Leave blank if unchanged" />
        </div>
        <button type="submit">Update Profile</button>
      </form>
      <div v-else class="readonly-profile">
        <div class="form-group">
          <label>Username</label>
          <p>{{ username }}</p>
        </div>
        <div class="form-group">
          <label>Email</label>
          <p>{{ email }}</p>
        </div>
      </div>
      <button @click="goBack">Back to Chat</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../store'
import { useRouter, useRoute } from 'vue-router'
import { getUserByIdAPI, updateUserAPI } from '../api/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const username = ref('')
const email = ref('')
const password = ref('')

const ownUserId = userStore.user ? userStore.user.id : null
const profileUserId = route.query.user_id ? parseInt(route.query.user_id) : ownUserId

const isOwnProfile = computed(() => profileUserId === ownUserId)

onMounted(async () => {
  const userData = await getUserByIdAPI(profileUserId)
  username.value = userData.username
  email.value = userData.email
})

const updateProfile = async () => {
  try {
    await updateUserAPI(ownUserId, username.value, email.value, password.value);
    alert('Profile updated successfully');
  } catch (error) {
    alert(error.message);
  }
}

const goBack = () => {
  router.push('/chat')
}
</script>

<style scoped>
.profile-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: white;
}
.profile-block {
  background-color: #f2f2f2;
  padding: 40px;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  text-align: center;
}
.profile-block h1 {
  margin-bottom: 20px;
  font-size: 2.5rem;
  color: #333;
}
.form-group {
  margin-bottom: 15px;
  text-align: left;
}
label {
  display: block;
  margin-bottom: 5px;
  font-size: 1.2rem;
  color: #555;
}
input {
  width: 100%;
  padding: 10px;
  font-size: 1.1rem;
  background-color: #fff9e6;
  border: 1px solid #ccc;
  border-radius: 4px;
}
button {
  padding: 12px 20px;
  font-size: 1.2rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}
button:hover {
  background-color: #0056b3;
}
</style>
