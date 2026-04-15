const modal = document.getElementById('win');
const label = document.getElementById('congratulation');
const table = document.getElementById('computer');

function onEnd(e = -1) {
    if (e === -1)
        return
    modal.style.display = 'flex'
    modal.classList.add(e ? 'modal-win' : 'modal-lose')
    label.innerText = 'You ' + (e ? 'Win' : 'Lose') + '!'
}

function onBeginReply(e) {
    if (e.C === 1) {
        end[0]++
        end[1]++
    }
    onReply(e)
}

function start() {
    modal.style.display = 'none'
    modal.className = 'modal'
    label.innerText = ''
    history = []
    init()
    query(onBeginReply, seq++, 'GET', '/begin')
}

function onReset() {
    start()
}

function reset() {
    query(onReset, seq++, 'GET', '/reset')
}

const end = []
const field = []
const clazz = {0: 'free', 1: 'boat', 2: 'lose', 3: 'boom', 4: 'open'}
const click = []

function point(f, x, y) {
    let e = document.getElementById(f + '.' + y + '.' + x)
    if (f === 1) {
        e.onclick = function () {
            onClick(x, y)
        }
    }
    e.className = 'void'
    e.innerHTML = '&#x25CF;'
    return e
}

function init() {
    for (let f = 0; f < 2; f++) {
        end[f] = 0
        field[f] = []
        for (let y = 0; y < 10; y++) {
            field[f][y] = []
            for (let x = 0; x < 10; x++) {
                field[f][y][x] = point(f, x, y)
            }
        }
    }
}

function query(done, id, method = 'GET', path = '/', body) {
    socket.send([id.toString(), method, path, JSON.stringify(body)].join('\r\n'))
    handle.set(id.toString(), done)
}

function reply(p = ['']) {
    if (!handle.has(p[0]))
        return
    if (p[1] === '200')
        try {
            handle.get(p[0])(JSON.parse(p[3] ? p[3] : '{}'))
        } catch (e) {
        }
    if (p[2] === '#')
        handle.delete(p[0])
}

function onClick(x = -1, y = -1) {
    table.classList.add('disabled')
    for (const i of click) {
        i.className = i.className.toString().split(' ')[0]
    }
    query(onNewReply, seq++, 'GET', '/click', {X: x, Y: y, H: history.at(-1)})
}

async function onNewReply(e) {
    semaphore = semaphore.then(async () => {
        onReply(e, 'new')
        await sleep(20)
    }).catch(err => console.error("onNewReply:", err));
}

let history = []
let semaphore = Promise.resolve();

function onReply(e = {
    F: -1,
    X: -1,
    Y: -1,
    C: -1,
    H: ''
}, mod = '') {
    if (e.F === -1)
        return
    let style = [clazz[e.C]]
    if (e.H) {
        history.push(e.H)
        table.classList.remove('disabled')
    }
    if (e.C === 2 || e.C === 3) {
        style.push(mod)
        click.push(field[e.F][e.Y][e.X])
        if (e.C === 3)
            play(e.F)
    }
    field[e.F][e.Y][e.X].className = style.join(' ')
    field[e.F][e.Y][e.X].onclick = undefined
    if (e.C === 3) {
        end[e.F]--
        if (!end[0])
            onEnd(0)
        if (!end[1])
            onEnd(1)
    }
}

const effect = [
    new Audio('../waves/boom.wav'),
    new Audio('../waves/boom.wav'),
    new Audio('../waves/game.wav'),
    new Audio('../waves/lose.wav'),
];

function play(i) {
    const sound = effect[i].cloneNode();
    sound.play();
    sound.onended = () => sound.remove();
}

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

let seq = 0
let handle
let socket

function reload() {
    handle.clear()
    document.location.reload()
}

window.onload = function () {
    handle = new Map
    socket = new WebSocket(location.origin.replace(/^http/, 'ws'))
    socket.onopen = function () {
        start()
    }
    socket.onmessage = function (e) {
        reply(e.data.toString().split(/\r?\n/))
    }
    socket.onclose = function () {
        reload()
    }
    socket.onerror = function () {
        reload()
    }
}
