//master class for website
class content {
  constructor() {
    this.model = new model(this);
  }
  //function to swtich view in general
  switchView(view, context) {
    //change into appropriate page
    $("#content").html($("#" + view).html());
    //initialize proper controllers if necessary
    if (!context.controller) {
      if (view === "blogView") {
        context.controller = new blogController(context).init();
      }
    }
  }
  init() {
    //magical incantation
    var that = this;
    //initialize homepage on first load
    $("#content").html($("#aboutView").html());
    //handler for navbar to switch view based on name
    $("#nav span").on("click", function() {
      var view = $(this).attr("name") + "View";
      that.switchView(view, that);
      //TODO: Add jQuery animations
    });
  }
}
