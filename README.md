I'm hosting the source code for my personal website here!

It's nowhere near finished as of yet, in case it's not obvious.

The HTML is treated as a golang text template upon request, with the tags being
filled in with the necessary resources.

This way, the HTML and JavaScript can be organized into different files on the
server, but the client, upon accessing the site, will see a single resource
containing the entire front end. Furthermore, this also prevents there from
being a situation where the front end has access to the server filesystem,
since everything needed for the client is loaded at once.

Current plan is to host a wiki for Terra Duenuo here as well. ;)