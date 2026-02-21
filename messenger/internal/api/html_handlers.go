package api

import "net/http"

const mainPageTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Messenger</title>
    <style>
        #home, #chat { padding: 20px; }
        button { margin: 5px; padding: 10px; }
        #messages { border: 1px solid #ccc; height: 300px; overflow-y: auto; padding: 10px; margin: 10px 0; }
        #chatInput { width: 300px; padding: 5px; }
        .message { margin: 5px 0; }
    </style>
</head>
<body>
    <div id="home">
        <h2>Выберите чат:</h2>
        <button onclick="selectChat(1)">Чат 1</button>
        <button onclick="selectChat(2)">Чат 2</button>
        <button onclick="selectChat(3)">Чат 3</button>
    </div>
    
    <div id="chat" style="display:none;">
        <h2 id="chatTitle">Чат</h2>
        <div id="messages"></div>
        <div>
            <input id="chatInput" type="text" placeholder="Введите сообщение...">
            <button onclick="sendMessage()">Отправить</button>
            <button onclick="showHome()">← Назад</button>
        </div>
    </div>
    
    <script>
        let ws = null;
        let currentChatId = null;
        let wasConnected = false; // флаг успешного открытия соединения
        
        function initWebSocket() {
            ws = new WebSocket('ws://' + window.location.host + '/messenger');
            
            ws.onopen = function() {
                console.log('WebSocket connected');
                wasConnected = true;
            };
            
            ws.onmessage = function(e) {
                console.log('Received:', e.data);
                addMessage('Сервер', e.data);
            };
            
            ws.onclose = function(event) {
                console.log('WebSocket disconnected', event);
                // Если соединение так и не открылось (нет сессии), редирект на логин
                if (!wasConnected) {
                    window.location.href = '/auth/login';
                }
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        function selectChat(chatId) {
            currentChatId = chatId;
            
            document.getElementById('home').style.display = 'none';
            document.getElementById('chat').style.display = 'block';
            document.getElementById('chatTitle').textContent = 'Чат ' + chatId;
            document.getElementById('messages').innerHTML = '';
            
            ws.send(JSON.stringify({
                method: "join",
                payload: { chat_id: chatId }
            }));
            
            document.getElementById('chatInput').focus();
        }
        
        function showHome() {
            currentChatId = null;
            document.getElementById('chat').style.display = 'none';
            document.getElementById('home').style.display = 'block';
        }
        
        function sendMessage() {
            if (!currentChatId) return;
            
            const input = document.getElementById('chatInput');
            const text = input.value.trim();
            
            if (!text) return;
            
            ws.send(JSON.stringify({
                method: "send",
                payload: {
                    chat_id: currentChatId,
                    text: text
                }
            }));
            
            addMessage('Вы', text);
            input.value = '';
            input.focus();
        }
        
        function addMessage(sender, text) {
            const messagesDiv = document.getElementById('messages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message';
            messageDiv.innerHTML = '<strong>' + sender + ':</strong> ' + text;
            messagesDiv.appendChild(messageDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }
        
        document.addEventListener('DOMContentLoaded', function() {
            initWebSocket();
            
            document.getElementById('chatInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') sendMessage();
            });
        });
    </script>
</body>
</html>`

const loginPageTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Вход - Messenger</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background-color: #f5f5f5;
        }
        .login-container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            width: 300px;
        }
        h2 {
            margin-top: 0;
            color: #333;
        }
        input {
            width: 100%;
            padding: 10px;
            margin: 10px 0;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
        }
        button {
            width: 100%;
            padding: 10px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #0056b3;
        }
        .error {
            color: red;
            font-size: 14px;
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Вход в Messenger</h2>
        <form method="POST" action="/auth/login">
            <input type="text" name="user_name" placeholder="Имя пользователя" required>
            <input type="password" name="password" placeholder="Пароль" required>
            <button type="submit">Войти</button>
        </form>
    </div>
</body>
</html>`

func (ap *API) getHomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(mainPageTemplate))
}

func (ap *API) getLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(loginPageTemplate))
}
