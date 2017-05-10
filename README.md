I'm hosting the source code for my personal website here!

The HTML is treated as a golang text template upon request, with the tags being
filled in with the necessary resources.

This way, the HTML and JavaScript can be organized into different files on the
server, but the client, upon accessing the site, will see a single resource
containing the entire front end. Furthermore, this also prevents there from
being a situation where the front end has access to the server filesystem,
since everything needed for the client is loaded at once.

For the front end, a master site controller instance is initialized upon render,
and this contains facilitation for other controllers to be loaded in upon click
events, all treated as members of the master controller instance.

Current plan is to host a wiki for Terra Duenuo on this server as well. ;)
