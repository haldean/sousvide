var primary, secondary, highlight

function setTheme(p, s, h) {
	primary = p
	secondary = s
	highlight = h
	reapplyTheme();
}

function setBlue() {   setTheme("#5116FE", "#8A98FE", "#000"); }
function setOrange() { setTheme("#FE4F00", "#FFB213", "#000"); }

function reapplyTheme() {
	d3.selectAll('.bg-primary').transition().style('background-color', primary)
	d3.selectAll('.bg-secondary').transition().style('background-color', secondary)
	d3.selectAll('.fg-primary').transition().style('color', primary)
	d3.selectAll('.fg-secondary').transition().style('color', secondary)
	d3.selectAll('.targetline').transition().style('stroke', secondary)
	d3.selectAll('.line').transition().style('stroke', highlight)
}

