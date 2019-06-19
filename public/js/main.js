var matrixCanvas = document.getElementById("led-matrix");
var ctx = matrixCanvas.getContext("2d");

var id = ctx.createImageData(1, 1);
var data = id.data;

var PixelPitch = 22;
var Gutter = 1;
var margin = 7;

function drawMatrix(matrix) {

    ctx.clearRect(0, 0, matrixCanvas.width, matrixCanvas.height);

    for (var col = 0; col < 32; col++) {
        for (var row = 0; row < 32; row++) {

            var x = col * PixelPitch + Gutter;
            var y = row * PixelPitch + Gutter;

            x += margin * 2;
            y += margin * 2;

            var colorData = matrix[col][row];

            drawPixel(ctx, data, x, y, colorData.R, colorData.G, colorData.B)
        }
    }
}

function drawPixel(context, data, x, y, r, g, b) {

    data[0] = r;
    data[1] = g;
    data[2] = b;
    data[3] = 255;

    context.beginPath();
    context.arc(x, y, 9, 0, 2 * Math.PI, false);

    if (r == 0 && g == 0 && b == 0) {
        context.fillStyle = 'rgb(24,24,24)';
    } else {
        context.fillStyle = 'rgb(' + r + ',' + g + ',' + b + ')';
    }

    context.fill();
    context.stroke();
}

$("#command-form").submit(function (e) {
    e.preventDefault();
    var widget = $(this).find('#widget').val();
    var url = $(this).find('#url').val();

    $.ajax({
        url: "http://localhost:8081/exec?widget=" + widget + '&url=' + url,
        crossDomain: true,
    });
});

window.addEventListener("load", function (evt) {

    var ws = new WebSocket("ws:localhost:8082/pixel");

    ws.onopen = function (evt) {
    };
    ws.onclose = function (evt) {
        ws = null;
    };

    ws.onmessage = function (evt) {
        drawMatrix(JSON.parse(evt.data))
    };

    ws.onerror = function (evt) {
    }
});
