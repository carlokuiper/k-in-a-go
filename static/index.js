const startForm = document.getElementById('start-form');
const moveForm = document.getElementById('move-form');

startForm.addEventListener('submit', async function (event) {
    event.preventDefault()
    const response = await fetch(document.URL + "start", {
        method: "POST",
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            m: parseInt(document.getElementById('m').value, 10),
            n: parseInt(document.getElementById('n').value, 10),
            k: parseInt(document.getElementById('k').value, 10)
        })
    });
    if (!response.ok) {
        const message = `An error has occured: ${response.status}`;
        throw new Error(message);
    }
    game = await response.json()
    const rows = game.board.map((r) =>
    "<o>" + r + "</p>").join("")
    document.getElementById('output').innerHTML = rows
});

moveForm.addEventListener('submit', async function (event) {
    event.preventDefault()
    const response = await fetch(document.URL + "move", {
        method: "POST",
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            x: parseInt(document.getElementById('x').value, 10),
            y: parseInt(document.getElementById('y').value, 10),
        })
    });
    debugger;
    if (!response.ok) {
        const message = `An error has occured: ${response.status}`;
        throw new Error(message);
    }
    game = await response.json()
    const rows = game.board.map((r) =>
        "<o>" + r + "</p>").join("")
    document.getElementById('output').innerHTML = rows
});