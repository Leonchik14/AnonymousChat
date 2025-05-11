<template>
  <LoginForm
    v-if="screen === 'login'"
    @login="handleLogin"
    @signup="screen = 'register'"
  />
  <RegisterForm
    v-else-if="screen === 'register'"
    @signup="handleSignUp"
    @back-to-login="screen = 'login'"
  />
  <EmailConfirmation
    v-else-if="screen === 'emailConfirm'"
    @confirm="confirmCode"
    @resend="resendEmail"
  />
  <ChatHistoryLayout
    v-else-if="screen === 'findPartner'"
    :chats="chatHistory"
    :selected-chat-id="selectedChatId"
    @find-partner="screen = 'chat'"
    @logout="logout"
    @select-chat="selectChat"
  />
  <ChatSearch
    v-else-if="screen === 'chat'"
    @cancel="cancelSearch"
    @partner-found="selectChatAndGoToHistory"
  />
</template>

<script setup>
import { ref, watch } from 'vue'
import LoginForm from './components/LoginForm.vue'
import RegisterForm from './components/RegisterForm.vue'
import EmailConfirmation from './components/EmailConfirmation.vue'
import ChatHistoryLayout from './components/ChatHistoryLayout.vue'
import ChatSearch from './components/ChatSearch.vue'
import { getChatListUrl, getChatHistoryUrl, getWsChatUrl } from './config/api'

const screen = ref('login')
const chatHistory = ref([])
const selectedChatId = ref(null)

watch(screen, async (newScreen) => {
  if (newScreen === 'findPartner') {
    await loadChats()
  }
})

async function loadChats() {
  const accessToken = localStorage.getItem('accessToken')
  const response = await fetch(getChatListUrl(), {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
    },
  })
  if (response.ok) {
    const chats = await response.json()
    chatHistory.value = chats.map(chat => {
      const myId = Number(localStorage.getItem('userId'))
      const partnerId = chat.User1ID === myId ? chat.User2ID : chat.User1ID
      const partnerLabel = partnerId ? `Partner #${String(partnerId).slice(-4)}` : `Partner`
      return {
        id: chat.ID,
        name: `Chat with ${partnerLabel}`,
        date: '',
        preview: '',
        messages: [],
        ...chat
      }
    })
    // Загружаем историю для всех чатов
    await Promise.all(chatHistory.value.map(chat => loadChatHistory(chat.id)))
    // Открываем первый чат
    if (chatHistory.value.length > 0) {
      selectedChatId.value = chatHistory.value[0].id
    }
  } else {
    chatHistory.value = []
  }
}

function handleLogin({ email, password }) {
  screen.value = 'findPartner'
}
function handleSignUp({ email, password, confirmPassword }) {
  screen.value = 'emailConfirm'
}
function confirmCode(code) {
  screen.value = 'login'
}
function resendEmail() {
  alert('Confirmation email resent!')
}
function logout() {
  screen.value = 'login'
}
function cancelSearch() {
  screen.value = 'findPartner'
}
async function selectChat(id) {
  selectedChatId.value = id
  await loadChatHistory(id)
  const chat = chatHistory.value.find(c => c.id === id)
  if (chat) chat.hasNew = false
}

async function loadChatHistory(chatId) {
  const accessToken = localStorage.getItem('accessToken')
  const response = await fetch(`${getChatHistoryUrl()}/${chatId}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
    },
  })
  if (response.ok) {
    const messages = await response.json()
    const chat = chatHistory.value.find(c => c.id === chatId)
    if (chat) {
      chat.messages = messages.map(msg => ({
        fromMe: msg.sender_id === Number(localStorage.getItem('userId')),
        text: msg.content,
        time: new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        ...msg
      }))
      if (chat.messages.length > 0) {
        const lastMsg = chat.messages[chat.messages.length - 1]
        chat.preview = lastMsg.text
        chat.date = lastMsg.time
      }
    }
  }
}

function handleIncomingMessage(msg) {
  const chat = chatHistory.value.find(c => c.id === msg.chat_id)
  if (chat) {
    chat.messages.push({
      fromMe: msg.sender_id === Number(localStorage.getItem('userId')),
      text: msg.content,
      time: new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      ...msg
    })
    chat.preview = msg.content
    if (selectedChatId.value !== chat.id) {
      chat.hasNew = true
    }
  }
}

function selectChatAndGoToHistory(chatId) {
  selectedChatId.value = chatId
  screen.value = 'findPartner'
}
</script>

<style>
body, html, #app {
  min-height: 100vh;
  min-width: 100vw;
  background: #181d29;
  margin: 0;
  padding: 0;
}
</style>