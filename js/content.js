//master class for website
class content {
  //function to swtich view in general
  switchView(view) {
    //condition for which view was chosen
    if (view === "blogView") {
      this.controller = new blogController(this);
    } else if (view === "addBlogView") {
      this.controller = new addBlogController(this);
    }
    //change into appropriate page and initialize if needed
    $("#content").html($("#" + view).html());
    if (this.controller) {
      this.controller.init();
    }
  }
  init() {
    this.model = new model();
    //magical incantation
    var that = this;
    //initialize homepage on first load
    $("#content").html($("#aboutView").html());
    //handler for navbar to switch view based on href
    $("#nav span").on("click", function() {
      //create selector from name atribute
      var view = $(this).attr("name") + "View";
      that.switchView(view);
      //TODO: Add jQuery animations
    });
  }
}
var website;
$(document).ready(function() {
  website = new content().init();
});
