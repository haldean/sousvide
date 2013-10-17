var primary, secondary

function setTheme(p, s) {
	primary = p
	secondary = s
	reapplyTheme();
}

function setBlue() {
	setTheme("#5116FE", "#45317E");
}

function setOrange() {
	setTheme("#FE4F00", "#FFB213");
}

function reapplyTheme() {
	d3.selectAll('.bg-primary').transition().style('background-color', primary)
	d3.selectAll('.bg-secondary').transition().style('background-color', secondary)
	d3.selectAll('.fg-primary').transition().style('color', primary)
	d3.selectAll('.fg-secondary').transition().style('color', secondary)
	d3.select('.targetline').transition().style('stroke', secondary)
}

