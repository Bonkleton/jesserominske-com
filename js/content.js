//master class for website
class content {
  constructor() {
    this.model = new model(this);
  }
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
    //magical incantation
    var that = this;
    //initialize homepage on first load
    $("#content").html($("#aboutView").html());
    //handler for navbar to switch view based on href
    $("#nav span").on("click", function() {
      var view = $(this).attr("name") + "View";
      that.switchView(view);
      //TODO: Add jQuery animations
    });
  }
}
