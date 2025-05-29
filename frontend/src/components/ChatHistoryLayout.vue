<template>
  <div class="chat-history-layout">
    <div class="sidebar">
      <h2>Chat History</h2>
      <ul>
        <li
          v-for="chat in chats"
          :key="chat.id"
          :class="{ active: chat.id === selectedChatId, 'has-new': chat.hasNew }"
          @click="emit('select-chat', chat.id)"
        >
          <div class="chat-date">{{ chat.date }}</div>
          <div class="chat-name">{{ chat.name }}</div>
          <div class="chat-preview">
            {{ chat.preview }}
            <span v-if="chat.hasNew" class="new-indicator"></span>
          </div>
        </li>
      </ul>
      <button @click="handleFindPartner" class="find-partner-btn">Find Partner</button>
      <button @click="$emit('logout')" class="logout-btn">Logout</button>
    </div>
    <ChatPanel :chat="selectedChat" />
    <div v-if="showMatchAlert" class="alert-overlay">
      <div class="alert-box">
        <button class="close-alert" @click="showMatchAlert = false" aria-label="Close">&times;</button>
        <div>Собеседник найден!</div>
        <button class="go-to-chat-btn" @click="goToChat">Перейти к чату</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import ChatPanel from './ChatPanel.vue'
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { getMatchmakingUrl } from '../config/api'

const props = defineProps(['chats', 'selectedChatId'])
const emit = defineEmits(['find-partner', 'logout', 'select-chat'])
const selectedChat = computed(() =>
  props.chats.find(chat => chat.id === props.selectedChatId)
)

const timer = ref(0)
let interval = null
let abortController = null
const showMatchAlert = ref(false)
const matchedChatId = ref(null)

onMounted(() => {
  timer.value = 0
  interval = setInterval(() => timer.value++, 1000)
  startLongPolling()
})
onUnmounted(() => {
  clearInterval(interval)
  if (abortController) abortController.abort()
})

const formattedTime = computed(() => {
  const min = String(Math.floor(timer.value / 60)).padStart(2, '0')
  const sec = String(timer.value % 60).padStart(2, '0')
  return `${min}:${sec}`
})

function startLongPolling() {
  abortController = new AbortController()
  const userId = localStorage.getItem('userId')
  fetch(getMatchmakingUrl(), {
    method: 'GET',
    headers: {
      'X-User-ID': userId,
      'Content-Type': 'application/json',
    },
    signal: abortController.signal,
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
      if (err.name !== 'AbortError') {
        // Можно повторить запрос, если нужно
        // startLongPolling()
      }
    })
}

const handleFindPartner = async () => {
  emit('find-partner')
  try {
    const accessToken = localStorage.getItem('accessToken')
    const response = await fetch(getMatchmakingUrl(), {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
    })
    if (!response.ok) {
      throw new Error('Failed to start matchmaking')
    }
  } catch (err) {
    alert(err.message)
  }
}

function goToChat() {
  showMatchAlert.value = false
  if (matchedChatId.value) {
    emit('select-chat', matchedChatId.value)
    emit('find-partner')
  }
}
</script>
<style scoped>
.chat-history-layout {
  display: flex;
  width: 900px;
  height: 600px;
  background: #181d29ee;
  border-radius: 18px;
  box-shadow: 0 8px 32px #000a;
  overflow: hidden;
  margin: 40px auto;
  animation: fadein 0.8s;
}

@keyframes fadein {
  from { transform: translateY(40px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.sidebar {
  width: 300px;
  background: linear-gradient(135deg, #181d29 60%, #23283a 100%);
  border-right: 1px solid #23283a;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  box-shadow: 4px 0 24px #0003;
  position: relative;
  z-index: 2;
}

.sidebar h2 {
  color: #b6d6ff;
  text-align: left;
  padding: 1.5rem 1.2rem 1.2rem 1.2rem;
  margin: 0;
  font-size: 1.7rem;
  font-weight: 800;
  letter-spacing: 1.5px;
  border-bottom: 1px solid #23283a;
  text-shadow: 0 2px 8px #0004;
}

.sidebar ul {
  list-style: none;
  padding: 0 0.5rem;
  margin: 0;
  flex: 1;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #4a90e2 #181d29;
}
.sidebar ul::-webkit-scrollbar {
  width: 7px;
}
.sidebar ul::-webkit-scrollbar-thumb {
  background: #4a90e2;
  border-radius: 6px;
}
.sidebar ul::-webkit-scrollbar-track {
  background: #181d29;
}

.sidebar li {
  margin: 1rem 0 0.7rem 0;
  padding: 1.1rem 1.1rem 0.9rem 1.1rem;
  border-radius: 14px;
  background: rgba(36, 44, 70, 0.85);
  box-shadow: 0 2px 16px #0002, 0 1.5px 0 #23283a;
  cursor: pointer;
  border: none;
  transition: background 0.22s, box-shadow 0.22s, transform 0.18s;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  position: relative;
}

.sidebar li.active,
.sidebar li:hover {
  background: linear-gradient(90deg, #23283a 60%, #4a90e2 100%);
  box-shadow: 0 4px 24px #4a90e244, 0 1.5px 0 #23283a;
  transform: translateY(-2px) scale(1.03);
}

.sidebar li:not(:last-child)::after {
  content: '';
  display: block;
  height: 1px;
  background: #23283a;
  opacity: 0.5;
  margin: 0.7rem -1.1rem 0 -1.1rem;
}

.chat-date {
  color: #7fa7d6;
  font-size: 0.93rem;
  font-weight: 500;
  margin-bottom: 0.1rem;
  letter-spacing: 0.5px;
}

.chat-name {
  color: #cbe6ff;
  font-weight: 700;
  font-size: 1.13rem;
  margin: 0.1rem 0 0.1rem 0;
  letter-spacing: 0.2px;
}

.chat-preview {
  color: #7fa7d6;
  font-size: 0.98rem;
  font-weight: 400;
  opacity: 0.92;
  margin-top: 0.1rem;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.find-partner-btn, .logout-btn {
  margin: 1.2rem 1.2rem 0.5rem 1.2rem;
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  width: calc(100% - 2.4rem);
  display: block;
  box-shadow: 0 2px 12px #0003;
}

.find-partner-btn {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
}

.find-partner-btn:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.03);
}

.logout-btn {
  background: linear-gradient(90deg, #23283a 0%, #4a5370 100%);
  color: #ff6a6a;
  margin-top: 0.5rem;
}

.logout-btn:hover {
  background: linear-gradient(90deg, #4a5370 0%, #23283a 100%);
  color: #fff;
  transform: translateY(-2px) scale(1.03);
}

.has-new .chat-name, .has-new .chat-preview {
  font-weight: bold;
  color: #b6ffb6;
}

.new-indicator {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #4caf50;
  border-radius: 50%;
  margin-left: 6px;
  box-shadow: 0 0 8px #4caf50;
}

.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  animation: bgfade 0.3s;
}

@keyframes bgfade {
  from { opacity: 0; }
  to { opacity: 1; }
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
  position: relative;
  text-align: center;
  animation: fadein 0.5s;
}

.close-alert {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  color: #7fa7d6;
  font-size: 1.5rem;
  cursor: pointer;
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
