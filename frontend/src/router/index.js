import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ChatView from '../views/ChatView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import ChatInfoView from '../views/ChatInfoView.vue'
import SearchView from '../views/SearchView.vue'

const routes = [
  { path: '/', name: 'Login', component: LoginView },
  { path: '/register', name: 'Register', component: RegisterView },
  { path: '/chat', name: 'Chat', component: ChatView },
  { path: '/profile', name: 'Profile', component: UserProfileView },
  { path: '/chatinfo/:id', name: 'ChatInfo', component: ChatInfoView, props: true },
  { path: '/search', name: 'Search', component: SearchView }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
