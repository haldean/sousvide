var targetInputElem, displayToggle, editorElem;
var editorVisible = false;

function swapClasses(elemId, start, end) {
  if (animEnabled) {
    d3.select(elemId)
      .transition()
      .style('opacity', '0')
      .each('end', function() {
        d3.select(elemId)
        .classed(start, false)
        .classed(end, true)
        .transition().style('opacity', '1');
      });
  } else {
    d3.select(elemId).classed(start, false).classed(end, true);
  }
}

function initEditor() {
	targetInputElem = document.getElementById("target_input");
	displayToggle = document.getElementById("editor_expand");
	editorElem = document.getElementById("editor")

	$(displayToggle).click(function(e) {
		e.preventDefault();
		if (editorVisible) {
			d3.select('#editor').style('margin-left', '0px')
				.transition().style('margin-left', '-300px').duration(500);
      swapClasses('#editor_expand', 'icon-double-angle-left', 'icon-double-angle-right');
		} else {
			d3.select('#editor').style('margin-left', '-300px')
				.transition().style('margin-left', '0px').duration(500);
      swapClasses('#editor_expand', 'icon-double-angle-right', 'icon-double-angle-left');
		}
		editorVisible = !editorVisible;
	});

	formAjax(document.getElementById("param_form"), "/params", function(data) {
		return !!data;
	});
}
