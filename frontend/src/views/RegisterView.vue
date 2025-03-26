<template>
  <div class="auth-container">
    <div class="auth-block">
      <h1>Registration</h1>
      <form @submit.prevent="register">
        <div class="form-group">
          <label for="username">Username</label>
          <input v-model="username" type="text" id="username" required />
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input v-model="email" type="email" id="email" required />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input v-model="password" type="password" id="password" required />
        </div>
        <button type="submit">Register</button>
      </form>
      <p class="redirect-text">
        Already have an account?
        <router-link to="/">Login</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { registerAPI } from '../api/api'

const router = useRouter()
const username = ref('')
const email = ref('')
const password = ref('')

const register = async () => {
  try {
    await registerAPI(username.value, email.value, password.value)
    alert('Registration successful. Please log in.')
    router.push('/')
  } catch (error) {
    alert(error.message)
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: white;
}
.auth-block {
  background-color: #f2f2f2;
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
  background-color: #fff9e6;
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
