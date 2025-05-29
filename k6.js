import ws from 'k6/ws';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '10s', target: 2000 },  // 2000 пользователей в секунду для теста
    ],
    thresholds: {
        'ws_connecting': ['p(95)<5000'],
        'ws_sessions': ['count>10000'],
        'ws_msgs_sent': ['count>20000'],
    },
};

export default function () {
    const chatId = Math.floor(Math.random() * 15000) + 1;
    const timestamp = new Date().toISOString();

    const url = `ws://localhost/ws/chat/${chatId}`;

    const res = ws.connect(url, {}, function (socket) {
        socket.on('open', function () {
            console.log(`WebSocket opened for chatId ${chatId}`);
            const payload = JSON.stringify({
                senderId: userId,
                content: `Test message ${Math.random().toString(36).substring(7)}`,
                timestamp: timestamp,
            });
            socket.send(payload);
        });

        socket.on('message', function (data) {
            const msg = JSON.parse(data);
            console.log(`Message received for chatId ${chatId}: ${JSON.stringify(msg)}`);
            check(msg, {
                'message received': () => msg.content && msg.sender_id,
            });
            socket.close();
        });

        socket.on('error', function (e) {
            console.log(`WebSocket error for chatId ${chatId}: ${e.error()}`);
        });

        socket.on('close', function () {
            console.log(`WebSocket closed for chatId ${chatId}`);
        });
    });

    check(res, {
        'webSocket connected': (r) => r && r.status === 101,
    });

    sleep(1);
}