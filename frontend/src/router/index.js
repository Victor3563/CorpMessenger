import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ChatView from '../views/ChatView.vue'

const routes = [
  { 
    path: '/', 
    name: 'Login', 
    component: LoginView 
  },
  { path: '/register', 
    name: 'Register', 
    component: RegisterView 
  },
  { path: '/chat', 
    name: 'Chat', 
    component: ChatView 
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
