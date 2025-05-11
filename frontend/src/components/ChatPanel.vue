<template>
  <div class="chat-panel" v-if="chat">
    <div class="chat-header">
      <span class="chat-title">
        {{ chat.name || 'Anonymous Chat' }}
      </span>
    </div>
    <div class="chat-messages">
      <div
        v-for="(msg, idx) in chat.messages"
        :key="idx"
        :class="['chat-message', msg.fromMe ? 'from-me' : 'from-them']"
      >
        <span class="msg-text">{{ msg.text }}</span>
        <span class="msg-time">{{ msg.time }}</span>
      </div>
    </div>
    <div class="chat-input-row">
      <input class="chat-input" v-model="newMessage" placeholder="Write a message..." @keyup.enter="sendMessage" />
      <button class="send-btn" @click="sendMessage">Send</button>
    </div>
  </div>
</template>
<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { getWsChatUrl } from '../config/api'

const props = defineProps(['chat'])
const ws = ref(null)
const newMessage = ref('')
const userId = Number(localStorage.getItem('userId'))

watch(() => props.chat?.id, (chatId) => {
  if (ws.value) ws.value.close()
  if (chatId) {
    ws.value = new WebSocket(getWsChatUrl() + '/' + chatId)
    ws.value.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      props.chat.messages.push({
        fromMe: msg.sender_id === userId,
        text: msg.content,
        time: new Date(msg.created_at).toLocaleTimeString(),
        ...msg
      })
    }
  }
})

function sendMessage() {
  if (ws.value && newMessage.value.trim()) {
    ws.value.send(JSON.stringify({
      senderId: userId,
      content: newMessage.value,
      timestamp: new Date().toISOString()
    }))
    newMessage.value = ''
  }
}

onUnmounted(() => {
  if (ws.value) ws.value.close()
})
</script>
<style scoped>
.chat-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #181d29;
  height: 100%;
}

.chat-header {
  padding: 1rem;
  border-bottom: 1px solid #233;
  color: #b6d6ff;
  font-size: 1.3rem;
  font-weight: bold;
}

.chat-messages {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.chat-message {
  max-width: 70%;
  padding: 0.7rem 1rem;
  border-radius: 12px;
  font-size: 1.1rem;
  position: relative;
  display: flex;
  flex-direction: column;
  background: #23283a;
  color: #cbe6ff;
}

.from-me {
  align-self: flex-end;
  background: #2a3142;
  color: #b6d6ff;
}

.from-them {
  align-self: flex-start;
  background: #23283a;
  color: #cbe6ff;
}

.msg-text {
  margin-bottom: 0.3rem;
}

.msg-time {
  font-size: 0.85rem;
  color: #7fa7d6;
  align-self: flex-end;
}

.chat-input-row {
  display: flex;
  padding: 1rem;
  border-top: 1px solid #233;
  background: #181d29;
}

.chat-input {
  flex: 1;
  padding: 0.7rem 1rem;
  border-radius: 8px;
  border: 1px solid #2a3142;
  background: #10131a;
  color: #cbe6ff;
  font-size: 1rem;
  margin-right: 1rem;
}

.send-btn {
  background: #2a3142;
  color: #b6d6ff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.send-btn:hover {
  background: #3a425a;
}
</style>
