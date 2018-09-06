let app = undefined;
let container = undefined;
let image = undefined;

// Default image size
let width = 1177;
let height = 942;
let ratio = 1;

let zoom = 1;

let red = undefined;
let green = undefined;
let cross = undefined;

document.addEventListener("wheel", onZoom, false);

function onZoom(event) {
	zoom *= Math.pow(2, -event.deltaY / 10)
	if (zoom < 1) zoom = 1;
	if (zoom > 512) zoom = 16;

	let mouse = app.renderer.plugins.interaction.mouse.global;
	container.scale.set(zoom);
	container.x = - (zoom-1) * mouse.x;
	container.y = - (zoom-1) * mouse.y;

	let active = container.getChildByName("active");
	if (!active) return;

	active.children.forEach((p) => {
		if (p.name !== "rect") {
			p.scale.set(2/zoom);
			p.x = p.realx - 9/zoom;
			p.y = p.realy - 9/zoom;
		} else {
			let left = p.x;
			let top = p.y;
			let w = p.width - p.lineWidth;
			let h = p.height - p.lineWidth;
			p.destroy();

			g = new PIXI.Graphics();
			g.lineStyle(4/zoom, 0xff0000);
			g.drawRect(0, 0, w, h);
			g.name = "rect";

			g.x = left;
			g.y = top;

			active.addChild(g);
		}
	});
}

window.onload = () => {
	// 15 FPS. It's not a game, just an image viewer
	PIXI.settings.TARGET_FPMS = 15 / 1000;

	// Texture for a red circle
	red = (new PIXI.Graphics())
		.lineStyle(2, 0xff0000)
		.drawCircle(0, 0, 4)
		.generateCanvasTexture(2, 4);

	// Texture for a green circle
	green = (new PIXI.Graphics())
		.lineStyle(2, 0x00ff00)
		.drawCircle(0, 0, 4)
		.generateCanvasTexture(2, 4);

	// Texture for a red cross
	// (Drawing a 4 pointed star with a small inner radius)
	cross = (new PIXI.Graphics())
		.lineStyle(2, 0xff0000)
		.beginFill()
		.drawStar(0, 0, 4, 6, 0.1, Math.PI/4)
		.endFill()
		.generateCanvasTexture(2, 4);

	// Get window size
	let e = window, a = 'inner';
	if ( !( 'innerWidth' in window ) )
	{
	a = 'client';
	e = document.documentElement || document.body;
	}
	let maxwidth = e[ a+'Width' ] - 50;
	let maxheight = e[ a+'Height' ] - 50;

	// Scale to the width
	ratio  *= maxwidth / width;
	height *= maxwidth / width;
	width   = maxwidth;

	// If the image size is too big, make the app window smaller
	if (height > maxheight) {
		ratio *= maxheight / height;
		width *= maxheight / height;
		height = maxheight;
	}

	width = Math.ceil(width);
	height = Math.ceil(height);

	// Create app
	app = new PIXI.Application(width, height, {backgroundColor : 0xfff});

	//Add the canvas that Pixi automatically created for you to the HTML document
	document.getElementById("img").appendChild(app.view);

	//Create the sprite
	container = new PIXI.Container();
	container.interactive = true;
	container
		.on('mousedown', onDown)
		.on('mousemove', onMove);

	//Add the image to the stage
	app.stage.addChild(container);
};

let imagename = undefined;

function loadImage(img) {
	container.removeChildren();

	imagename = img;
	let image = new PIXI.Sprite.fromImage("images/" + img);

	image.anchor.set(0.5);
	image.scale.set(ratio);
	image.x = app.screen.width/2;
	image.y = app.screen.height/2;
	image.name = "image";

	container.addChild(image);
}

