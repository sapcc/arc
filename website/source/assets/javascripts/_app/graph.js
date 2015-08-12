var width = 500,
    height = 250;

var color = d3.scale.category10();

var force = d3.layout.force()
    .charge(-120)
    .linkDistance(100)
    .size([width, height]);

var svg = d3.select(".jumbo-logo").append("svg")
    .attr("width", width)
    .attr("height", height);

d3.json("nodes.json", function(error, graph) {
  if (error) throw error;

  force
      .nodes(graph.nodes)
      .links(graph.links)
      .start();

  var link = svg.selectAll(".link")
      .data(graph.links)
    .enter().append("line")
      .attr("class", "link")
      .style("stroke-width", function(d) { return Math.sqrt(d.value); });

	var elem = svg.selectAll(".node").data(graph.nodes)
	var elemEnter = elem.enter()
			.append("g");			
	var circle = elemEnter
			.append("circle")
      .attr("class", "node")
      .attr("r", 10)
      .style("fill",
    function(d) {
      if(d.group == 2){return "red"}
      if(d.group == 1){return "black"}
      return color(d.group);
    });
	var label = elemEnter
			.append("text")
      .attr("class", "node-text")
      .attr("font-size", 14)
			.attr("font-family", "sans-serif")
			.attr("text-anchor", "middle")
      .style("fill", "white")
			.text(function(d) { return d.name; })
      .call(force.drag);

	/*var node = svg.selectAll(".node")
      .data(graph.nodes)
    .enter().append("circle")
      .attr("class", "node")
      .attr("r", 5)
      .style("fill", function(d) { return color(d.group); })
      .call(force.drag);*/
	
  /*node.append("title")
  		.text(function(d) { return d.name; });*/

  force.on("tick", function() {
    link.attr("x1", function(d) { return d.source.x; })
        .attr("y1", function(d) { return d.source.y; })
        .attr("x2", function(d) { return d.target.x; })
        .attr("y2", function(d) { return d.target.y; });
		
				elemEnter.attr("transform", function(d){return "translate("+d.x+","+d.y+")"});
		
    /*node.attr("cx", function(d) { return d.x; })
        .attr("cy", function(d) { return d.y; });*/
  });
});