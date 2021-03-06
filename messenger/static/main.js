const GREEN_CIRCLE_EMOJI = '&#128994;';
const RED_CIRCLE_EMOJI = '&#128308';

let socket;

function connect() {
    if (socket != null) {
        socket.close(1000);
        document.getElementById('messages_area').value = '';
    }

    let toUserId = document.getElementById('toUserId').value;

    socket = new WebSocket(`ws://localhost:7004/api/v1/messenger/connect?toUserId=${toUserId}`);

    console.log('Attempting Connection...');

    socket.onopen = () => {
        console.log('Successfully Connected');

        sendAccessToken();

        document.getElementById('connection_status').innerHTML = GREEN_CIRCLE_EMOJI;
    };

    socket.onmessage = message => {
        let json = JSON.parse(message.data);

        if (json.constructor === Array) {
            json.forEach(msg => {
                document.getElementById('messages_area').value += `${msg.from}: ${msg.text}\r\n`;
            });
        } else {
            document.getElementById('messages_area').value += `${json.from}: ${json.text}\r\n`;
        }
    };

    socket.onclose = event => {
        console.log('Socket Closed Connection: ', event);
        document.getElementById('connection_status').innerHTML = RED_CIRCLE_EMOJI + ' ' + event.reason;
    };

    socket.onerror = error => {
        console.log('Socket Error: ', error);
        document.getElementById('connection_status').innerHTML = RED_CIRCLE_EMOJI + ' ' + error.code;
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

function sendAccessToken() {
    let accessToken = document.getElementById('accessToken').value;

    if (socket == null || accessToken === '') {
        return;
    }

    socket.send(accessToken);
}
