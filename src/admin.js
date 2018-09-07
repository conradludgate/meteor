socket = new WebSocket("wss://mesa.cwp.io/adminws/");
socket.onmessage = function (event) {
	d = JSON.parse(event.data);
	if (d.type == 0) {
		d.msg.forEach((m) => {
			span = document.createElement("span");
			text = document.createTextNode(m);
			span.appendChild(text);
			span.classList.add("log");
			document.getElementById("log").appendChild(span);
			// document.getElementById("log").innerHTML = 
			// 	"<span class='log'>" + m + "</span>" + 
			// 	document.getElementById("log").innerHTML;
		});
	} else if (d.type == 1) {
		document.getElementById("proc").innerHTML = "Processed " + d.msg[0] + "/" + d.msg[1] + " images";
	} else if (d.type == 2) {
		Object.keys(d.msg).forEach((email) => {
			span = document.createElement("span");
			text = document.createTextNode(email);
			span.appendChild(text);
			span.classList.add("user");
			span.classList.add(d.msg[email] ? "online" : "offline");
			document.getElementById("users").appendChild(span);
		});
	}
};