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
  animation: fadein 0.8s;
}

@keyframes fadein {
  from { transform: translateY(40px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.chat-header {
  padding: 1.2rem;
  border-bottom: 1px solid #23283a;
  color: #b6d6ff;
  font-size: 1.3rem;
  font-weight: 700;
  letter-spacing: 1px;
  background: #181d29;
}

.chat-messages {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: linear-gradient(135deg, #10131a 0%, #181d29 100%);
}

.chat-message {
  max-width: 70%;
  padding: 0.8rem 1.2rem;
  border-radius: 12px;
  font-size: 1.1rem;
  position: relative;
  display: flex;
  flex-direction: column;
  animation: messageAppear 0.3s;
  box-shadow: 0 2px 12px #0003;
}

@keyframes messageAppear {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.from-me {
  align-self: flex-end;
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
}

.from-them {
  align-self: flex-start;
  background: #23283a;
  color: #cbe6ff;
}

.msg-text {
  margin-bottom: 0.3rem;
  line-height: 1.4;
}

.msg-time {
  font-size: 0.85rem;
  color: #7fa7d6;
  align-self: flex-end;
  opacity: 0.8;
}

.chat-input-row {
  display: flex;
  padding: 1.2rem;
  border-top: 1px solid #23283a;
  background: #181d29;
  gap: 1rem;
}

.chat-input {
  flex: 1;
  padding: 0.8rem 1.2rem;
  border-radius: 8px;
  border: 1.5px solid #23283a;
  background: #10131a;
  color: #cbe6ff;
  font-size: 1rem;
  transition: all 0.2s;
}

.chat-input:focus {
  outline: none;
  border-color: #7fa7d6;
  box-shadow: 0 0 0 2px rgba(127, 167, 214, 0.2);
}

.send-btn {
  background: linear-gradient(90deg, #7f4ad6 0%, #4a90e2 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.8rem 1.5rem;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 12px #0003;
}

.send-btn:hover {
  background: linear-gradient(90deg, #4a90e2 0%, #7f4ad6 100%);
  transform: translateY(-2px) scale(1.03);
}

.send-btn:active {
  transform: translateY(0) scale(0.98);
}
</style>
