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
	console.log(event)
	zoom *= Math.pow(2, -event.deltaY / 10)
	if (zoom < 1) zoom = 1;

	let mouse = app.renderer.plugins.interaction.mouse.global;
	container.scale.set(zoom);
	container.x = - (zoom-1) * mouse.x;
	container.y = - (zoom-1) * mouse.y;
}

window.onload = function() {
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
	let maxwidth = e[ a+'Width' ] - 100;
	let maxheight = e[ a+'Height' ] - 100;

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

	// width = Math.ceil(width);
	// height = Math.ceil(height);

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

	loadImage("images/image0000.jpg")
}

function loadImage(img) {
	container.removeChildren();

	let image = new PIXI.Sprite.fromImage(img);

	image.anchor.set(0.5);
	image.scale.set(ratio);
	image.x = app.screen.width/2;
	image.y = app.screen.height/2;

	container.addChild(image);
}

function onMove() {
	// Get mouse position
	let mouse = app.renderer.plugins.interaction.mouse.global;

	// Is mouse in the app window?
	if (mouse.x < 0 || mouse.x > app.screen.width ||
		mouse.y < 0 || mouse.y > app.screen.height) { 

		// If not, reset
		this.scale.set(1);
		this.x = 0;
		this.y = 0;
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
	circle.x = mouse.x - 4 - 1;
	circle.y = mouse.y - 4 - 1;

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
		if (p.x + 5 < left) left = p.x + 5;
		if (p.x + 5 > right) right = p.x + 5;
		if (p.y + 5 < top) top = p.y + 5;
		if (p.y + 5 > bottom) bottom = p.y + 5;
	});

	// Delete any points inside of the bounds
	active.children.forEach((p) => {
		if (p.x + 5 > left && p.x + 5 < right &&
			p.y + 5 > top && p.y + 5 < bottom) p.destroy();
	});

	// Draw a rectangle with those bounds
	g = new PIXI.Graphics();
	g.lineStyle(2, 0xff0000);
	g.drawRect(0, 0, right-left, bottom-top);
	g.name = "rect";

	g.x = left;
	g.y = top;

	active.addChild(g);
}