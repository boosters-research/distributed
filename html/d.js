// call the backend api
function callAPI(ep, req, fn, err) {
	const http = new XMLHttpRequest();
	const url = "/api/" + ep;

	// when there is a change 
	http.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			var data = JSON.parse(this.responseText);
			fn(data);
		} else if (this.status == 500) {
			var data = JSON.parse(this.responseText);
			var err = JSON.parse(data.error)
			var error = document.getElementById("error");
			error.innerText = err.Detail;
		}
	}

	http.open("POST", url, true);
	http.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
	http.send(JSON.stringify(req));
}

function delCookie(cname) {
	setCookie(cname, "", 0);
}

function getCookie(cname) {
	let name = cname + "=";
	let decodedCookie = decodeURIComponent(document.cookie);
	let ca = decodedCookie.split(';');

	for(let i = 0; i <ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) == ' ') {
		c = c.substring(1);
			}
		if (c.indexOf(name) == 0) {
			return c.substring(name.length, c.length);
		}
	}

	return "";
}

function setCookie(cname, cvalue, expiry) {
	const d = new Date(expiry * 1000);
	let expires = "expires="+ d.toUTCString();
	document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function timeSince(timestamp) {
	var date = Date.parse(timestamp);

	var seconds = Math.floor((new Date() - date) / 1000);

	var interval = seconds / 31536000;

	if (interval > 1) {
		return Math.floor(interval) + " years";
	}

	interval = seconds / 2592000;
	if (interval > 1) {
		return Math.floor(interval) + " months";
	}

	interval = seconds / 86400;
	if (interval > 1) {
		return Math.floor(interval) + " days";
	}

	interval = seconds / 3600;
	if (interval > 1) {
		return Math.floor(interval) + " hours";
	}

	interval = seconds / 60;
	if (interval > 1) {
		return Math.floor(interval) + " minutes";
	}
	return Math.floor(seconds) + " seconds";
}

// load the about page
function loadAbout() {
	var content = document.getElementById("content");

	// set the content
	content.innerHTML = `
		<p>
		  Distributed is an open source tool for remote teams to stay in sync.
		  Rather than using real time communication tools like Slack and 
		  Discord, Distributed takes an async approach. The truth is, we 
		  do our best work when we're alone, with no one around to disturb
		  us.
		</p>
		<p>
		  So shouldn't our tools take that into consideration?
		</p>
		<p>
		  Going async means changing the workflow. Rather than being immersed
		  in streams of messages with no context, we flip it so that everything 
		  is driven by being context first. Distributed uses the concept of 
		  "Boards" like notice boards, kanban boards and message boards that 
		  we all know and love and turns it into the main form of communication 
		  for remote teams.
		</p>
		<p>
		  Anything that's not urgent goes on the board. Anything else we can 
		  revert back to real time chat or video. But the key is to go offline 
		  first and let people take back control of their schedule and workflow.
		  Like most boards, in Distributed, boards consist of posts and comments.
		  Each board has a name like "announcements", "checkins", "general" and 
		  contains posts with a title to indicate the subject.
		</p>
		<p>
		  Posts and comments can be upvoted like most message boards. This is 
		  how we surface the highest priority items. And most things age out 
		  much like Reddit or Hacker News. The global popularity of these 
		  communities show us something that could also be applied to work.
		</p>
		<p>
		  Like all things, Distributed is a work in progress, but hopefully something 
		  useful for everyone living in a globally distributed world.
		</p>
	`;
}

// load all the posts for a board
function loadBoard(name) {
	// don't process zero length names
	if (name.length == 0) {
		return;
	}

	// lowercase the name
	name = name.toLowerCase();

	var content = document.getElementById("content");

	// clear the content
	content.innerHTML = "";

	callAPI("posts", {"board": name }, function(data) {
		console.log("Got the data for " +  name + ": ",  data);
		if (data.records.length == 0) {
			content.innerHTML = "<p>There are no posts on this board</p>";		
			return;
		}

		for (i = 0; i < data.records.length; i++) {
			var el = document.createElement("div");
			var post = data.records[i];

			// set the title
			var title = document.createElement("h4");
			title.innerText = post.title;
			el.appendChild(title);
	
			// below title nav
			var info = document.createElement("span");
			info.style.fontSize = "small";

			if (post.url.length > 0) {
				info.innerHTML = "<a href='"+ post.url + "'>Link</a> | ";
			}

			// posted by
			info.innerHTML += "Posted by " + post.userName + " " + timeSince(post.created) + " ago";

			// add board name if all
			if (name == "all") {
				var a = "<a href='#" + post.board + "'>" + post.board + "</a>";
				info.innerHTML = info.innerHTML + " to " + a;
			}

			// append the content
			el.appendChild(info);
			content.appendChild(el);
		}
	});
}

// load the login page
function loadLogin() {
	var content = document.getElementById("content");

	content.innerHTML = `
		<form id="login-form", action="#login" onsubmit="login(true); return false">
		<p>
		  <input id="username" name="username" placeholder="Username" type=text minlength="1" required />
		</p>
		<p>
		  <input id="password" name="password" type="password" placeholder="Password" minlength="8" required />
		</p>
		<button>Submit</button>
		</form>

		<h3>Signup</h3>

		<form id="signup-form", action="#signup" onsubmit="signup(true); return false">
		<p>
		  <input id="username" name="username" placeholder="Username" type=text minlength="1" required />
		</p>
		<p>
		  <input id="password" name="password" type="password" placeholder="Password" minlength="8" required />
		</p>
		<p>
		  <input id="email" name="email" type="email" placeholder="Email" required />
		</p>
		<button>Submit</button>
		</form>
	`;
}


// login or authenticate the user
function login(submit) {
	// check if its a login
	if (submit == true) {
		console.log("Login event");
		
		var el = document.getElementById("login-form").elements;

		var username = el['username'].value;
		var password = el['password'].value;

		callAPI("login", {"username": username, "password": password }, function(rsp) {
			var expires = parseInt(rsp.session.expires)
			setCookie("sess", rsp.session.id, expires);

			window.location.href = "/";
			window.location.hash = "";
		})

		return;
	}

	var session = getCookie("sess");

	if (session.length == 0) {
		console.log("bad session", session);
		return
	}

	callAPI("readSession", {"sessionId": session}, function(rsp) {
		var date = new Date();
		var now = Math.floor(Date.now() / 1000)
		var expires = parseInt(rsp.session.expires);

		// session expired
		if (expires < now) {
			return
		}

		var lg = document.getElementById("login")
		lg.innerText = "Logout"
		lg.href = "#logout"
	})
}

// logout the user
function logout() {
	var session = getCookie("sess");

	if (session.length == 0) {
		console.log("bad session", session);
		return;
	}

	callAPI("logout", {"sessionId": session}, function(rsp) {
		delCookie("sess");
		var lg = document.getElementById("login")
		lg.innerText = "Login"
		lg.href = "#login"
		window.location.href = "/";
		window.location.hash = "";
	});
}

function newPost(submit) {
	if (submit == true) {
		console.log("Post event");
		
		var session = getCookie("sess");

		var el = document.getElementById("new-post").elements;

		var title = el['title'].value;
		var board = el['board'].value;
		var url = el['link'].value;
		var text = el['text'].value;

		callAPI("post", {
			"post": {
				"title": title,
				"board": board,
				"url": url,
				"content": text
			},
			"sessionId": session,
		}, function(rsp) {
			window.location.href = "/#" + board;
		})

		return;
	}

	// render the form
	var content = document.getElementById("content");
	content.innerHTML = `
		<form id="new-post", action="#new-post" onsubmit="newPost(true); return false">
		<p>
		  <input id="title" name="title" placeholder="Title" type=text minlength="1" required />
		</p>
		<p>
		  <input id="board" name="board" placeholder="Board" type=text minlength="1" required />
		</p>
		<p>
		  <input id="link" name="link" placeholder="Link" type=url />
		</p>
		<p>
		  <textarea id="text" name="text" placeholder="Text" type=text rows=10 cols=25 /></textarea>
		</p>
		<button>Submit</button>
		</form>
		
	`;
}

// executes on hash reload and first load
function reload() {
	// clear the error on reload
	var error = document.getElementById("error");
	error.innerHTML = "";

	var hash = window.location.hash;

	// get the board name
	var name = hash.substring(1);
	var heading = name.replace("-", " ");

	// get the route
	var route = routes.get(name);

	var title = document.getElementById("board");
	// set the title
	title.innerText = heading;

	// load the route
	if (route != undefined) {
		console.log("Loading route: " + name);
		route();
		return
	}

	// load all boards
	if (hash.length == 0) {
		name = "all";
		title.innerText = "All";
	}

	console.log("Loading board: " + name);

	// load the board
	loadBoard(name);
}

function signup(submit) {
	// check if its a login
	if (submit == true) {
		console.log("Signup event");
		
		var el = document.getElementById("signup-form").elements;

		var username = el['username'].value;
		var password = el['password'].value;
		var email = el['email'].value;

		callAPI("signup", {"username": username, "password": password, "email": email }, function(rsp) {
			setCookie("sess", rsp.session.id, rsp.session.expires);

			window.location.href = "/";
			window.location.hash = "";
		})

		return;
	}
}

// the global router
var routes = new Map();
routes.set("about", loadAbout);
routes.set("login", loadLogin);
routes.set("logout", logout);
routes.set("new-post", newPost);

// when the page is ready, start loading content
document.addEventListener("DOMContentLoaded", function(event) {
	login();
	reload();
})

window.addEventListener('hashchange', function() {
	login();
	reload();
});
