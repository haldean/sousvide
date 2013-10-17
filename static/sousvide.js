var loader
var tempElem, absErrElem, targetElem, heatingElem, plotElem, accErrElem
var targetDisplayElem, targetChangeElem, targetInputElem
var pInputElem, iInputElem, dInputElem
var enabledElem, maxErrElem
var timerElem, timerAudio

var enableButton, disableButton, lastEnabled = undefined

var loadUntilTrue = undefined;

function primeTempCache() {
	$.ajax({
		url: '/json',
		type: 'json',
		success: function(resp) {
			var temps = Array();
			for (var i = 0; i < resp.length; i++) {
				temps.push(resp[i].Temp)
			}
			console.log("initialized temps to:")
			console.log(temps)
			addTemps(temps)

			initChart()
			window.onresize = function() {
				console.log("window onresize");
				document.getElementsByTagName("svg")[0].remove();
				initChart();
			};
			getApiData()
		}
	})
}

function getApiData() {
	$.ajax({
		url: '/api_data',
		type: 'json',
		success: function(resp) {
			displayData(resp)
		}
	})
	plotElem.setAttribute('src', '/plot?' + (new Date()).getTime())
	setTimeout(getApiData, 1000)
}

function displayData(data) {
	if (loadUntilTrue && loadUntilTrue(data)) {
		$(loader).css('display', 'none')
		loadUntilTrue = undefined
	}

	var temp = data.Temp,
		target = data.Target,
		err = temp - target;

	if (temp != undefined) {
		pushTemp(temp);
	}

	$(tempElem).text(temp.toFixed(1));
	$(targetElem).text(target.toFixed(2));
	$(absErrElem).text((err >= 0 ? '+' : '') + err.toFixed(2));

	if (data.Enabled && (lastEnabled == false || lastEnabled == undefined)) {
		console.log("dis -> en")
		$(enableButton).addClass('fg-secondary')
		$(enableButton).removeClass('fg-primary')
		$(disableButton).addClass('fg-primary')
		$(disableButton).removeClass('fg-secondary')
		reapplyTheme()
	} else if (!data.Enabled && (lastEnabled == true || lastEnabled == undefined)) {
		console.log("en -> dis")
		$(disableButton).addClass('fg-secondary')
		$(disableButton).removeClass('fg-primary')
		$(enableButton).addClass('fg-primary')
		$(enableButton).removeClass('fg-secondary')
		reapplyTheme()
	}
	lastEnabled = data.Enabled

	$(accErrElem).text(data.AccError.toFixed(2))
	$(maxErrElem).text(data.MaxError.toFixed(2))

	pInputElem.setAttribute('value', data.Pid.P)
	iInputElem.setAttribute('value', data.Pid.I)
	dInputElem.setAttribute('value', data.Pid.D)
}

function getTimerData() {
	$.ajax({
		url: '/timers',
		type: 'json',
		success: function(resp) {
			displayTimers(resp)
		}
	})
	setTimeout(getTimerData, 1000)
}

function durationFormat(nano) {
	var neg = nano < 0
	if (neg) nano *= -1
	var sec = Math.floor(nano / 1e9)
	var min = Math.floor(sec / 60)
	sec -= min * 60
	var hr = Math.floor(min / 60)
	min -= hr * 60
	if (min < 10) min = '0' + min
	if (sec < 10) sec = '0' + sec
	return (neg ? '-' : '') + hr + 'h' + min + 'm' + sec + 's'
}

function makeTimer(timer) {
	tr = document.createElement('tr')

	td = document.createElement('td')
	td.innerHTML = timer.Name + ' (' + durationFormat(timer.SetTime) + ')'

	form = document.createElement('form')
	form.setAttribute('method', 'POST')
	form.setAttribute('action', 'delete_timer')
	del = document.createElement('input')
	del.setAttribute('type', 'submit')
	del.setAttribute('value', 'dismiss')
	form.appendChild(del)
	hid = document.createElement('input')
	hid.setAttribute('type', 'hidden')
	hid.setAttribute('name', 'id')
	hid.setAttribute('value', timer.Id)
	form.appendChild(hid)
	td.appendChild(form)

	tr.appendChild(td)

	td = document.createElement('td')
	td.innerHTML = durationFormat(timer.TimeRemaining)
	$(td).addClass('val')
	if (timer.Expired) {
		$(td).addClass('expired')
	}
	tr.appendChild(td)

	return tr
}

function displayTimers(data) {
	timerElem.innerHTML = ''
	for (var i = 0; i < data.length; i++) {
		timer = data[i]
		console.log(timer)
		if (timer.Expired) {
			timerAudio.play()
		}
		timerElem.appendChild(makeTimer(timer));
	}
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
	//setTheme("#FFF", "#000")

	tempElem = document.getElementById('temp')
	absErrElem = document.getElementById('abs_err')
	targetElem = document.getElementById('target')

	enableButton = document.getElementById('button_enable')
	attachRequest(enableButton, "/enable", function(data) {
		return data.Enabled
	})
	disableButton = document.getElementById('button_disable')
	attachRequest(disableButton, "/disable", function(data) {
		return !data.Enabled
	})

	loader = document.getElementById('loader')

	heatingElem = document.getElementById('heating')
	plotElem = document.getElementById('plot')
	accErrElem = document.getElementById('acc_err')
	maxErrElem = document.getElementById('max_err')
	enabledElem = document.getElementById('enabled')

	targetChangeElem = document.getElementById('target_change')
	targetDisplayElem = document.getElementById('target_display')
	targetInputElem = document.getElementById('target_input')

	pInputElem = document.getElementById('pid_p')
	iInputElem = document.getElementById('pid_i')
	dInputElem = document.getElementById('pid_d')

	targetElem.onclick = function() {
		$(targetDisplayElem).css('display', 'none')
		$(targetChangeElem).css('display', 'inline')
		targetInputElem.setAttribute('value', $(targetElem).text())
	}

	primeTempCache()

	timerElem = document.getElementById('timers')
	timerAudio = document.getElementById('timernoise')
	audioEnable = document.getElementById('enable_audio')
	audioEnable.onclick = function() {
		$(audioEnable).css('display', 'none')
		timerAudio.play()
	}
	getTimerData()
})
