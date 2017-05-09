//controller for add blog view
class addBlogController {
  constructor(website) {
    this.website = website;
  }
  //wraps blog body in HTML p tags tail-recursively
  paragraphize(br) {
    var index = br[0].indexOf("\n");
    //if there are no newlines, end recursion
    if (index === -1) {
      return br[1] + "<p class='lead'>\n" + br[0] + "\n</p>\n";
    } else {
      var wrapped = "<p class='lead'>\n" + br[0].substring(0, index) + "\n</p>\n";
      //if the newline is the end of the body, end recusrion
      if (br[0].substring(index).length <= 1) {
        return br[1] + wrapped;
      }
      //chop off current paragraph, add wrapped paragraph to result, and call again
      return this.paragraphize([br[0].substring(index + 1), br[1] + wrapped]);
    }
  }
  init() {
    //attach to blog form
    this.titleInput = $("#i-blogTitle");
    this.bodyInput = $("#i-blogBody");
    this.unameInput = $("#i-blogUname");
    this.pwordInput = $("#i-blogPword");
    //magical incantation
    var that = this;
    //event handler for blog submit
    $("#addBlogSubmit").on("click", function() {
      var d = new Date();
      var blogDate = d.getMonth().toString() + "/"
                     + d.getDate().toString() + "/"
                     + d.getFullYear().toString();
      var requestBody = {
        title: that.titleInput.val(),
        date: blogDate,
        body: that.paragraphize([String(that.bodyInput.val()), ""]),
        uname: that.unameInput.val(),
        pword: that.pwordInput.val()
      };
      var submitCallback = function() {
        alert("ayy");
      };
      that.website.model.addBlog(requestBody, submitCallback, this);
    });
    //event hanbler for close button
    $("#addBlogCloseButton").on("click", function() {
      //switch to the blog view
      that.website.switchView("blogView");
    });
  }
}
