<template>
  <div class="searching-bg">
    <div class="searching-card">
      <div class="spinner"></div>
      <h1 class="searching-label">Searching for partner…</h1>
      <div class="timer">{{ formattedTime }}</div>
      <div class="searching-subtitle">Please wait while we find someone for you</div>
      <button class="cancel-btn" @click="$emit('cancel')">Cancel</button>
      <div v-if="showMatchAlert" class="alert-overlay">
        <div class="alert-box">
          <button class="close-alert" @click="showMatchAlert = false" aria-label="Close">&times;</button>
          <div style="margin-bottom: 1rem;">Собеседник найден!</div>
          <button class="go-to-chat-btn" @click="goToChat">Перейти к чату</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { getMatchmakingUrl } from '../config/api'
const timer = ref(0)
let interval = null

const showMatchAlert = ref(false)
const matchedChatId = ref(null)

onMounted(() => {
  timer.value = 0
  interval = setInterval(() => timer.value++, 1000)
  startLongPolling()
})

function startLongPolling() {
  const accessToken = localStorage.getItem('accessToken')
  const userId = localStorage.getItem('userId')
  fetch(getMatchmakingUrl(), {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'X-User-ID': userId,
      'Content-Type': 'application/json',
    },
  })
    .then(async (response) => {
      if (!response.ok) throw new Error('Ошибка поиска собеседника')
      const data = await response.json()
      if (data.event === 'match') {
        matchedChatId.value = data.data
        showMatchAlert.value = true
      }
    })
    .catch((err) => {
      // обработка ошибок
    })
}

function goToChat() {
  showMatchAlert.value = false
  if (matchedChatId.value) {
    // Эмитим событие для App.vue, чтобы открыть нужный чат
    // Можно назвать событие partner-found или chat-found
    // и передать chatId
    // Например:
    // $emit('partner-found', matchedChatId.value)
    // Но в <script setup> используем defineEmits:
    emit('partner-found', matchedChatId.value)
  }
}

const emit = defineEmits(['cancel', 'partner-found'])

onUnmounted(() => {
  clearInterval(interval)
})

const formattedTime = computed(() => {
  const min = String(Math.floor(timer.value / 60)).padStart(2, '0')
  const sec = String(timer.value % 60).padStart(2, '0')
  return `${min}:${sec}`
})
</script>

<style scoped>
.searching-bg {
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
.searching-card {
  background: #181d29ee;
  border-radius: 18px;
  box-shadow: 0 8px 32px #000a;
  padding: 3rem 2.5rem 2.5rem 2.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 350px;
  min-height: 350px;
  position: relative;
  animation: fadein 0.8s;
}
@keyframes fadein {
  from { transform: translateY(40px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
.spinner {
  width: 54px;
  height: 54px;
  border: 6px solid #23283a;
  border-top: 6px solid #7f4ad6;
  border-right: 6px solid #4a90e2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 2rem;
  box-shadow: 0 2px 12px #0003;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.searching-label {
  font-size: 2rem;
  color: #b6d6ff;
  font-weight: 700;
  margin-bottom: 0.5rem;
  text-align: center;
  letter-spacing: 1px;
}
.timer {
  font-size: 1.3rem;
  color: #7fa7d6;
  margin-bottom: 1rem;
  font-family: 'Fira Mono', 'Consolas', monospace;
  letter-spacing: 1px;
}
.searching-subtitle {
  color: #7fa7d6;
  font-size: 1.05rem;
  margin-bottom: 2rem;
  text-align: center;
  opacity: 0.95;
}
.cancel-btn {
  background: linear-gradient(90deg, #ff6a6a 0%, #ff3a3a 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.9rem 2.5rem;
  font-size: 1.1rem;
  cursor: pointer;
  margin-top: 1rem;
  transition: background 0.2s, transform 0.1s;
  font-weight: 600;
  box-shadow: 0 2px 12px #0003;
}
.cancel-btn:hover {
  background: linear-gradient(90deg, #ff3a3a 0%, #ff6a6a 100%);
  transform: translateY(-2px) scale(1.03);
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
  animation: bgfade 0.3s;
}
.alert-box {
  background: #181d29ee;
  color: #b6d6ff;
  padding: 2.5rem 2.2rem;
  border-radius: 18px;
  box-shadow: 0 8px 32px #000a;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.2rem;
  min-width: 260px;
  min-height: 120px;
  position: relative;
  text-align: center;
  justify-content: center;
  animation: fadein 0.5s;
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
  transition: color 0.2s;
}
.close-alert:hover {
  color: #ff6b6b;
}
.go-to-chat-btn {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.2rem;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 1rem;
  box-shadow: 0 2px 12px #0003;
  transition: all 0.2s;
}
.go-to-chat-btn:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.03);
}
</style>
