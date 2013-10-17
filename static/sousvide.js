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
	if (loadUntilTrue && loadUntilTrue(data)) {
		$(loader).css('display', 'none');
		loadUntiltrue = undefined;
	}
	hideMask();
}

function startLoader(blinkUntil) {
	loadUntilTrue = blinkUntil;
	loader.style.webkitAnimationName = ""
	loader.style.mozAnimationName = ""
	loader.style.webkitAnimationName = "loader-anim"
	loader.style.mozAnimationName = "loader-anim"
	$(loader).css('display', 'block')
}

function attachRequest(elem, path, blinkUntil) {
	$(elem).click(function(e) {
		e.preventDefault();
		startLoader(blinkUntil);
		$.ajax({
			url: path,
			dataType: 'text',
			type: 'post',
			success: function(resp) {
				console.log("got response to " + path + ": " + resp);
			}
		})
	});
}

function findInputs(root, so_far) {
	if (so_far == undefined) {
		so_far = new Array();
	}
	var children = root.children;
	for (var i = 0; i < children.length; i++) {
		var child = children[i];
		if (child.tagName == "INPUT") {
			so_far.push(child);
		} else if (child.children) {
			so_far = findInputs(child, so_far);
		}
	}
	return so_far;
}

function formAjax(formElem, path, blinkUntil) {
	var submit = (function(e) {
		if (e) e.preventDefault();
		startLoader(blinkUntil);
		var data = {};
		var inputs = findInputs(formElem);
		for (var i = 0; i < inputs.length; i++) {
			if (inputs[i].attributes.type.value == "text") {
				data[inputs[i].attributes.name.value] = inputs[i].value;
			}
		}
		$.ajax({
			url: path,
			data: data,
			type: 'post',
			dataType: 'text',
			success: function(resp) {
				console.log("got response to " + path + " with data " + data + ": " + resp);
			},
			error: function(xhr, stat, err) {
				console.log("error " + err + " on backend: " + xhr.responseText);
			}
		});
	});

	$(formElem).submit(submit);
	var inputs = findInputs(formElem);
	for (var i = 0; i < inputs.length; i++) {
		$(inputs[i]).keyup(function(e) {
			if (e.keyCode == 13) {
				submit();
				e.currentTarget.blur();
			}
		});
	}
}

$(document).ready(function() {
	loader = document.getElementById('loader');
	mask = document.getElementById('mask');

	startMaskTimeout();

	initTempElems();
	primeTempCache();

	initTimerElems();
	getTimerData();

	initEditor();
})
