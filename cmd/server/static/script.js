document.querySelector('.game-tile').addEventListener('click', function() {
    document.getElementById('modal').style.display = 'flex';
    document.body.classList.add('modal-open');
    document.getElementById('game-iframe').src = 'game_tictactoe/index.html';
});

document.getElementById('close-btn').addEventListener('click', function() {
    document.getElementById('modal').style.display = 'none';
    document.body.classList.remove('modal-open');
});

