<template>
  <div class="login-form-bg">
    <div class="login-form">
      <h2>Login</h2>
      <input v-model="email" placeholder="Email address" type="email" />
      <input v-model="password" placeholder="Password" type="password" />
      <div v-if="error" class="error-message">{{ error }}</div>
      <button @click="handleLogin" :disabled="isLoading">{{ isLoading ? 'Logging in...' : 'Log In' }}</button>
      <button class="secondary" @click="$emit('signup')">Sign Up</button>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { getLoginUrl } from '../config/api'

const emit = defineEmits(['login', 'signup'])
const email = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  error.value = ''
  isLoading.value = true
  try {
    const response = await fetch(getLoginUrl(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email: email.value, password: password.value }),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.message || 'Login failed')
    }
    if (data.accessToken) {
      localStorage.setItem('accessToken', data.accessToken)
      localStorage.setItem('refreshToken', data.refreshToken)
      localStorage.setItem('userId', data.userId)
      emit('login', { email: email.value, password: password.value })
    } else {
      throw new Error('No access token received')
    }
  } catch (err) {
    error.value = err.message
  } finally {
    isLoading.value = false
  }
}
</script>
<style scoped>
.login-form-bg {
  min-height: 100vh;
  min-width: 100vw;
  background: #181d29;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-form {
  background: #23283a;
  color: #cbe6ff;
  padding: 2rem;
  border-radius: 12px;
  width: 340px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  box-shadow: 0 2px 24px #000a;
}
input {
  background: #10131a;
  border: 1px solid #2a3142;
  border-radius: 6px;
  padding: 0.5rem;
  color: #cbe6ff;
}
button {
  background: #2a3142;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.2s;
}
button.secondary {
  background: #23283a;
  color: #7fa7d6;
}
button:hover {
  background: #3a425a;
}
.error-message {
  color: #ff6b6b;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  text-align: center;
}
</style>
