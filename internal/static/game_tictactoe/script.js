

    resizeBoard();

    const svgElement = document.getElementById('board');
    const count = 15;

    let board = document.getElementById('board');
    let height = board.scrollHeight;
    let width = board.scrollWidth;

    let minSize = Math.min(height, width);
    let size = minSize/count;
    let move = false;

    svgElement.addEventListener('click', (e) => {
        const rect = svgElement.getBoundingClientRect();
        const px = e.clientX - rect.left;
        const py = e.clientY - rect.top;

        const x = Math.floor(px/size) * size + size/2;
        const y = Math.floor(py/size) * size + size/2;


        move = !move;
        //textElement.setAttribute('fill', move ? "#0F2D70" : "#B33B1D");
        if (move) {
            drawCross(x, y, size);
        } else {
            drawCircle(x, y, size);
        }

        svgElement.appendChild(textElement);
    });

    drawBoard(count, size);


function drawCross(x, y, size) {
    const halfSize = size / 3.5; 

    drawLine(x - halfSize, y - halfSize, x + halfSize, y + halfSize, "#B33B1D", '3');
    drawLine(x + halfSize, y - halfSize, x - halfSize, y + halfSize, "#B33B1D", '3');
}

function drawCircle(x, y, size) {
    const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
    circle.setAttribute('cx', x);
    circle.setAttribute('cy', y);
    circle.setAttribute('r', size / 3.5); 
    circle.setAttribute('stroke', "#0F2D70");
    circle.setAttribute('stroke-width', '3');
    circle.setAttribute('fill', 'none');

    svgElement.appendChild(circle);
}

function drawLine(x1, y1, x2, y2, color, sw = '1') {
    const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
    line.setAttribute('x1', x1);
    line.setAttribute('y1', y1);
    line.setAttribute('x2', x2);
    line.setAttribute('y2', y2);
    line.setAttribute('stroke', color);
    line.setAttribute('stroke-width', sw);

    svgElement.appendChild(line);
}

function drawBoard(count, size){
    for (let i = 0; i < count+1; i++) {
        for (let j = 0; j < count+1; j++) {
            drawLine(i*size, 0, i*size, count*size, '#696969');
            drawLine(0, j*size, count*size, j*size, '#696969');
        }
    }
}

function resizeBoard() {
	const container = document.getElementById('canvas-container');
	const board = document.getElementById('board');
	const containerWidth = container.clientWidth;
	const containerHeight = container.clientHeight;
	const size = Math.min(containerWidth, containerHeight);

	board.setAttribute('width', size);
	board.setAttribute('height', size);
}
