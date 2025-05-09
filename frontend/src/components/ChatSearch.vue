<template>
  <div class="searching-bg">
    <div class="searching-card">
      <div class="spinner"></div>
      <h1 class="searching-label">Searching for partner…</h1>
      <div class="timer">{{ formattedTime }}</div>
      <div class="searching-subtitle">Please wait while we find someone for you</div>
      <button class="cancel-btn" @click="$emit('cancel')">Cancel</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
const timer = ref(0)
let interval = null
onMounted(() => {
  timer.value = 0
  interval = setInterval(() => timer.value++, 1000)
})
onUnmounted(() => clearInterval(interval))

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
  background: #181d29;
  display: flex;
  align-items: center;
  justify-content: center;
}
.searching-card {
  background: #23283a;
  border-radius: 18px;
  box-shadow: 0 2px 24px #000a;
  padding: 3rem 2.5rem 2.5rem 2.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 350px;
  min-height: 350px;
  position: relative;
}
.spinner {
  width: 48px;
  height: 48px;
  border: 5px solid #2a3142;
  border-top: 5px solid #b6d6ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 2rem;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.searching-label {
  font-size: 2rem;
  color: #b6d6ff;
  font-weight: bold;
  margin-bottom: 0.5rem;
  text-align: center;
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
}
.cancel-btn {
  background: #ff6a6a;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.9rem 2.5rem;
  font-size: 1.1rem;
  cursor: pointer;
  margin-top: 1rem;
  transition: background 0.2s;
  font-weight: 500;
  box-shadow: 0 2px 12px #0003;
}
.cancel-btn:hover {
  background: #ff3a3a;
}
</style>
