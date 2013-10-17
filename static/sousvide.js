var MASK_TIMEOUT_MS = 500

var loader;
var loadUntilTrue = undefined;

var mask;
var maskHidden = true;

function hideMask() {
	if (maskHidden) return;
	maskHidden = true;
		d3.select('#mask')
			.transition()
			.style('opacity', '0')
			.remove()
}

function startMaskTimeout() {
	maskHidden = false;
	setTimeout(function() {
		if (!maskHidden) {
			d3.select('#mask')
				.style('display', 'block')
				.style('opacity', '0')
				.transition()
				.style('opacity', '1')
		}
	}, MASK_TIMEOUT_MS);
}

function loaded(data) {
	console.log(data);
	if (loadUntilTrue && loadUntilTrue(data)) {
		$(loader).css('display', 'none');
		loadUntiltrue = undefined;
	}
	hideMask();
}

function attachRequest(elem, path, blinkUntil) {
	$(elem).click(function(e) {
		e.preventDefault();
		loadUntilTrue = blinkUntil

		loader.style.webkitAnimationName = ""
		loader.style.mozAnimationName = ""
		loader.style.webkitAnimationName = "loader-anim"
		loader.style.mozAnimationName = "loader-anim"
		$(loader).css('display', 'block')

		$.ajax({
			url: path,
			type: 'html',
			success: function(resp) {
				console.log("got response to " + path + ": " + resp);
			}
		})
	});
	console.log("sent onclick for elem to " + path)
}

$(document).ready(function() {
	loader = document.getElementById('loader');
	mask = document.getElementById('mask');

	startMaskTimeout();

	initTempElems();
	primeTempCache();

	initTimerElems();
	getTimerData();
})
