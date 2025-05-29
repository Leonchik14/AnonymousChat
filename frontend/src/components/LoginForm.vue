<template>
  <div class="login-form-bg">
    <div class="login-form animated">
      <h2>Welcome Back</h2>
      <div class="input-group">
        <span class="input-icon">@</span>
        <input v-model="email" placeholder="Email address" type="email" />
      </div>
      <div class="input-group">
        <span class="input-icon">&#128274;</span>
        <input v-model="password" placeholder="Password" type="password" />
      </div>
      <div v-if="error" class="error-message">{{ error }}</div>
      <button class="login-btn" @click="handleLogin" :disabled="isLoading">
        {{ isLoading ? 'Logging in...' : 'Log In' }}
      </button>
      <button class="signup-btn" @click="$emit('signup')">Sign Up</button>
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
  background: linear-gradient(135deg, #10131a 0%, #181d29 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: bgfade 1.2s;
}
@keyframes bgfade {
  from { filter: blur(8px) opacity(0.5);}
  to { filter: blur(0) opacity(1);}
}
.login-form {
  background: #181d29ee;
  color: #cbe6ff;
  padding: 2.5rem 2.2rem 2rem 2.2rem;
  border-radius: 18px;
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  box-shadow: 0 8px 32px #000a;
  animation: fadein 0.8s;
}
@keyframes fadein {
  from { transform: translateY(40px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
h2 {
  color: #b6d6ff;
  margin-bottom: 1.2rem;
  text-align: center;
  font-weight: 700;
  letter-spacing: 1px;
}
.input-group {
  display: flex;
  align-items: center;
  background: #10131a;
  border: 1.5px solid #23283a;
  border-radius: 8px;
  padding: 0.5rem 0.8rem;
  margin-bottom: 0.2rem;
  transition: border 0.2s;
}
.input-group:focus-within {
  border-color: #7fa7d6;
}
.input-icon {
  color: #7fa7d6;
  font-size: 1.2rem;
  margin-right: 0.7rem;
}
input {
  background: transparent;
  border: none;
  outline: none;
  color: #cbe6ff;
  font-size: 1rem;
  flex: 1;
  padding: 0.3rem 0;
}
.login-btn {
  background: linear-gradient(90deg, #23283a 0%, #4a5370 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 0.5rem;
  box-shadow: 0 2px 12px #0003;
  transition: background 0.2s, transform 0.1s;
}
.login-btn:hover {
  background: linear-gradient(90deg, #4a5370 0%, #23283a 100%);
  transform: translateY(-2px) scale(1.03);
}
.signup-btn {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 0.2rem;
  box-shadow: 0 2px 12px #0003;
  transition: background 0.2s, transform 0.1s;
}
.signup-btn:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.03);
}
.error-message {
  color: #ff6b6b;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
  text-align: center;
  font-weight: 500;
}
</style>
