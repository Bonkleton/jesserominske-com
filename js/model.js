//controller handling data modeling and routing
class model {
	constructor(website) {
    this.website = website;
  }
	//generalized AJAX request
	sendRequest(type, route, data, callback, context) {
		$.ajax({
			url: route,
			type: type,
			contentType: "application/json",
			dataType: "json",
			headers: {},
			data: JSON.stringify(data),
		}).done(function(data, text, request) {
			if (callback) {
				callback.apply(context, arguments);
			}
		});
	}
	//wrapper for getting blogs
	getBlogs(callback, context) {
		this.sendRequest("GET", "blog/post", {}, callback, context);
	}
	//wrapper for submitting a blog
	addBlog(data, callback, context) {
		this.sendRequest("POST", "blog/post", data, callback, context);
	}
	//TODO: add edit blog functionality
}
