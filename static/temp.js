var ENABLE_TEMP = true;

var tempElem, absErrElem, targetElem, heatingElem, accErrElem
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

	pInputElem = document.getElementById('pid_p')
	iInputElem = document.getElementById('pid_i')
	dInputElem = document.getElementById('pid_d')
}

function primeTempCache() {
	if (!ENABLE_TEMP) {
		loaded(undefined);
		return;
	}

	$.ajax({
		url: '/json',
		type: 'json',
		success: function(resp) {
			var temps = Array();
			var heating = Array();
			for (var i = 0; i < resp.length; i++) {
				temps.push(resp[i].Temp);
				heating.push(resp[i].Heating);
			}
			addTemps(temps, heating);

			initChart();
			window.onresize = function() {
				console.log("window onresize");
				document.getElementsByTagName("svg")[0].remove();
				initChart();
				reapplyTheme();
			};
			getApiData();
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
	loaded(data);

	var temp = data.Temp,
		target = data.Target,
		err = temp - target;

	setTarget(target);
	if (temp != undefined) {
		pushTemp(temp, data.Heating);
	}

	$(tempElem).text(temp.toFixed(1));
	$(targetElem).text(target.toFixed(2));
	$(absErrElem).text((err >= 0 ? '+' : '') + err.toFixed(2));

	if (document.activeElement != targetInputElem) {
		targetInputElem.value = target.toFixed(2);
	}

	if (data.Enabled && (lastEnabled == false || lastEnabled == undefined)) {
		$(enableButton).removeClass('fg-primary');
		$(disableButton).removeClass('fg-secondary');
		$(enableButton).addClass('fg-secondary');
		$(disableButton).addClass('fg-primary');
		setOrange();
	} else if (!data.Enabled && (lastEnabled == true || lastEnabled == undefined)) {
		$(enableButton).removeClass('fg-secondary');
		$(disableButton).removeClass('fg-primary');
		$(enableButton).addClass('fg-primary');
		$(disableButton).addClass('fg-secondary');
		setBlue();
	}
	lastEnabled = data.Enabled;

	$(accErrElem).text(data.AccError.toFixed(2));

	if (document.activeElement != pInputElem &&
			document.activeElement != iInputElem &&
			document.activeElement != dInputElem) {
		pInputElem.setAttribute('value', data.Pid.P);
		iInputElem.setAttribute('value', data.Pid.I);
		dInputElem.setAttribute('value', data.Pid.D);
	}
}
