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
		document.getElementById("users").innerHTML = "";
		Object.keys(d.msg).forEach((email) => {
			if (email != "admin") {
				div = document.createElement("div");
				div.onmouseover = () => {
					this.children[2].hidden = false;
				}
				div.onmouseout = () => {
					this.children[2].hidden = true;
				}

				span = document.createElement("span");
				text = document.createTextNode(email);
				span.appendChild(text);
				span.classList.add("user");
				span.classList.add(d.msg[email].online ? "online" : "offline");

				div.appendChild(span);

				button = document.createElement("buttom");
				text = document.createTextNode("Delete");
				button.appendChild(text);
				button.classList.add("delete");
				button.onclick = () => {
					deluser(email);
				}

				div.appendChild(button);

				document.getElementById("users").appendChild(div);
			}
		});
	}
};

function deluser(email) {
	socket.send("{'type': 1, 'data': '" + email + "'}");
}

function adduser() {
	socket.send("{'type': 0, 'data': '" + document.getElementById("email").value + "'}");

	return false;
}