function onMove() {
	// Get mouse position
	let mouse = app.renderer.plugins.interaction.mouse.global;

	// Is mouse in the app window?
	if (mouse.x < 0 || mouse.x > app.screen.width ||
		mouse.y < 0 || mouse.y > app.screen.height) { 

		if (zoom === 1) return;

		// If not, reset
		this.scale.set(1);
		this.x = 0;
		this.y = 0;

		let active = container.getChildByName("active");
		if (!active) return;
		
		active.children.forEach((p) => {
			if (p.name !== "rect") {
				p.scale.set(2);
				p.x = p.realx - 9;
				p.y = p.realy - 9;
			} else {
				let left = p.x;
				let top = p.y;
				let w = p.width - p.lineWidth;
				let h = p.height - p.lineWidth;
				p.destroy();

				g = new PIXI.Graphics();
				g.lineStyle(4, 0xff0000);
				g.drawRect(0, 0, w, h);
				g.name = "rect";

				g.x = left;
				g.y = top;

				active.addChild(g);
			}
		});
		return; 
	}
	
	this.scale.set(zoom);
	this.x = - (zoom-1) * mouse.x;
	this.y = - (zoom-1) * mouse.y;
}

function onDown() {
	let mouse = app.renderer.plugins.interaction.mouse.global;

	// Get currently active selection container
	let active = this.getChildByName("active");

	// If there isn't one, make one
	if (!active) {
		active = new PIXI.Container();
		active.name = "active";
		this.addChild(active);
	}

	// Make a circle where the cursor is
	let circle = new PIXI.Sprite(red);
	circle.realx = mouse.x;
	circle.realy = mouse.y;
	circle.scale.set(2/zoom);
	circle.x = mouse.x - 9/zoom;
	circle.y = mouse.y - 9/zoom;

	circle.interactive = true;
	circle
		// If clicked on, delete the circle
		.on("mousedown", () => {
			circle.destroy();
			update_rect(active);
		})
		// If moused over, change texture to a cross to make it obvious
		// that it can be deleted
		.on("mouseover", () => {
			circle.texture = cross;
		})
		// Fix the texture when you move the cursor off the circle
		.on("mouseout", () => {
			circle.texture = red;
		});

	active.addChild(circle);

	update_rect(active);
}

function update_rect(active) {
	// Clear the current rectangle if any
	let rect = active.getChildByName("rect");
	if (!!rect) rect.destroy();

	// If less than 2 points, then there's nothing to do
	if (active.children.length < 2) return;
	
	// Get rectangle bounds
	let left = width, right = 0, top = height, bottom = 0;
	active.children.forEach((p) => {
		if (p.realx < left) left = p.realx;
		if (p.realx > right) right = p.realx;
		if (p.realy < top) top = p.realy;
		if (p.realy > bottom) bottom = p.realy;
	});

	// Delete any points inside of the bounds
	active.children.forEach((p) => {
		if (p.realx > left && p.realx < right &&
			p.realy > top && p.realy < bottom) p.destroy();
	});

	// Draw a rectangle with those bounds
	g = new PIXI.Graphics();
	g.lineStyle(4/zoom, 0xff0000);
	g.drawRect(0, 0, right-left, bottom-top);
	g.name = "rect";

	g.x = left;
	g.y = top;

	active.addChild(g);
}

let id = 0;

function save() {
	// Get currently active selection container
	let active = container.getChildByName("active");

	// If there isn't one, nothing to do
	if (!active) return;

	r = active.getChildByName("rect");

	if (!r) return;

	r.destroy();
	let left = width, right = 0, top = height, bottom = 0;
	active.children.forEach((p) => {
		if (p.realx < left) left = p.realx;
		if (p.realx > right) right = p.realx;
		if (p.realy < top) top = p.realy;
		if (p.realy > bottom) bottom = p.realy;
	});

	// Draw a rectangle with those bounds
	g = new PIXI.Graphics();
	g.lineStyle(2, 0x00ff00);
	g.drawRect(0, 0, right-left, bottom-top);
	g.name = (++a).toString();

	g.x = left;
	g.y = top;

	container.addChild(g);
	active.destroy();
}

function submit() {
	save();
	let l = container.children.length - 1;

	let data = {"image": imagename, "meteors": Array(l)};
	for (let i = 0; i < l; i++) {
		r = container.getChildByName((i+1).toString());
		data.meteors[i] = {"l": r.x, "t": r.y, "r": r.width-2+r.x, "b": r.height-2+r.y};
	}

	let XHR = new XMLHttpRequest();
	XHR.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			d = JSON.parse(XHR.responseText);
			if (d.error == 0) {
				loadImage(d.msg);
			} else {
				console.log(d);
			}
		}
	};
	XHR.open("PUSH", "/submit/", true);
	XHR.send(data);
}