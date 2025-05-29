<template>
  <div class="register-bg">
    <div class="register-form animated">
      <template v-if="!showVerify">
        <h2>Sign Up</h2>
        <div class="input-group">
          <span class="input-icon">@</span>
          <input v-model="email" placeholder="Email address" type="email" />
        </div>
        <div class="input-group">
          <span class="input-icon">&#128274;</span>
          <input v-model="password" placeholder="Password" type="password" />
        </div>
        <div class="input-group">
          <span class="input-icon">&#128274;</span>
          <input v-model="confirmPassword" placeholder="Confirm password" type="password" />
        </div>
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
        <button class="signup-btn" @click="handleSignUp" :disabled="isLoading">
          {{ isLoading ? 'Signing up...' : 'Sign Up' }}
        </button>
        <button class="back-login-btn" @click="emit('back-to-login')" type="button">Back to Login</button>
        <div v-if="showSuccessAlert" class="alert-overlay">
          <div class="alert-box enhanced-alert">
            <button class="close-alert enhanced-close" @click="showSuccessAlert = false" aria-label="Close">&times;</button>
            <div class="alert-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect width="24" height="24" rx="12" fill="#23283a"/>
                <path d="M6 8l6 5 6-5" stroke="#7fa7d6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <rect x="6" y="8" width="12" height="8" rx="2" stroke="#7fa7d6" stroke-width="2"/>
              </svg>
            </div>
            <div class="alert-message">{{ alertMessage }}</div>
            <button v-if="!emailSent" class="verify-btn enhanced-verify-btn" @click="sendEmailVerification">Verify email</button>
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
.register-form {
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
.back-login-btn {
  width: 100%;
  background: linear-gradient(90deg, #23283a 0%, #4a90e2 100%);
  color: #b6d6ff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 1.1rem;
  box-shadow: 0 2px 12px #0003;
  transition: background 0.2s, color 0.2s, transform 0.1s;
  outline: none;
  display: block;
  margin-left: auto;
  margin-right: auto;
}
.back-login-btn:hover, .back-login-btn:focus {
  background: linear-gradient(90deg, #4a90e2 0%, #23283a 100%);
  color: #fff;
  transform: translateY(-2px) scale(1.03);
}
.error-message {
  color: #ff6b6b;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
  text-align: center;
  font-weight: 500;
}
.success-message {
  color: #6ee7b7;
  font-size: 1.08rem;
  margin-top: 0.5rem;
  font-weight: 600;
  text-align: center;
  text-shadow: 0 1px 8px #1a3a2a44;
  letter-spacing: 0.2px;
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
.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(16, 19, 26, 0.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadein 0.4s;
}
.alert-box {
  background: #23283ad9;
  color: #b6d6ff;
  padding: 2.5rem 2.7rem 2.2rem 2.7rem;
  border-radius: 18px;
  box-shadow: 0 8px 32px #000a, 0 2px 24px #4a90e255;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.2rem;
  position: relative;
  min-width: 340px;
  max-width: 90vw;
  animation: popin 0.5s cubic-bezier(.22,1,.36,1);
}
@keyframes popin {
  from { transform: scale(0.92) translateY(30px); opacity: 0; }
  to { transform: scale(1) translateY(0); opacity: 1; }
}
.alert-icon {
  margin-bottom: 0.2rem;
  animation: iconpop 0.7s cubic-bezier(.22,1,.36,1);
}
@keyframes iconpop {
  from { transform: scale(0.7) rotate(-10deg); opacity: 0; }
  to { transform: scale(1) rotate(0); opacity: 1; }
}
.alert-message {
  color: #cbe6ff;
  font-size: 1.08rem;
  text-align: center;
  margin-bottom: 0.2rem;
  line-height: 1.5;
  font-weight: 500;
}
.enhanced-verify-btn {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  font-weight: 600;
  font-size: 1.08rem;
  padding: 0.7rem 1.6rem;
  border-radius: 8px;
  box-shadow: 0 2px 12px #0003;
  margin-top: 0.5rem;
  border: none;
  transition: background 0.2s, transform 0.1s;
}
.enhanced-verify-btn:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.04);
}
.enhanced-close {
  position: absolute;
  top: -1.3rem;
  right: -1.3rem;
  background: #23283a;
  border: 2px solid #7fa7d6;
  color: #b6d6ff;
  font-size: 2.3rem;
  border-radius: 50%;
  width: 2.6rem;
  height: 2.6rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  line-height: 1;
  z-index: 10;
  box-shadow: 0 2px 12px #0005;
  transition: background 0.2s, color 0.2s, border 0.2s;
}
.enhanced-close:hover {
  background: #ff6b6b;
  color: #fff;
  border: 2px solid #ff6b6b;
}
</style>
