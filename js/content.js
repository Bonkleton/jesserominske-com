//initialize homepage on first load
$("#content").html(document.getElementById("aboutView").innerHTML);
            
//attach handler to navbar to handle loading view content into the DOM
$("#nav a").on("click", function() {
    var view = $(this).attr("href") + "View"; //get the href attribute to identify the item clicked
    $("#content").html(document.getElementById(view).innerHTML); //change into appropriate page
    //TODO: Add jQuery animations
    return false; //disable default href behavior
});