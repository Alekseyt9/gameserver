document.querySelector('.game-tile').addEventListener('click', function() {
    document.getElementById('modal').style.display = 'flex';
    document.body.classList.add('modal-open');
    //document.getElementById('game-iframe').src = 'content/matching/index.html';
    document.getElementById('game-iframe').src = 'content/games/tictactoe/index.html';
});

document.getElementById('close-btn').addEventListener('click', function() {
    document.getElementById('modal').style.display = 'none';
    document.body.classList.remove('modal-open');
});

class WebSocketManager {
    constructor(url) {
        this.url = url;
        this.websocket = null;
        this.eventHandlers = {
            onopen: [],
            onmessage: [],
            onclose: []
        };
    }

    getConnection() {
        if (!this.websocket || this.websocket.readyState === WebSocket.CLOSED) {
            this.websocket = new WebSocket(this.url);

            this.websocket.onopen = (event) => {
                this.eventHandlers.onopen.forEach(handler => handler(event));
            };

            this.websocket.onmessage = (event) => {
                this.eventHandlers.onmessage.forEach(handler => handler(event));
            };

            this.websocket.onclose = (event) => {
                this.eventHandlers.onclose.forEach(handler => handler(event));
            };
        }
        return this.websocket;
    }

    subscribe(event, handler) {
        if (this.eventHandlers[event]) {
            this.eventHandlers[event].push(handler);
        } else {
            console.warn(`Event ${event} is not supported.`);
        }
    }
}

var wsUrlElement = document.getElementById("ws-url");
var wsUrl = wsUrlElement ? wsUrlElement.textContent : "ws://default-url";

window.ws = new WebSocketManager(wsUrl);
const wsInstance = window.ws.getConnection();

window.ws.subscribe('onopen', (event) => {
    console.log("Connection opened!");
    wsInstance.send("Hello, server!");
});

window.ws.subscribe('onmessage', (event) => {
    console.log("Received: " + event.data);
});

window.ws.subscribe('onclose', (event) => {
    console.log("Connection closed");
});