var tempElem, absErrElem, targetElem, heatingElem, accErrElem
var targetDisplayElem, targetChangeElem, targetInputElem
var enableButton, disableButton, lastEnabled = undefined
var pInputElem, iInputElem, dInputElem

function initTempElems() {
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

	heatingElem = document.getElementById('heating')
	accErrElem = document.getElementById('acc_err')

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
}

function primeTempCache() {
	$.ajax({
		url: '/json',
		type: 'json',
		success: function(resp) {
			var temps = Array();
			for (var i = 0; i < resp.length; i++) {
				temps.push(resp[i].Temp)
			}
			addTemps(temps)

			initChart()
			window.onresize = function() {
				console.log("window onresize");
				document.getElementsByTagName("svg")[0].remove();
				initChart();
				reapplyTheme();
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
	setTimeout(getApiData, 1000)
}

function displayData(data) {
	console.log("got data!");
	loaded(data);

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
		console.log("dis -> en");
		$(enableButton).removeClass('fg-primary');
		$(disableButton).removeClass('fg-secondary');
		$(enableButton).addClass('fg-secondary');
		$(disableButton).addClass('fg-primary');
		setOrange();
	} else if (!data.Enabled && (lastEnabled == true || lastEnabled == undefined)) {
		console.log("en -> dis");
		$(enableButton).removeClass('fg-secondary');
		$(disableButton).removeClass('fg-primary');
		$(enableButton).addClass('fg-primary');
		$(disableButton).addClass('fg-secondary');
		setBlue();
	}
	lastEnabled = data.Enabled

	$(accErrElem).text(data.AccError.toFixed(2))

	pInputElem.setAttribute('value', data.Pid.P)
	iInputElem.setAttribute('value', data.Pid.I)
	dInputElem.setAttribute('value', data.Pid.D)
}
