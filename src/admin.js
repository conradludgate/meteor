let id = 0;

socket = new WebSocket("wss://mesa.cwp.io/adminws/");
socket.onmessage = function (event) {
	d = JSON.parse(event.data);
	if (d.type == 0) {
		d.msg.forEach((m) => {
			span = document.createElement("span");
			text = document.createTextNode(m);
			span.appendChild(text);
			span.classList.add("log");
			span.classList.add((id++) % 2 == 0 ? "even" : "odd");
			//document.getElementById("log").appendChild(span);
			document.getElementById("log").innerHTML = 
				span.outerHTML + 
				document.getElementById("log").innerHTML;
		});
	} else if (d.type == 1) {
		document.getElementById("proc").innerHTML = "Processed " + d.msg[0] + "/" + d.msg[1] + " images";
	} else if (d.type == 2) {
		Object.keys(d.msg).forEach((email) => {
			span = document.createElement("span");
			text = document.createTextNode(email);
			span.appendChild(text);
			span.classList.add("user");
			span.classList.add(d.msg[email].online ? "online" : "offline");
			document.getElementById("users").appendChild(span);
		});
	}
};