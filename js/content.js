// handles loading of view content into DOM
$(function() {
    // initialize homepage on first load
    $("#content").html(document.getElementById("homeView").innerHTML);
    
    // attach handler to navbar
    $("#nav li a").on("click", function() {
        var view = $(this).attr("href") + "View"; // get the href attribute to identify the item clicked
        $("#content").html(document.getElementById(view).innerHTML); // change into appropriate page
        
        return false; // disable default href behavior
    });
});