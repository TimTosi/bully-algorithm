var drawZone = d3.select(".demo-content");

function update(data) {
    var nodes = drawZone.selectAll(".node").data(data)

    data.forEach(function(element) {
        nodeID = "node-" + element.nodeID;
        n = d3.select("#" + nodeID);
        if (n.empty() == true) {
            inner = drawZone.append("div")
                .attr("class", "demo-updates mdl-card mdl-shadow--2dp mdl-cell mdl-cell--4-col mdl-cell--4-col-tablet mdl-cell--4-col-desktop")
            inner.append("div")
                .attr("id", nodeID)
                .attr("class", "mdl-cell--12-col node pulse")
                .text( function (d) { return "Leader " + element.leaderID; });
            inner.append("div")
                .attr("class", "mdl-card__actions mdl-card--border")
                .append("p")
                .text("Node-" + element.nodeID);
        } else if (element.leaderID == -1) {
            n.classed("pulse", false)
                .classed("dead", true)
                .text("X");
        } else {
            n.classed("pulse", true)
                .classed("dead", false)
                .text( function (d) { return "Leader " + element.leaderID; });
        }
    }, this);
}

d3.interval(function() {
    d3.json("http://localhost:8080/data", update);
}, 1500);
