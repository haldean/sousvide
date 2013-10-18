// used to implement a drop-n lowpass filter
var lowpass = 0
var drop = 1

// store an hour of data for graphing
var n = 30 * 60 / drop;
var data = new Array(n);
var target = new Array(n);
var heating = new Array(n);

var HEAT_HEIGHT = 5;

for (var i = 0; i < n; i++) {
	data[i] = 0;
	target[i] = 0;
	heating[i] = 0;
}

var width, height, x, y;

var line = d3.svg.line()
	.interpolate("basis")
    .x(function(d, i) { return x(i); })
    .y(function(d, i) { return y(d); });

var svg, path, targetpath, heatingpath;

function initChart() {
	width = window.innerWidth;
	height = window.innerHeight;

	x = d3.scale.linear()
		.domain([0, n - 1])
		.range([0, width]);

	y = d3.scale.linear()
		.domain([0, 100])
		.range([height, 0]);

	svg = d3.select("body").append("svg")
		.attr("width", width)
		.attr("height", height)
		.append("g");

	targetpath = svg.append("g")
		.append("path")
		.datum(target)
		.attr("class", "targetline")
		.attr("d", line);

	heatingpath = svg.append("g")
		.append("path")
		.datum(heating)
		.attr("class", "heatingline")
		.attr("d", line);

	path = svg.append("g")
		.append("path")
		.datum(data)
		.attr("class", "line")
		.attr("d", line);
}

function addTemps(temps, heatings) {
	for (var i = 0; i < temps.length; i += drop) {
		data.push(temps[i]);
		heating.push(heatings[i] ? HEAT_HEIGHT : 0);
		if (data.length > n) {
			data.shift();
			heating.shift();
		}
	}
	heating[0] = 0;
	heating[n-1] = 0;
}

function pushTemp(temp, isheating) {
	if (lowpass++ % drop) return

	data.push(temp);

	heating.pop();
	heating.push(isheating ? HEAT_HEIGHT : 0);
	heating.push(0);

	heatingpath
		.attr("d", line)
		.attr("transform", null)
		.transition()
		.duration(1000)
		.ease("linear")
		.attr("transform", "translate(" + x(-1) + ",0)")
	path
		.attr("d", line)
		.attr("transform", null)
		.transition()
		.duration(1000)
		.ease("linear")
		.attr("transform", "translate(" + x(-1) + ",0)")

	data.shift();
	heating.shift();
	heating[0] = 0;
}

function setTarget(newT) {
	for (var i = 0; i < n; i++) {
		target[i] = newT;
	}
	targetpath.transition().duration(1000).attr("d", line);
}
