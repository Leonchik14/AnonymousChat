<template>
  <div class="chat-history-layout">
    <div class="sidebar">
      <h2>Chat History</h2>
      <ul>
        <li
          v-for="chat in chats"
          :key="chat.id"
          :class="{ active: chat.id === selectedChatId }"
          @click="$emit('select-chat', chat.id)"
        >
          <div class="chat-date">{{ chat.date }}</div>
          <div class="chat-name">{{ chat.name }}</div>
          <div class="chat-preview">
            {{ chat.messages.length ? chat.messages[chat.messages.length - 1].text : '' }}
          </div>
        </li>
      </ul>
      <button @click="$emit('find-partner')" class="find-partner-btn">Find Partner</button>
      <button @click="$emit('logout')" class="logout-btn">Logout</button>
    </div>
    <ChatPanel :chat="selectedChat" />
  </div>
</template>
<script setup>
import ChatPanel from './ChatPanel.vue'
import { computed } from 'vue'
const props = defineProps(['chats', 'selectedChatId'])
const selectedChat = computed(() =>
  props.chats.find(chat => chat.id === props.selectedChatId)
)
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
</style>
