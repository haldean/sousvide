var tempElem, absErrElem, targetElem, absErrTdElem, plotElem
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

	var i = data.Temps.length - 1
	var temp = data.Temps[i],
		target = data.Targets[i],
		err = temp - target;

	$(tempElem).text(temp.toFixed(2));
	$(targetElem).text(target.toFixed(2));
	$(absErrElem).text((err >= 0 ? '+' : '') + err.toFixed(2));

	pInputElem.setAttribute('value', data.PidParams.P)
	iInputElem.setAttribute('value', data.PidParams.I)
	dInputElem.setAttribute('value', data.PidParams.D)

	if (err > 0) {
		$(absErrTdElem).removeClass('cold')
		$(absErrTdElem).addClass('hot')
	} else if (err < 0) {
		$(absErrTdElem).removeClass('hot')
		$(absErrTdElem).addClass('cold')
	} else {
		$(absErrTdElem).removeClass('hot')
		$(absErrTdElem).removeClass('cold')
	}
}

$(document).ready(function() {
	tempElem = document.getElementById('temp')
	absErrElem = document.getElementById('abs_err')
	targetElem = document.getElementById('target')
	absErrTdElem = document.getElementById('err_td')
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
