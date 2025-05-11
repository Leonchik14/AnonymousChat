const REGISTER_URL = import.meta.env.REGISTER_URL || 'http://localhost/api/auth/register'
const EMAIL_VERIFICATION_URL = import.meta.env.EMAIL_VERIFICATION_URL || 'http://localhost/api/auth/send-verification'
const LOGIN_URL = import.meta.env.LOGIN_URL || 'http://localhost/api/auth/login'
const MATCHMAKING_URL = import.meta.env.MATCHMAKING_URL || 'http://localhost/api/matchmaking/start'
const CHAT_LIST_URL = import.meta.env.CHAT_LIST_URL || 'http://localhost/api/chat/all'
const CHAT_HISTORY_URL = import.meta.env.CHAT_HISTORY_URL || 'http://localhost/api/chat/history'
const WS_CHAT_URL = import.meta.env.WS_CHAT_URL || 'ws://localhost/ws/chat'
export const getRegisterUrl = () => {
  return REGISTER_URL
}

export const getEmailVerificationUrl = () => {
  return EMAIL_VERIFICATION_URL
}

export const getLoginUrl = () => {
  return LOGIN_URL
}

export const getMatchmakingUrl = () => {
  return MATCHMAKING_URL
}

export const getChatListUrl = () => {
  return CHAT_LIST_URL
}

export const getChatHistoryUrl = () => {
  return CHAT_HISTORY_URL
}

export const getWsChatUrl = () => {
  return WS_CHAT_URL
}

export const API_ENDPOINTS = {
  register: 'localhost/api/auth/register'
}

// Для отладки можно добавить логирование
if (import.meta.env.DEV) {
  console.log('Register endpoint URL:', REGISTER_URL)
} 