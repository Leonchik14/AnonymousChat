<template>
  <div class="email-confirm-bg">
    <div class="email-confirm-card">
      <h2>Email Confirmation</h2>
      <p>
        A confirmation email has been sent to your email address.<br>
        Please enter the code from the email below.
      </p>
      <input
        v-model="confirmationCode"
        placeholder="Enter confirmation code"
        class="code-input"
      />
      <div class="button-row">
        <button @click="$emit('confirm', confirmationCode)">Confirm</button>
        <button class="secondary" @click="$emit('resend')">Resend Email</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const confirmationCode = ref('')
</script>

<style scoped>
.email-confirm-bg {
  min-height: 100vh;
  min-width: 100vw;
  background: #181d29;
  display: flex;
  align-items: center;
  justify-content: center;
}
.email-confirm-card {
  background: #23283a;
  color: #cbe6ff;
  padding: 2.5rem 2rem;
  border-radius: 12px;
  width: 340px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  box-shadow: 0 2px 24px #000a;
  align-items: center;
}
h2 {
  color: #b6d6ff;
  margin-bottom: 0.5rem;
  text-align: center;
}
p {
  color: #7fa7d6;
  font-size: 1rem;
  text-align: center;
  margin-bottom: 1rem;
}
.code-input {
  background: #10131a;
  border: 1px solid #2a3142;
  border-radius: 6px;
  padding: 0.5rem;
  color: #cbe6ff;
  width: 100%;
  margin-bottom: 1rem;
  font-size: 1.1rem;
  text-align: center;
  letter-spacing: 2px;
}
.button-row {
  display: flex;
  gap: 1rem;
  width: 100%;
  justify-content: center;
}
button {
  background: #2a3142;
  color: #b6d6ff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.2s;
}
button.secondary {
  background: #23283a;
  color: #7fa7d6;
}
button:hover {
  background: #3a425a;
}
</style>
