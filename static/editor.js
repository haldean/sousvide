var targetInputElem, displayToggle, editorElem;
var editorVisible = false;

function initEditor() {
	targetInputElem = document.getElementById("target_input");
	displayToggle = document.getElementById("editor_expand");
	editorElem = document.getElementById("editor")

	$(displayToggle).click(function(e) {
		e.preventDefault();
		if (editorVisible) {
			d3.select('#editor').style('margin-left', '0px')
				.transition().style('margin-left', '-300px').duration(500);
			d3.select('#editor_expand')
				.transition()
				.style('opacity', '0')
				.each('end', function() {
					d3.select('#editor_expand')
						.classed("icon-double-angle-right", true)
						.classed("icon-double-angle-left", false)
						.transition().style('opacity', '1');
				});
		} else {
			d3.select('#editor').style('margin-left', '-300px')
				.transition().style('margin-left', '0px').duration(500);
			d3.select('#editor_expand')
				.transition()
				.style('opacity', '0')
				.each('end', function() {
					d3.select('#editor_expand')
						.classed("icon-double-angle-left", true)
						.classed("icon-double-angle-right", false)
						.transition().style('opacity', '1');
				});
		}
		editorVisible = !editorVisible;
	});

	formAjax(document.getElementById("target_form"), "/target", function(data) {
		return !!data;
	});

	formAjax(document.getElementById("pid_form"), "/pid", function(data) {
		return !!data;
	});
}
