## go live reload

This is minimal example of server side rendered website
written in golang that has live auto-reload - meaning:
whenever you change the go code or html template code, 
website will reload automatically in the web browser.

It works by using [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
and [watchexec](https://github.com/watchexec/watchexec) utility.

## dependencies

Install locally:
- go compiler
- [watchexec](https://github.com/watchexec/watchexec): utility for watching filesystem
- [just](https://github.com/casey/just): make replacement

## how to run

To build and run the website in watch mode execute in your terminal:
```
just serve
```

In web browser go to: http://localhost:8080

Then change the code and see website refreshing by itself.