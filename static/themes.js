var primary, secondary
function setTheme(p, s) {
	primary = p
	secondary = s
	reapplyTheme();
}

function reapplyTheme() {
	$('.bg-primary').css('background', primary)
	$('.bg-secondary').css('background', secondary)
	$('.fg-primary').css('color', primary)
	$('.fg-secondary').css('color', secondary)
}
