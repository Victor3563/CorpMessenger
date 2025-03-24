<template>
  <div class="login-container">
    <div class="auth-block">
      <h1>Login</h1>
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="username">Username</label>
          <input v-model="username" type="text" id="username" required />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input v-model="password" type="password" id="password" required />
        </div>
        <button type="submit">Login</button>
      </form>
      <p class="redirect-text">
        Don't have an account?
        <router-link to="/register">Register</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store'
import { loginAPI } from '../api/api'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const password = ref('')

const login = async () => {
  try {
    const data = await loginAPI(username.value, password.value)
    userStore.setUser(data)
    router.push('/chat')
  } catch (error) {
    alert(error.message)
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: white;
  padding: 0 20px;
}

.auth-block {
  background-color: #f2f2f2; /* светло-серый блок */
  padding: 40px;
  border-radius: 8px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  text-align: center;
}

.auth-block h1 {
  font-size: 2.5rem;
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
  text-align: left;
}

label {
  display: block;
  font-size: 1.2rem;
  margin-bottom: 5px;
  color: #555;
}

input {
  width: 100%;
  padding: 10px;
  font-size: 1.1rem;
  background-color: #fff9e6; /* мягкий желтый фон */
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 12px;
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

.redirect-text {
  margin-top: 15px;
  font-size: 0.9rem;
  color: black;
}

.redirect-text a {
  color: #007bff;
  text-decoration: none;
}

.redirect-text a:hover {
  text-decoration: underline;
}
</style>
