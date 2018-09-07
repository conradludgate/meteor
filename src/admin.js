socket = new WebSocket("wss://mesa.cwp.io/adminws/");
socket.onmessage = function (event) {
	d = JSON.parse(event.data);
	if (d.type == 0) {
		d.msg.forEach((m) => {
			// span = document.createElement("span");
			// text = document.createTextNode(m);
			// span.appendChild(text);
			// span.class = "log";
			// document.getElementById("log").appendChild(span);
			document.getElementById("log").innerHTML = 
				"<span class='log'>" + m + "</span>" + 
				document.getElementById("log");
		});
	} else {
		document.getElementById("proc").innerHTML = "Processed " + d.msg[0] + "/" + d.msg[1] + " images";
	}
};