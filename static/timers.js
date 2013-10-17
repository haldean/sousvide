var timerElem, timerAudio

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
