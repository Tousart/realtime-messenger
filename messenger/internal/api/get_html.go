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
        
        function initWebSocket() {
            ws = new WebSocket('ws://' + window.location.host + '/messenger');
            
            ws.onopen = function() {
                console.log('WebSocket connected');
            };
            
            ws.onmessage = function(e) {
                console.log('Received:', e.data);
                addMessage('Сервер', e.data);
            };
            
            ws.onclose = function() {
                console.log('WebSocket disconnected');
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
                data: { chat_id: chatId }
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
                data: {
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

func (ap *API) getHomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(mainPageTemplate))
}
