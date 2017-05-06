//controller for the blog view
class blogController {
  constructor(website) {
    this.website = website;
  }
  init() {
    //magical incantation
    var that = this;
    $("#addBlogButton").on("click", function() {
      that.website.switchView("addBlogView");
    });
  }
}
//controller for add blog view
class addBlogController {
  constructor(website) {
    this.website = website;
  }
  //wraps blog body in HTML p tags
  //takes an array: [body, result] with result being "" to start
  paragraphize(br) {
    var index = br[0].indexOf("\n");
    //if there are no newlines, just wrap the body
    if (index === -1) {
      return "<p class='lead'>\n" + body + "\n</p>\n"
    } else {
      var wrapped = "<p class='lead'>\n" + br[0].subStr(0, index) + "\n</p>\n";
      //if the newline is the end of the body, we're done
      if (br[0].subStr(index).length <= 1) {
        return br[1] + wrapped;
      }
      //chop off current paragraph, add wrapped paragraph to result, and call again
      return paragraphize([br[0].subStr(index + 1), br[1] + wrapped]);
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
      var requestBody = {
        title: that.titleInput.val(),
        date: "Today's date",
        body: that.paragraphize([that.bodyInput.val(), ""]),
        uname: that.unameInput.val(),
        pword: that.pwordInput.val()
      };
      var submitCallback = function() {
        alert(JSON.stringify(requestBody));
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
