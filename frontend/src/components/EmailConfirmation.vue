<template>
  <div class="email-confirm-bg">
    <div class="email-confirm-card">
      <div v-if="isLoading" class="confirm-message loading">Подтверждение...</div>
      <div v-else-if="confirmResult" :class="['confirm-message', confirmResult.success ? 'success' : 'error']">
        {{ confirmResult.message }}
      </div>
      <button v-if="!isLoading && confirmResult && confirmResult.success" class="go-login-btn" @click="$emit('go-to-login')">
        Войти
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const confirmationCode = ref('')
const confirmResult = ref(null)
const isLoading = ref(false)
const token = ref(null)

// Получаем токен из query-параметра
onMounted(() => {
  const params = new URLSearchParams(window.location.search)
  const t = params.get('token')
  if (t) {
    token.value = t
    confirmEmailByToken(t)
  }
})

async function confirmEmailByToken(token) {
  isLoading.value = true
  try {
    // Замените URL на ваш реальный endpoint подтверждения
    const response = await fetch(`/api/auth/verify?token=${encodeURIComponent(token)}`, {
      method: 'GET',
    })
    if (response.ok) {
      confirmResult.value = { success: true, message: 'Ваша почта успешно подтверждена! Теперь вы можете войти.' }
    } else {
      const data = await response.json().catch(() => ({}))
      confirmResult.value = { success: false, message: data.message || 'Ошибка подтверждения почты.' }
    }
  } catch (e) {
    confirmResult.value = { success: false, message: 'Ошибка сети. Попробуйте позже.' }
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.email-confirm-bg {
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
.email-confirm-card {
  background: #181d29ee;
  color: #cbe6ff;
  padding: 2.5rem 2.2rem 2rem 2.2rem;
  border-radius: 18px;
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  box-shadow: 0 8px 32px #000a;
  align-items: center;
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
p {
  color: #7fa7d6;
  font-size: 1.05rem;
  text-align: center;
  margin-bottom: 1rem;
  line-height: 1.5;
}
.code-input {
  background: #10131a;
  border: 1.5px solid #23283a;
  border-radius: 8px;
  padding: 0.7rem 0.8rem;
  color: #cbe6ff;
  width: 100%;
  margin-bottom: 1rem;
  font-size: 1.1rem;
  text-align: center;
  letter-spacing: 2px;
  transition: border 0.2s, box-shadow 0.2s;
}
.code-input:focus {
  outline: none;
  border-color: #7fa7d6;
  box-shadow: 0 0 0 2px rgba(127, 167, 214, 0.2);
}
.button-row {
  display: flex;
  gap: 1rem;
  width: 100%;
  justify-content: center;
}
button {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  box-shadow: 0 2px 12px #0003;
  transition: background 0.2s, transform 0.1s;
}
button.secondary {
  background: linear-gradient(90deg, #23283a 0%, #4a5370 100%);
  color: #7fa7d6;
}
button:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.03);
}
button.secondary:hover {
  background: linear-gradient(90deg, #4a5370 0%, #23283a 100%);
  color: #fff;
}
.confirm-message {
  font-size: 1.08rem;
  font-weight: 600;
  text-align: center;
  margin: 1.2rem 0 0.5rem 0;
  padding: 0.7rem 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 12px #0002;
}
.confirm-message.success {
  color: #6ee7b7;
  background: #1a3a2a33;
  text-shadow: 0 1px 8px #1a3a2a44;
}
.confirm-message.error {
  color: #ff6b6b;
  background: #3a232833;
  text-shadow: 0 1px 8px #3a232844;
}
.confirm-message.loading {
  color: #7fa7d6;
  background: #23283a33;
}
.go-login-btn {
  width: 100%;
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 2.2rem;
  box-shadow: 0 2px 12px #0003;
  transition: background 0.2s, color 0.2s, transform 0.1s;
  outline: none;
  display: block;
}
.go-login-btn:hover, .go-login-btn:focus {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  color: #fff;
  transform: translateY(-2px) scale(1.03);
}
</style>
