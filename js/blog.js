//controller for the blog view
class blogController {
  constructor(website) {
    this.website = website;
  }
  init() {
    //magical incantation
    var that = this;
    //TODO: make this spawn a modal instead, make addBlog a sub-controller
    $("#addBlogButton").on("click", function() {
      that.website.switchView("addBlogView");
    });
  }
}
