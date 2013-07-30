var tempElem, absErrElem, targetElem, absErrTdElem, heatingElem, plotElem
var targetDisplayElem, targetChangeElem, targetInputElem
var pInputElem, iInputElem, dInputElem

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
	console.log(data)

	var temp = data.Temp,
		target = data.Target,
		err = temp - target;

	$(tempElem).text(temp.toFixed(2));
	$(targetElem).text(target.toFixed(2));
	$(absErrElem).text((err >= 0 ? '+' : '') + err.toFixed(2));

	pInputElem.setAttribute('value', data.Pid.P)
	iInputElem.setAttribute('value', data.Pid.I)
	dInputElem.setAttribute('value', data.Pid.D)

	if (data.Heating) {
		$(heatingElem).addClass('hot')
		$(heatingElem).removeClass('cold')
		$(heatingElem).text('ON')
	} else {
		$(heatingElem).addClass('cold')
		$(heatingElem).removeClass('hot')
		$(heatingElem).text('OFF')
	}
}

$(document).ready(function() {
	tempElem = document.getElementById('temp')
	absErrElem = document.getElementById('abs_err')
	targetElem = document.getElementById('target')
	absErrTdElem = document.getElementById('err_td')
	heatingElem = document.getElementById('heating')
	plotElem = document.getElementById('plot')

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

	getApiData()
})
