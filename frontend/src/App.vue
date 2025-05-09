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
  />
</template>

<script setup>
import { ref } from 'vue'
import LoginForm from './components/LoginForm.vue'
import RegisterForm from './components/RegisterForm.vue'
import EmailConfirmation from './components/EmailConfirmation.vue'
import ChatHistoryLayout from './components/ChatHistoryLayout.vue'
import ChatSearch from './components/ChatSearch.vue'

const screen = ref('login')
const chatHistory = ref([
  {
    id: 1,
    name: 'John Smith',
    date: '15 June',
    preview: 'Text message',
    messages: [
      { fromMe: false, text: 'Hi! How are you?', time: '10:35 am' },
      { fromMe: true, text: 'Hi! How are you?', time: '10:36 am' },
      { fromMe: false, text: 'Another message!', time: '10:37 am' }
    ]
  },
  {
    id: 2,
    name: 'John Smith',
    date: '10 June',
    preview: 'Partner',
    messages: [
      { fromMe: false, text: 'Hello again!', time: '11:00 am' }
    ]
  }
])
const selectedChatId = ref(chatHistory.value[0]?.id ?? null)

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
function selectChat(id) {
  selectedChatId.value = id
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