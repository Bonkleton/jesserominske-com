//controller for the blog view
class blogController {
  constructor(website) {
    this.website = website;
  }
  //loads blogs into blog view body
  renderBlogs(blogs) {
    var blogsHTML = "";
    for (var i = 0; i < blogs.length; i++) {
      blogsHTML = '<div class="blogContainer">\n'
                  +'<div class="blogHeader">\n'
                  + '<h2 class="header pull-left">' + blogs[i].Title + '</h2>\n</div>'
                  + '<div class="blogDate pull-right">' + blogs[i].Date + '</div>'
                  + '<div class="blogBody pull-left">\n'
                  + blogs[i].Body
                  + '</div>\n</div>' + blogsHTML;
    }
    $("#blogEntriesContainer").html($.parseHTML(blogsHTML));
  }
  init() {
    //get blogs from server
    this.website.model.getBlogs(function(response) {
      this.renderBlogs(response);
    }, this);
    //magical incantation
    var that = this;
    //event handler to initialize modal and controller
    $("#addBlogButton").on("click", function() {
      that.formControl = new addBlogController(that.website).init();
      $("#addBlogModal").modal();
    });
  }
}
