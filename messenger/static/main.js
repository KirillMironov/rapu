let socket;

function connect() {
    if (socket != null) {
        socket.close();
        document.getElementById('messages_area').value = '';
    }

    socket = new WebSocket('ws://localhost:7004/connect' +
        '?userId='
        + document.getElementById('userId').value
        + '&toUserId='
        + document.getElementById('toUserId').value);

    console.log('Attempting Connection...');

    socket.onopen = () => {
        console.log('Successfully Connected');
    };

    socket.onmessage = message => {
        let json = JSON.parse(message.data);

        if (json.constructor === Array) {
            json.forEach(msg => {
                document.getElementById('messages_area').value += `${msg.from}: ${String(msg.text)}\r\n`;
            });
        } else {
            document.getElementById('messages_area').value += `${json.from}: ${String(json.text)}\r\n`;
        }
    };

    socket.onclose = event => {
        console.log('Socket Closed Connection: ', event);
    };

    socket.onerror = error => {
        console.log('Socket Error: ', error);
    };
}

function sendMessage() {
    let message = document.getElementById('message_input').value.toString();

    if (socket == null || message === '') {
        return;
    }

    socket.send(message);
    document.getElementById('message_input').value = '';
}
