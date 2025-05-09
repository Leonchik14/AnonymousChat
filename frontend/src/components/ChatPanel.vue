<template>
  <div class="chat-panel" v-if="chat">
    <div class="chat-header">
      <span class="chat-title">{{ chat.name }}</span>
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
      <input class="chat-input" placeholder="Write a message..." />
      <button class="send-btn">Send</button>
    </div>
  </div>
</template>
<script setup>
const props = defineProps(['chat'])
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
