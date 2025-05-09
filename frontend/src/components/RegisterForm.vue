<template>
  <div class="register-bg">
    <div class="register-form">
      <template v-if="!showVerify">
        <h2>Sign Up</h2>
        <input v-model="email" placeholder="Email address" type="email" />
        <input v-model="password" placeholder="Password" type="password" />
        <input v-model="confirmPassword" placeholder="Confirm password" type="password" />
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
        <button @click="handleSignUp" :disabled="isLoading">
          {{ isLoading ? 'Signing up...' : 'Sign Up' }}
        </button>
        <button class="back-login-btn" @click="emit('back-to-login')" type="button">Back to Login</button>
        <div v-if="showSuccessAlert" class="alert-overlay">
          <div class="alert-box">
            <button class="close-alert" @click="showSuccessAlert = false" aria-label="Close">&times;</button>
            <div>{{ alertMessage }}</div>
            <button v-if="!emailSent" class="verify-btn" @click="sendEmailVerification">Verify email</button>
          </div>
        </div>
      </template>
      <template v-else>
        <h2>Email Verification</h2>
        <div class="verify-instruction">Введите код подтверждения, который был отправлен на вашу почту</div>
        <input v-model="verificationCode" placeholder="Verification code" type="text" />
        <button class="verify-btn" @click="handleVerifyCode">Submit code</button>
        <button class="verify-btn back-btn" @click="showVerify = false" type="button">Back</button>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getRegisterUrl, getEmailVerificationUrl } from '../config/api'

const emit = defineEmits(['signup-success', 'back-to-login'])

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const successMessage = ref('')
const confirmSuccessMessage = ref('')
const isLoading = ref(false)
const showVerify = ref(false)
const verificationCode = ref('')
const showSuccessAlert = ref(false)
const emailSent = ref(false)
const alertMessage = ref('')

let lastRegisteredEmail = ''

const handleSignUp = async () => {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  try {
    isLoading.value = true
    error.value = ''
    successMessage.value = ''
    
    const response = await fetch(getRegisterUrl(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || 'Registration failed')
    }

    // Показываем сообщение об успешной регистрации
    successMessage.value = data.message
    alertMessage.value = "Registration successful! Verify your email to receive access to the extended features or continue without verification."
    emailSent.value = false
    showSuccessAlert.value = true
    lastRegisteredEmail = email.value
    
    // Очищаем форму
    email.value = ''
    password.value = ''
    confirmPassword.value = ''

    emit('signup-success')
  } catch (err) {
    error.value = err.message
  } finally {
    isLoading.value = false
  }
}

const sendEmailVerification = async () => {
  try {
    const response = await fetch(getEmailVerificationUrl(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email: lastRegisteredEmail }),
    })
    if (!response.ok) {
      throw new Error('Failed to send verification email')
    }
    alertMessage.value = "Message sent. Check your email and follow the instructions."
    emailSent.value = true
  } catch (err) {
    alertMessage.value = err.message
  }
}

const handleVerifyCode = () => {
  // Здесь будет логика отправки кода подтверждения на бэкенд
  alert('Введённый код: ' + verificationCode.value)
}
</script>

<style scoped>
.register-bg {
  min-height: 100vh;
  min-width: 100vw;
  background: #181d29;
  display: flex;
  align-items: center;
  justify-content: center;
}
.register-form {
  background: #23283a;
  color: #cbe6ff;
  padding: 2.5rem 2rem;
  border-radius: 12px;
  width: 340px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  box-shadow: 0 2px 24px #000a;
}
h2 {
  color: #b6d6ff;
  margin-bottom: 1rem;
  text-align: center;
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
  margin-top: 0.5rem;
}
button:hover {
  background: #3a425a;
}
.error-message {
  color: #ff6b6b;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}
button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
.success-message {
  color: #4caf50;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}
.verify-btn {
  background: #3a425a;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 0.5rem;
  transition: background 0.2s;
}
.verify-btn:hover {
  background: #4a5370;
}
.verify-instruction {
  color: #b6d6ff;
  margin-bottom: 1rem;
  text-align: center;
}
.back-btn {
  background: #23283a;
  margin-left: 0.5rem;
}
.back-login-btn {
  background: #23283a;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 0.5rem;
  transition: background 0.2s;
}
.back-login-btn:hover {
  background: #3a425a;
}
.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.alert-box {
  background: #23283a;
  color: #b6d6ff;
  padding: 2rem 2.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 24px #000a;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  position: relative;
}
.close-alert {
  position: absolute;
  top: -1.2rem;
  right: -1.2rem;
  background: none;
  border: none;
  color: #b6d6ff;
  font-size: 2.2rem;
  cursor: pointer;
  line-height: 1;
  z-index: 10;
}
.close-alert:hover {
  color: #ff6b6b;
}
</style>
