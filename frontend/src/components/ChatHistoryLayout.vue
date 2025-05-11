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
  background: #181d29;
  border-radius: 12px;
  box-shadow: 0 2px 24px #000a;
  overflow: hidden;
  margin: 40px auto;
}

.sidebar {
  width: 300px;
  background: #181d29;
  border-right: 1px solid #233;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.sidebar h2 {
  color: #b6d6ff;
  text-align: left;
  padding: 1rem;
  margin: 0;
  font-size: 1.5rem;
  border-bottom: 1px solid #233;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.sidebar li {
  padding: 1rem;
  cursor: pointer;
  border-bottom: 1px solid #23283a;
  transition: background 0.2s;
}

.sidebar li.active,
.sidebar li:hover {
  background: #23283a;
}

.chat-date {
  color: #7fa7d6;
  font-size: 0.9rem;
}

.chat-name {
  color: #cbe6ff;
  font-weight: bold;
  font-size: 1.1rem;
}

.chat-preview {
  color: #7fa7d6;
}

.find-partner-btn, .logout-btn {
  margin: 1rem;
  padding: 0.7rem 1.2rem;
  background: #2a3142;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
  width: calc(100% - 2rem);
  display: block;
}

.find-partner-btn:hover, .logout-btn:hover {
  background: #3a425a;
}

.logout-btn {
  background: #23283a;
  color: #ff6a6a;
  margin-top: 0;
}

.logout-btn:hover {
  background: #ff6a6a;
  color: #fff;
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
  text-align: center;
}

.close-alert {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

.go-to-chat-btn {
  background: #2a3142;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 1rem;
  transition: background 0.2s;
}
.go-to-chat-btn:hover {
  background: #3a425a;
}
</style>
