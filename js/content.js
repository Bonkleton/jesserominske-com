//master class for website
class content {
  constructor() {
    this.model = new model(this);
  }
  init() {
    //initialize homepage on first load
    $("#content").html($("#aboutView").html());
    //magical incantation
    var that = this;
    //handler for navbar to switch view based on name
    $("#nav span").on("click", function() {
      //change into appropriate page
      var view = $(this).attr("name") + "View";
      $("#content").html($("#" + view).html());
      //initialize proper controllers if necessary
      if (!that.controller) {
        if (view === "blogView") {
          that.controller = new blogController(that).init();
        }
      }
      //TODO: Add jQuery animations
    });
  }
}
