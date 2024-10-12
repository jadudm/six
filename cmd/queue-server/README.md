# queue

A queue server.

## routers

https://www.alexedwards.net/blog/which-go-router-should-i-use

We will probably want host-based routing at some point. (Possibly.)

We'll try chi.

https://github.com/go-chi/chi

// https://medium.com/hprog99/working-with-json-in-golang-a-comprehensive-guide-5a94ca5961a1

/enqueue/{queue}
/dequeue/{queue}

Queues and deques JSON messages.

http PUT http://localhost:6000/enqueue/HEAD type=job domain=www.fac.gov
http GET http://localhost:6000/dequeue/HEAD
