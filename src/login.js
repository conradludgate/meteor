function changeTab(tab) {
	if (tab === 0) {
		document.getElementById("confirm").style.display = "none";
		document.getElementById("type").value="login";
	} else {
		document.getElementById("confirm").style.display = "block";
		document.getElementById("type").value="create";
	}

	// Get all elements with class="tablinks" and remove the class "active"
	let tablinks = document.getElementsByClassName("tablinks");
	tablinks[1-tab].className = tablinks[1-tab].className.replace(" active", "");

	tablinks[tab].className += " active";
}